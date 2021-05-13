package base

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-core/host"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
)

var unreachableAddr = "/ip4/127.0.0.1/tcp/44783/p2p/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy48tpPV95XdKu1BcvMA"

func setupHost(ctx context.Context, t *testing.T) host.Host {
	return bhost.New(swarmt.GenSwarm(t, ctx, swarmt.OptDisableReuseport))
}

func TestReachableDict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setupHost(ctx, t)
	s := setupHost(ctx, t)
	reachable := fmt.Sprintf("%s/p2p/%s", s.Addrs()[0].String(), s.ID().Pretty())
	asm := ReachableAssembler{}
	d := xr.Dict{
		Tag: "reachable",
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "key"}, Value: xr.String{Value: "234"}},
			xr.Pair{Key: xr.String{Value: "NotAddr1"}, Value: xr.String{Value: "/ip4/multiaddr1"}},
			xr.Pair{Key: xr.String{Value: "NotAddr2"}, Value: xr.String{Value: "/ip4/multiaddr2"}},
			xr.Pair{Key: xr.String{Value: "reachable"}, Value: xr.String{Value: reachable}},
			xr.Pair{Key: xr.String{Value: "unreachable"}, Value: xr.String{Value: unreachableAddr}},
		},
	}
	ds, err := asm.Assemble(ir.AssemblerContext{Grammar: BaseGrammar, Host: c}, d)
	if err != nil {
		t.Fatal(err)
	}

	// Verify reachable list.
	r := ds.(*Reachable).Reachable.(*ir.Dict).Get(&ir.String{Value: "reachable"})
	if len(ds.(*Reachable).Reachable.(*ir.Dict).Pairs) != 1 && r == nil {
		t.Fatal("Reachable entry not added correctly to reachable dict")
	}
	// Verify discarded unreachable from user dict
	if ds.(*Reachable).User.(*ir.Dict).Get(&ir.String{Value: "unreachable"}) != nil {
		t.Fatal("Unreachable entry was not removed from user pairs")
	}

	// Verify right number of entries after disassembling
	dsa := ds.Disassemble()
	if len(dsa.(xr.Dict).Pairs) != len(d.Pairs)-1 {
		t.Fatal("Disassembly of reachable dict was not correct")
	}

	// Print result for convenience
	var w bytes.Buffer
	dsa.WritePretty(&w)
	fmt.Println(w.String())
}

func TestSetDict(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setupHost(ctx, t)
	s := setupHost(ctx, t)
	reachable := fmt.Sprintf("%s/p2p/%s", s.Addrs()[0].String(), s.ID().Pretty())
	asm := ReachableAssembler{}
	d := xr.Set{
		Tag: "reachable",
		Elements: xr.Nodes{
			xr.String{Value: "234"},
			xr.String{Value: "/ip4/multiaddr1"},
			xr.String{Value: "/ip4/multiaddr2"},
			xr.String{Value: reachable},
			xr.String{Value: unreachableAddr},
		},
	}
	ds, err := asm.Assemble(ir.AssemblerContext{Grammar: BaseGrammar, Host: c}, d)
	if err != nil {
		t.Fatal(err)
	}

	// Verify reachable list.
	r := ds.(*Reachable).Reachable.(*ir.Set).Elements.IndexOf(&ir.String{Value: reachable})
	if len(ds.(*Reachable).Reachable.(*ir.Set).Elements) != 1 && r < 0 {
		t.Fatal("Reachable entry not added correctly to reachable set")
	}

	// Verify discarded unreachable from user set
	if ds.(*Reachable).User.(*ir.Set).Elements.IndexOf(&ir.String{Value: "unreachable"}) >= 0 {
		t.Fatal("Unreachable entry was not removed from user set")
	}

	// Verify right number of entries after disassembling
	dsa := ds.Disassemble()
	if len(dsa.(xr.Set).Elements) != len(d.Elements)-1 {
		t.Fatal("Disassembly of reachable dict was not correct")
	}

	// Print result for convenience
	var w bytes.Buffer
	dsa.WritePretty(&w)
	fmt.Println(w.String())
}
