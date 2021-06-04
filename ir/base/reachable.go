package base

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-routing-language/parse"
	xr "github.com/libp2p/go-routing-language/syntax"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/libp2p/go-smart-record/ir"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

// Reachable is a smart node. It detects if the multiaddrs is
// connected or dialable.
type Reachable struct {
	// Multiaddr checked
	addr ma.Multiaddr
	// What do we need to verify?
	verifyConn bool
	verifyDial bool
	// Has it been verified?
	verifiedConn     bool
	verifiedDial     bool
	verifiedFail     bool
	verifiedFailConn bool

	metadataCtx *meta.Meta
}

// Reachable disassembles to reachable xr.Predicate of the form
// dialed(address=MULTIADDRESS:STRING) if dial checks
// connected(address=MULTIADDRESS:STRING) if connected checks.
// notConnected(address=MULTIADDRESS:STRING) if connected fails.
// notDialable(address=MULTIADDRESS:STRING) if dial fails.
func (r Reachable) Disassemble() xr.Node {
	var tag string

	// Set right tag for predicate
	if r.verifiedConn {
		tag = "connected"
	} else if r.verifiedDial {
		tag = "dialed"
	} else if r.verifiedFail {
		tag = "notDialable"
	} else if r.verifiedFailConn {
		tag = "notConnected"
	} else {
		// This means that nothing has been verified
		// Disassemble the predicate as-is
		if r.verifyConn {
			tag = "connectivity"

		} else if r.verifyDial {
			tag = "dialable"
		}
	}

	return xr.Predicate{
		Tag: tag,
		Named: xr.Pairs{
			xr.Pair{
				xr.String{"address"}, xr.String{r.addr.String()},
			},
		},
	}
}

func (r *Reachable) Metadata() meta.MetadataInfo {
	return r.metadataCtx.Get()
}

func (r *Reachable) WritePretty(w io.Writer) error {
	return r.Disassemble().WritePretty(w)
}

func (r *Reachable) UpdateWith(ctx ir.UpdateContext, with ir.Node) error {
	w, ok := with.(*Reachable)
	if !ok {
		return fmt.Errorf("cannot update with a non-reachable node")
	}

	// Update value
	*r = *w
	// Update metadata
	r.metadataCtx.Update(w.metadataCtx)

	return nil
}

// getNamed returns the xr.Node in a key.
// NOTE: Consider adding this as a function of xr.Predicates
// in the routing-language, and remove it from here.
func getNamed(p xr.Predicate, key xr.Node) xr.Node {
	for _, ps := range p.Named {
		if xr.IsEqual(ps.Key, key) {
			return ps.Value
		}
	}
	return nil
}

type ReachableAssembler struct{}

// Reachable assemble expects a predicate of the form:
// connectivity(address=MULTIADDRESS) or
// dialable(address=MULTIADDRESS)
// See Disassemble() for more info on the resulting predicates after check.
func (ReachableAssembler) Assemble(ctx ir.AssemblerContext, srcNode xr.Node, metadata ...meta.Metadata) (ir.Node, error) {
	// Reachable receives a predicate
	p, ok := srcNode.(xr.Predicate)
	if !ok {
		return nil, fmt.Errorf("smart-tags must be predicates")
	}

	// Get tag and positional arguments.
	tag := p.Tag
	addr := getNamed(p, xr.String{"address"})
	// Check tag
	if tag != "connectivity" && tag != "dialable" {
		return nil, fmt.Errorf("not a reachable smart tag")
	}

	// Check multiaddress
	maddr, err := parse.ParseMultiaddr(&parse.ParseCtx{}, addr)
	if err != nil {
		return nil, fmt.Errorf("no valid multiaddr provided")
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}

	return &Reachable{
		addr:        maddr,
		verifyConn:  tag == "connectivity",
		verifyDial:  tag == "dialable",
		metadataCtx: m,
	}, nil
}

// TriggerReachable triggers the execution of Reachable verifications
// over a dict and adds the appropiate flag to Nodes that don't pass the verification.
func TriggerReachable(d *ir.Dict, h host.Host) {
	// For each pair.
	for k := len(d.Pairs) - 1; k >= 0; k-- {
		triggerReachable(d.Pairs[k].Key, h)
		triggerReachable(d.Pairs[k].Value, h)
	}
}

func triggerReachable(n ir.Node, h host.Host) {
	switch n1 := n.(type) {
	case *Reachable:
		n1.verify(h)
	case *ir.Dict:
		TriggerReachable(n1, h)
	case *ir.List:
		verifyList(n1, h)
	}

}

func verifyList(s *ir.List, h host.Host) {
	// For each element
	for k := len(s.Elements) - 1; k >= 0; k-- {
		triggerReachable(s.Elements[k], h)
	}
}

// trigger the verification of reachable
// updates flags according to the verification result
func (r *Reachable) verify(h host.Host) {
	info, err := peer.AddrInfoFromP2pAddr(r.addr)
	// If there is an error, the verification is not successful.
	if err != nil {
		r.verifiedFail = true
		return
	}

	// If dialable verification enabled and not checked.
	if r.verifyDial && !r.verifiedDial {
		// Set verifyFail if the verification failed.
		if c := checkIfDialable(h, *info); !c {
			r.verifiedFail = true
			return
		}
		// Set verified flag
		r.verifiedDial = true
	}

	// If connected verification enabled and not checked.
	if r.verifyConn && !r.verifiedConn {
		// Set verifyFail if the verification failed.
		if c := checkIfConnected(h, *info); !c {
			r.verifiedFailConn = true
			return
		}
		// Set verified flag
		r.verifiedConn = true
	}
}

// CheckIfdialable Checks if peer reachable with 5s timeout.
func checkIfDialable(h host.Host, i peer.AddrInfo) bool {
	// If self, consider as reachable and don't try to connect
	if h.ID() == i.ID {
		return true
	}
	reachable := false
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := h.Connect(ctx, i); err == nil {
		reachable = true
	}
	cancel()
	return reachable
}

// CheckIfConnected checks if we are currently connected to a peer.
func checkIfConnected(h host.Host, i peer.AddrInfo) bool {
	// If self, consider as connected and don't try to connect
	if h.ID() == i.ID {
		return true
	}
	// Check if we are connected to peer
	for _, p := range h.Network().Peers() {
		if p == i.ID {
			return true
		}
	}
	return false
}
