package vm

import (
	"context"
	"testing"
	"time"

	p2ptestutil "github.com/libp2p/go-libp2p-netutil"
	xr "github.com/libp2p/go-routing-language/syntax"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
)

var k = "234"
var gcPeriodOpt = GCPeriod(1 * time.Second)

func TestEmptyUpdate(t *testing.T) {
	ctx := ir.DefaultUpdateContext{}
	asmCtx := ir.AssemblerContext{Grammar: base.BaseGrammar}
	vm, _ := NewVM(context.Background(), ctx, asmCtx, gcPeriodOpt)
	p, _ := p2ptestutil.RandTestBogusIdentity()

	in := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "key"}, Value: xr.String{Value: "234"}},
			xr.Pair{Key: xr.String{Value: "fff"}, Value: xr.String{Value: "ff2"}},
		},
	}

	err := vm.Update(p.ID(), k, in)
	if err != nil {
		t.Fatal(err)
	}
	out := vm.Get(k)
	if !xr.IsEqual(in, *out[p.ID()]) {
		t.Fatal("Record not updated in empty key", in, out)
	}

	out = vm.Get("randomKey")
	if out[p.ID()] != nil {
		t.Fatal("Returned non emtpy record", out)
	}
}

func TestExistingUpdate(t *testing.T) {
	ctx := ir.DefaultUpdateContext{}
	asmCtx := ir.AssemblerContext{Grammar: base.BaseGrammar}
	vm, _ := NewVM(context.Background(), ctx, asmCtx, gcPeriodOpt)
	p, _ := p2ptestutil.RandTestBogusIdentity()

	in1 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "fff"}, Value: xr.String{Value: "ff2"}},
		},
	}
	in2 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "asdf"}, Value: xr.String{Value: "asfd"}},
		},
	}
	in := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "asdf"}, Value: xr.String{Value: "asfd"}},
			xr.Pair{Key: xr.String{Value: "fff"}, Value: xr.String{Value: "ff2"}},
		},
	}

	err := vm.Update(p.ID(), k, in1)
	if err != nil {
		t.Fatal(err)
	}
	err = vm.Update(p.ID(), k, in2)
	if err != nil {
		t.Fatal(err)
	}
	out := vm.Get(k)
	if !xr.IsEqual(in, *out[p.ID()]) {
		t.Fatal("Record not updated in existing key", in, out)
	}
}

func TestSeveralPeers(t *testing.T) {
	ctx := ir.DefaultUpdateContext{}
	asmCtx := ir.AssemblerContext{Grammar: base.BaseGrammar}
	vm, _ := NewVM(context.Background(), ctx, asmCtx, gcPeriodOpt)
	p1, _ := p2ptestutil.RandTestBogusIdentity()
	p2, _ := p2ptestutil.RandTestBogusIdentity()

	in1 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "fff"}, Value: xr.String{Value: "ff2"}},
		},
	}
	in2 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "asdf"}, Value: xr.String{Value: "asfd"}},
		},
	}
	in := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "asdf"}, Value: xr.String{Value: "asfd"}},
			xr.Pair{Key: xr.String{Value: "fff"}, Value: xr.String{Value: "ff2"}},
		},
	}

	err := vm.Update(p1.ID(), k, in1)
	if err != nil {
		t.Fatal(err)
	}
	err = vm.Update(p2.ID(), k, in2)
	if err != nil {
		t.Fatal(err)
	}
	out := vm.Get(k)
	if !xr.IsEqual(in1, *out[p1.ID()]) || !xr.IsEqual(in2, *out[p2.ID()]) {
		t.Fatal("Record not updated in existing key", in1, in2, out)
	}
	err = vm.Update(p2.ID(), k, in1)
	if err != nil {
		t.Fatal(err)
	}
	out = vm.Get(k)
	if !xr.IsEqual(in, *out[p2.ID()]) {
		t.Fatal("Record not updated in existing key", in1, in2, out)
	}
}

func TestGcProcess(t *testing.T) {
	ctx := ir.DefaultUpdateContext{}
	asmCtx := ir.AssemblerContext{Grammar: base.BaseGrammar}
	vm, _ := NewVM(context.Background(), ctx, asmCtx, gcPeriodOpt)
	p, _ := p2ptestutil.RandTestBogusIdentity()
	p2, _ := p2ptestutil.RandTestBogusIdentity()

	in1 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "fff"}, Value: xr.String{Value: "ff2"}},
		},
	}
	in2 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "asdf"}, Value: xr.String{Value: "asfd"}},
		},
	}

	// Small expiration for in1
	err := vm.Update(p.ID(), k, in1, []ir.Metadata{ir.TTL(1 * time.Second)}...)
	if err != nil {
		t.Fatal(err)
	}
	// Add it also in other peer
	err = vm.Update(p2.ID(), k, in1, []ir.Metadata{ir.TTL(1 * time.Second)}...)
	if err != nil {
		t.Fatal(err)
	}
	// Large expiration for in2
	err = vm.Update(p.ID(), k, in2, []ir.Metadata{ir.TTL(3000 * time.Second)}...)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	out := vm.Get(k)
	// In1 should have been garbage collected in p
	if !xr.IsEqual(in2, *out[p.ID()]) {
		t.Fatal("Record not garbage collected successfully", in2, *out[p.ID()])
	}
	// Everything should have been garbage collected in p2
	if out[p2.ID()] != nil {
		t.Fatal("Record not garbage collected successfully in p2", in2, *out[p.ID()])
	}

}

func TestGcFullDict(t *testing.T) {
	d := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "x"}, Value: xr.String{Value: "w"}},
			xr.Pair{Key: xr.String{Value: "w"}, Value: xr.String{Value: "h"}},
		},
	}

	ttl := ir.TTL(1 * time.Second)
	ds, err := ir.SyntacticGrammar.Assemble(ir.AssemblerContext{Grammar: ir.SyntacticGrammar},
		d, []ir.Metadata{ttl}...)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	if g := gcNode(ds); !g {
		t.Fatal("Dict should have been garbage collected", g, ds)
	}
}

func TestGcPartialDict(t *testing.T) {
	in1 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "fff"}, Value: xr.String{Value: "ff2"}},
		},
	}
	in2 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "asdf"}, Value: xr.String{Value: "asfd"}},
		},
	}
	// Small TTL
	ds1, err := ir.SyntacticGrammar.Assemble(ir.AssemblerContext{Grammar: ir.SyntacticGrammar},
		in1, []ir.Metadata{ir.TTL(1 * time.Second)}...)
	if err != nil {
		t.Fatal(err)
	}
	// Large TTL
	ds2, err := ir.SyntacticGrammar.Assemble(ir.AssemblerContext{Grammar: ir.SyntacticGrammar},
		in2, []ir.Metadata{ir.TTL(3000 * time.Second)}...)
	if err != nil {
		t.Fatal(err)
	}
	// Update
	err = ds1.UpdateWith(ir.DefaultUpdateContext{}, ds2)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)
	g := gcNode(ds1)
	if g {
		t.Fatal("Dict should not have been garbage collected", g, ds1)
	}
	if !ir.IsEqual(ds1, ds2) {
		t.Fatal("Dict not garbage collected partially", ds1, ds2)
	}

}
