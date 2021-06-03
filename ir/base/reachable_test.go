package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p-core/host"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	xr "github.com/libp2p/go-routing-language/syntax"

	"github.com/libp2p/go-smart-record/ir"
)

var unreachableAddr = "/ip4/127.0.0.1/tcp/44783/p2p/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy48tpPV95XdKu1BcvMA"

func setupHost(ctx context.Context, t *testing.T) host.Host {
	return bhost.New(swarmt.GenSwarm(t, ctx, swarmt.OptDisableReuseport))
}

func reachableNode(addr string, conn bool, dial bool) xr.Node {
	maddr1 := xr.Predicate{
		Tag:        "multiaddr",
		Positional: xr.Nodes{xr.String{addr}},
	}
	ops := xr.Nodes{}
	if conn {
		ops = append(ops, xr.String{"connected"})
	}
	if dial {
		ops = append(ops, xr.String{"dialable"})
	}

	return xr.Predicate{
		Tag: "reachable",
		Named: xr.Pairs{
			xr.Pair{Key: xr.String{"address"}, Value: maddr1},
			xr.Pair{Key: xr.String{"how"}, Value: xr.List{
				Elements: ops,
			},
			},
		},
	}
}

func TestAssembly(t *testing.T) {
	p := reachableNode(unreachableAddr, true, true)
	n, err := BaseGrammar.Assemble(ir.AssemblerContext{}, p)
	if err != nil {
		t.Errorf("assemble error: (%v)", err)
	}
	po, ok := n.(*Reachable)
	if !ok {
		t.Errorf("nodes assembled to something different from reachable")
	}
	if !(po.verifyConn && po.verifyDial) {
		t.Errorf("verify flags not set correctly")
	}
}

func TestTriggerReachableDial(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setupHost(ctx, t)
	s := setupHost(ctx, t)
	reachable := fmt.Sprintf("%s/p2p/%s", s.Addrs()[0].String(), s.ID().Pretty())

	// Getting data ready
	p1 := reachableNode(reachable, false, true)
	p2 := reachableNode(unreachableAddr, false, true)
	in := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "unreachable"}, Value: xr.List{xr.Nodes{p2}}},
		},
	}
	d := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "reachable"}, Value: p1},
			xr.Pair{Key: xr.String{Value: "list"}, Value: in},
		},
	}

	// Assemble data structure
	n, err := BaseGrammar.Assemble(ir.AssemblerContext{Grammar: BaseGrammar}, d)
	if err != nil {
		t.Errorf("assemble error: (%v)", err)
	}
	do, ok := n.(*ir.Dict)
	if !ok {
		t.Errorf("nodes assembled to something different from a dict")
	}

	// Trigger reachable.
	TriggerReachable(do, c)

	// Verifications
	r := do.Get(&ir.String{Value: "reachable"})
	if r == nil {
		t.Fatal("dialable node was removed")
	}
	// verifyDial was set
	if !r.(*Reachable).verifyDial || r.(*Reachable).verifyConn {
		t.Errorf("dialable node not verified successfully")
	}
	// List of unreachable was removed
	l := do.Get(&ir.String{Value: "list"}).(*ir.Dict).Get(&ir.String{Value: "unreachable"}).(*ir.List)
	if len(l.Elements) != 0 {
		t.Errorf("list with single unreachable node not removed")
	}
}

func TestTriggerReachableConnect(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setupHost(ctx, t)
	s := setupHost(ctx, t)
	w := setupHost(ctx, t)
	conn := fmt.Sprintf("%s/p2p/%s", s.Addrs()[0].String(), s.ID().Pretty())
	nonConn := fmt.Sprintf("%s/p2p/%s", w.Addrs()[0].String(), w.ID().Pretty())

	// Connect c-s
	if err := c.Connect(ctx, *host.InfoFromHost(s)); err != nil {
		t.Fatal(err)
	}

	// Getting data ready
	p1 := reachableNode(conn, true, false)
	p2 := reachableNode(unreachableAddr, true, false)
	p3 := reachableNode(nonConn, true, false)
	d := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "reachable"}, Value: p1},
			xr.Pair{Key: p2, Value: xr.String{Value: "unreachable"}},
			xr.Pair{Key: xr.String{Value: "not connected"}, Value: p3},
		},
	}

	// Assemble data structure
	n, err := BaseGrammar.Assemble(ir.AssemblerContext{Grammar: BaseGrammar}, d)
	if err != nil {
		t.Errorf("assemble error: (%v)", err)
	}
	do, ok := n.(*ir.Dict)
	if !ok {
		t.Errorf("nodes assembled to something different from a dict")
	}

	// Trigger reachable.
	TriggerReachable(do, c)

	// Verifications
	r := do.Get(&ir.String{Value: "reachable"})
	if r == nil {
		t.Fatal("dialable node was removed")
	}
	// verifyConn was set
	if !r.(*Reachable).verifyConn || r.(*Reachable).verifyDial {
		t.Errorf("dialable node not verified successfully")
	}
	r2 := do.Get(&ir.String{Value: "not connected"})
	if r2 != nil {
		t.Errorf("not connected node was not removed")
	}
	if len(do.Pairs) != 1 {
		t.Errorf("failed removing all pairs that didn't pass verification")
	}
}
