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
	// Has something been verified already?
	verifiedConn bool
	verifiedDial bool

	metadataCtx *meta.Meta
}

// Reachable disassembles to reachable xr.Predicate of the form
// reachable(how=["dialable", "connected"], address=MULTIADDRESS:STRING, checked=BOOL)
// including a checked flag to enable clients to understand if the multiaddr
// was checked by the smart record server.
// NOTE: Setting both ops ("connected" and "dialable") at the same time is quite redundant.
// We check dialability by connecting to the node. We should maybe consider making the
// how field and ir.String, and allow just one or the other.
func (r Reachable) Disassemble() xr.Node {
	how := xr.List{}
	p := xr.Predicate{
		Tag: "reachable",
		Named: xr.Pairs{
			xr.Pair{
				xr.String{"address"}, xr.String{r.addr.String()},
			},
		},
	}

	// Set checked if something has been checked already
	if r.verifiedConn || r.verifiedDial {
		p.Named = append(p.Named, xr.Pair{xr.String{"checked"}, xr.Bool{true}})
	} else {
		p.Named = append(p.Named, xr.Pair{xr.String{"checked"}, xr.Bool{false}})
	}
	// Set the how actions used.
	if r.verifyConn {
		how.Elements = append(how.Elements, xr.String{"connected"})
	}
	if r.verifyDial {
		how.Elements = append(how.Elements, xr.String{"dialable"})
	}
	p.Named = append(p.Named, xr.Pair{xr.String{"how"}, how})

	return p
}

func (r *Reachable) Metadata() meta.MetadataInfo {
	return r.metadataCtx.GetMeta()
}

func (r *Reachable) WritePretty(w io.Writer) error {
	return r.Disassemble().WritePretty(w)
}

func (r *Reachable) UpdateWith(ctx ir.UpdateContext, with ir.Node) error {
	w, ok := with.(*Reachable)
	if !ok {
		return fmt.Errorf("cannot update with a non-reachable node")
	}

	// Update flags, addr stays as it is.
	r.verifyConn = r.verifiedConn || w.verifiedConn
	r.verifyDial = r.verifyDial || w.verifyDial
	r.verifiedConn = r.verifiedConn || w.verifiedConn
	r.verifiedDial = r.verifiedDial || w.verifiedDial

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

// Reachable assemble expects a predicate of the form reachable(how=["dialable", "connected"], address=MULTIADDRESS)
func (ReachableAssembler) Assemble(ctx ir.AssemblerContext, srcNode xr.Node, metadata ...meta.Metadata) (ir.Node, error) {
	// Reachable receives a predicate
	p, ok := srcNode.(xr.Predicate)
	if !ok {
		return nil, fmt.Errorf("smart-tags must be predicates")
	}

	// Get tag and positional arguments.
	tag := p.Tag
	how := getNamed(p, xr.String{"how"})
	addr := getNamed(p, xr.String{"address"})
	// Check tag
	if tag != "reachable" {
		return nil, fmt.Errorf("not a reachable smart tag")
	}

	// Check how actions
	if how == nil {
		return nil, fmt.Errorf("how parameter not found")
	}
	lhow, ok := how.(xr.List)
	if !ok {
		return nil, fmt.Errorf("how parameter is not a list")
	}
	conn := lhow.Elements.IndexOf(xr.String{"connected"})
	dial := lhow.Elements.IndexOf(xr.String{"dialable"})
	if dial == -1 && conn == -1 {
		return nil, fmt.Errorf("no reachable action specified in how")
	}

	// Check multiaddress
	maddr, err := parse.ParseMultiaddr(&parse.ParseCtx{}, addr)
	if err != nil {
		return nil, fmt.Errorf("no valid multiaddr provided")
	}

	// Assemble metadata provided and update assemblyTime
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}

	return &Reachable{
		addr:        maddr,
		verifyConn:  conn >= 0,
		verifyDial:  dial >= 0,
		metadataCtx: m,
	}, nil
}

// TriggerReachable triggers the execution of Reachable verifications
// over a dict and removes Nodes that don't pass the verification.
// NOTE: A failed verification in either the key or the value
// triggers the removal of the whole Pair. We could keep nils
// in the failed verifications, but it would make them hard to manage
// and I can't find a good reason to do this.
func TriggerReachable(d *ir.Dict, h host.Host) {
	// For each pair.
	for k := len(d.Pairs) - 1; k >= 0; k-- {
		if v := triggerReachable(d.Pairs[k].Key, h); !v {
			// If verification failed remove the whole Pair
			d.Remove(d.Pairs[k].Key)
			// No need to check the value
			return
		}
		if v := triggerReachable(d.Pairs[k].Value, h); !v {
			// If verification failed remove the whole Pair
			d.Remove(d.Pairs[k].Key)
		}
	}
}

func triggerReachable(n ir.Node, h host.Host) bool {
	switch n1 := n.(type) {
	case *Reachable:
		return n1.verify(h)
	case *ir.Dict:
		TriggerReachable(n1, h)
	case *ir.List:
		verifyList(n1, h)
	}
	return true
}

func verifyList(s *ir.List, h host.Host) {
	// For each element
	for k := len(s.Elements) - 1; k >= 0; k-- {
		if v := triggerReachable(s.Elements[k], h); !v {
			// Remove element if verification failed
			s.Elements = append(s.Elements[:k], s.Elements[k+1:]...)
		}
	}
	// If all elements of the list removed we keep an empty list
	// no need to return anything.

	// Signal that list can be removed if all elements removed
	// if len(s.Elements) == 0 {
	//         return false
	// }
}

// trigger the verification of reachable
// returns true if the verification was successful,
// and false if it failed, signalling that node can be removed.
func (r *Reachable) verify(h host.Host) bool {
	info, err := peer.AddrInfoFromP2pAddr(r.addr)
	// If there is an error, the verification is not successful.
	if err != nil {
		return false
	}

	// If dialable verification enabled and not checked.
	if r.verifyDial && !r.verifiedDial {
		// Return false if the verification failed.
		// Node needs to be removed.
		if c := checkIfDialable(h, *info); !c {
			return false
		}
		// Set verified flag
		r.verifiedDial = true
	}

	// If connected verification enabled and not checked.
	if r.verifyConn && !r.verifiedConn {
		// Return false if the verification failed.
		// Node needs to be removed.
		if c := checkIfConnected(h, *info); !c {
			return false
		}
		// Set verified flag
		r.verifiedConn = true
	}
	return true
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
