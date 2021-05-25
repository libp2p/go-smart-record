package base

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
)

// Reachable is a smart node. It detects if there are multiaddrs in the node.
// If there are, it checks if they are dialable, and it only keeps those
// that are reachable, discarding the unreachable.
type Reachable struct {
	// Reachable only keeps multiaddrss that are reachable in Node.
	Reachable ir.Node
	// User holds user fields which are not multiaddrs.
	User ir.Node
	// Flag to discern if the Reachable node is of type "connected".
	// If this flag is not set it means is of type "dialable".
	isConn bool
}

func (r Reachable) Disassemble() xr.Node {
	_, dok := r.Reachable.(*ir.Dict)
	tag := "dialable"
	// Check the type of Reachable: "dialable" or "connected"
	if r.isConn {
		tag = "connected"
	}
	if dok {
		return (&ir.Dict{
			Tag: tag,
			Pairs: ir.MergePairs(
				r.Reachable.(*ir.Dict).Pairs, // List of reachable multiaddresses
				r.User.(*ir.Dict).Pairs,      // The rest of pairs which don't have multiaddrs.
			),
		}).Disassemble()
	} else {

		return (&ir.Set{
			Tag: tag,
			Elements: ir.MergeElements(
				r.Reachable.(*ir.Set).Elements, // List of reachable multiaddresses
				r.User.(*ir.Set).Elements,      // The rest of pairs which don't have multiaddrs.
			),
		}).Disassemble()
	}

}

func (r *Reachable) Metadata() ir.MetadataInfo {
	return r.User.Metadata()
}

func (r *Reachable) WritePretty(w io.Writer) error {
	return r.Disassemble().WritePretty(w)
}

func (r *Reachable) UpdateWith(ctx ir.UpdateContext, with ir.Node) error {
	w, ok := with.(*Reachable)
	if !ok {
		return fmt.Errorf("cannot update with a non-reachable node")
	}
	// Update each of the dicts for reachable and user straightaway
	var err error
	err = r.User.UpdateWith(ctx, w)
	if err != nil {
		return fmt.Errorf("Error updating user node: %s", err)
	}
	err = r.Reachable.UpdateWith(ctx, w)
	if err != nil {
		return fmt.Errorf("Error updating user node: %s", err)
	}

	return nil
}

type ReachableAssembler struct{}

func (ReachableAssembler) Assemble(ctx ir.AssemblerContext, srcNode xr.Node, metadata ...ir.Metadata) (ir.Node, error) {
	// Check if host set in context
	if ctx.Host == nil {
		return nil, fmt.Errorf("can't assemble reachable node without host in assembler context")
	}
	// Reachable can receive a Dict or Set as input.
	d, dok := srcNode.(xr.Dict)
	s, sok := srcNode.(xr.Set)
	if !dok && !sok {
		return nil, fmt.Errorf("expecting dict or set")
	}
	if dok {
		return reachableDictAssemble(ctx, d, metadata...)
	}
	return reachableSetAssemble(ctx, s, metadata...)

}

func reachableDictAssemble(ctx ir.AssemblerContext, d xr.Dict, metadata ...ir.Metadata) (ir.Node, error) {

	isConn := false
	if d.Tag != "connected" && d.Tag != "dialable" {
		return nil, fmt.Errorf("expecting tag 'connected' or 'dialable'")
	}
	// If the node is of type connected set flag
	if d.Tag == "connected" {
		isConn = true
	}

	u := xr.Dict{}
	r := xr.Dict{}
	for _, p := range d.Pairs {
		info := isValidMultiAddrNode(p.Value)
		// If not a multiaddr add to user set and continue
		if info == nil {
			// Add non-multiaddr to user-dict
			u.Pairs = append(u.Pairs, p)
			continue
		}
		// According to if connected or dialable
		if isConn {
			// If connected add pair with multiaddr to reachable
			if conn := checkIfConnected(ctx.Host, *info); conn {
				r.Pairs = append(r.Pairs, p)
			}
		} else {
			// If dialable add pair with multiaddr to reachable
			if dialable := checkIfDialable(ctx.Host, *info); dialable {
				r.Pairs = append(r.Pairs, p)
			}
		}
	}
	// Assemble reachable and user dicts.
	asm := ir.DictAssembler{}
	uasm, err := asm.Assemble(ctx, u, metadata...)
	if err != nil {
		return nil, fmt.Errorf("couldn't assemble user dict: %s", err)
	}
	rasm, err := asm.Assemble(ctx, r, metadata...)
	if err != nil {
		return nil, fmt.Errorf("couldn't assemble reachable dict: %s", err)
	}
	return &Reachable{
		Reachable: rasm,
		User:      uasm,
		isConn:    isConn,
	}, nil
}

func reachableSetAssemble(ctx ir.AssemblerContext, d xr.Set, metadata ...ir.Metadata) (ir.Node, error) {
	isConn := false
	if d.Tag != "connected" && d.Tag != "dialable" {
		return nil, fmt.Errorf("expecting tag 'connected' or 'dialable'")
	}
	// If the node is of type connected set flag
	if d.Tag == "connected" {
		isConn = true
	}

	u := xr.Set{}
	r := xr.Set{}
	for _, p := range d.Elements {
		info := isValidMultiAddrNode(p)
		// If not a multiaddr add to user set and continue
		if info == nil {
			// Add non-multiaddr to user-set
			u.Elements = append(u.Elements, p)
			continue
		}
		// According to if connected or dialable
		if isConn {
			// If connected add pair with multiaddr to reachable
			if conn := checkIfConnected(ctx.Host, *info); conn {
				r.Elements = append(r.Elements, p)
			}
		} else {
			// If dialable add pair with multiaddr to reachable
			if dialable := checkIfDialable(ctx.Host, *info); dialable {
				r.Elements = append(r.Elements, p)
			}
		}
	}
	// Assemble reachable and user dicts.
	asm := ir.SetAssembler{}
	uasm, err := asm.Assemble(ctx, u, metadata...)
	if err != nil {
		return nil, fmt.Errorf("couldn't assemble user dict: %s", err)
	}
	rasm, err := asm.Assemble(ctx, r, metadata...)
	if err != nil {
		return nil, fmt.Errorf("couldn't assemble reachable dict: %s", err)
	}
	return &Reachable{
		Reachable: rasm,
		User:      uasm,
		isConn:    isConn,
	}, nil
}

// isValidMultiAddrNode checks if the node is of type multiaddr+
// and returns its corresponding AddrInfo
func isValidMultiAddrNode(n xr.Node) *peer.AddrInfo {
	// Check if the value is of type string
	s, ok := n.(xr.String)
	if !ok {
		return nil
	}
	// Check if multiaddr and extract addrinfo
	info, err := peer.AddrInfoFromString(s.Value)
	if err != nil {
		return nil
	}
	return info
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
