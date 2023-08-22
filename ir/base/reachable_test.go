package base

import (
	"context"
	"fmt"
	"testing"

	"github.com/libp2p/go-libp2p/core/host"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	swarmt "github.com/libp2p/go-libp2p/p2p/net/swarm/testing"
	xr "github.com/libp2p/go-routing-language/syntax"

	"github.com/libp2p/go-smart-record/ir"
)

var unreachableAddr = "/ip4/127.0.0.1/tcp/44783/p2p/12D3KooWKRyzVWW6ChFjQjK4miCty85Niy48tpPV95XdKu1BcvMA"

func setupHost(ctx context.Context, t *testing.T) host.Host {
	h, err := bhost.NewHost(swarmt.GenSwarm(t, swarmt.OptDisableReuseport), nil)
	if err != nil {
		panic(err)
	}
	return h
}

func reachableNode(addr string, conn bool) xr.Node {
	var tag string
	maddr1 := xr.Predicate{
		Tag:        "multiaddr",
		Positional: xr.Nodes{xr.String{Value: addr}},
	}
	if conn {
		tag = "connectivity"
	} else {
		tag = "dialable"
	}

	return xr.Predicate{
		Tag: tag,
		Named: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "address"}, Value: maddr1},
		},
	}
}

func TestAssembly(t *testing.T) {
	p := reachableNode(unreachableAddr, true)
	n, err := BaseGrammar.Assemble(ir.AssemblerContext{}, p)
	if err != nil {
		t.Errorf("assemble error: (%v)", err)
	}
	po, ok := n.(*Reachable)
	if !ok {
		t.Errorf("nodes assembled to something different from reachable")
	}
	if !po.verifyConn {
		t.Errorf("verify flags not set correctly")
	}

	p = reachableNode(unreachableAddr, false)
	n, err = BaseGrammar.Assemble(ir.AssemblerContext{}, p)
	if err != nil {
		t.Errorf("assemble error: (%v)", err)
	}
	po, ok = n.(*Reachable)
	if !ok {
		t.Errorf("nodes assembled to something different from reachable")
	}
	if !po.verifyDial {
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
	p1 := reachableNode(reachable, false)
	p2 := reachableNode(unreachableAddr, false)
	in := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "unreachable"}, Value: xr.List{Elements: xr.Nodes{p2}}},
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
	if !r.(*Reachable).verifiedDial || r.(*Reachable).verifiedConn {
		t.Errorf("dialable node not verified successfully")
	}
	// List of unreachable failed
	l := do.Get(&ir.String{Value: "list"}).(*ir.Dict).Get(&ir.String{Value: "unreachable"}).(*ir.List)
	el, ok := (*l).Elements[0].(*Reachable)
	if !ok {
		t.Fatal("list element not of type reachable")
	}
	if !el.verifiedFail {
		t.Errorf("list with single unreachable flags not set correctly")
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
	p1 := reachableNode(conn, true)
	p2 := reachableNode(unreachableAddr, true)
	p3 := reachableNode(nonConn, true)
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
	// verifyConn was set
	if !r.(*Reachable).verifyConn || r.(*Reachable).verifyDial {
		t.Errorf("connected node not verified successfully")
	}

	// Verify unreachable and not connected.
	r2 := do.Get(&ir.String{Value: "not connected"})
	if !r2.(*Reachable).verifiedFailConn || r2.(*Reachable).verifiedConn {
		t.Errorf("not connected flags not set correctly")
	}
	r2 = do.Pairs[1].Key
	if !r2.(*Reachable).verifiedFailConn {
		t.Errorf("unreachable flags not set correctly")
	}
}

func TestTriggerDisassemble(t *testing.T) {
	c := setupHost(context.Background(), t)
	// Connectivity
	r := &Reachable{
		addr:       c.Addrs()[0],
		verifyConn: true,
	}
	d := r.Disassemble()
	do, ok := d.(xr.Predicate)
	if !ok {
		t.Fatal("reachable didn't disassemble to predicate", do)
	}
	if do.Tag != "connectivity" {
		t.Errorf("connectivity predicate didn't disassemble correctly")
	}

	// dialable
	r = &Reachable{
		addr:       c.Addrs()[0],
		verifyDial: true,
	}
	d = r.Disassemble()
	do, ok = d.(xr.Predicate)
	if !ok {
		t.Fatal("reachable didn't disassemble to predicate", do)
	}
	if do.Tag != "dialable" {
		t.Errorf("dialable predicate didn't disassemble correctly")
	}

	// dialed
	r = &Reachable{
		addr:         c.Addrs()[0],
		verifiedDial: true,
	}
	d = r.Disassemble()
	do, ok = d.(xr.Predicate)
	if !ok {
		t.Fatal("reachable didn't disassemble to predicate", do)
	}
	if do.Tag != "dialed" {
		t.Errorf("dialed predicate didn't disassemble correctly")
	}

	// not connected
	r = &Reachable{
		addr:             c.Addrs()[0],
		verifiedFailConn: true,
	}
	d = r.Disassemble()
	do, ok = d.(xr.Predicate)
	if !ok {
		t.Fatal("reachable didn't disassemble to predicate", do)
	}
	if do.Tag != "notConnected" {
		t.Errorf("notConnected predicate didn't disassemble correctly")
	}
}
