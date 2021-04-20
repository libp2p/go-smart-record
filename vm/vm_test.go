package vm

import (
	"testing"

	p2ptestutil "github.com/libp2p/go-libp2p-netutil"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
)

var k = "234"

func TestEmptyUpdate(t *testing.T) {
	ctx := ir.DefaultUpdateContext{}
	asm := base.BaseGrammar
	vm := NewVM(ctx, asm)
	p, _ := p2ptestutil.RandTestBogusIdentity()

	in := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "key"}, Value: ir.String{Value: "234"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}

	err := vm.Update(p.ID(), k, in)
	if err != nil {
		t.Fatal(err)
	}
	out := vm.Get(k)
	if !ir.IsEqual(in, *out[p.ID()]) {
		t.Fatal("Record not updated in empty key", in, out)
	}

	out = vm.Get("randomKey")
	if out[p.ID()] != nil {
		t.Fatal("Returned non emtpy record", out)
	}
}

func TestExistingUpdate(t *testing.T) {
	ctx := ir.DefaultUpdateContext{}
	asm := base.BaseGrammar
	vm := NewVM(ctx, asm)
	p, _ := p2ptestutil.RandTestBogusIdentity()

	in1 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}
	in2 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
		},
	}
	in := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
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
	if !ir.IsEqual(in, *out[p.ID()]) {
		t.Fatal("Record not updated in existing key", in, out)
	}
}

func TestSeveralPeers(t *testing.T) {
	ctx := ir.DefaultUpdateContext{}
	asm := base.BaseGrammar
	vm := NewVM(ctx, asm)
	p1, _ := p2ptestutil.RandTestBogusIdentity()
	p2, _ := p2ptestutil.RandTestBogusIdentity()

	in1 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}
	in2 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
		},
	}
	in := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
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
	if !ir.IsEqual(in1, *out[p1.ID()]) || !ir.IsEqual(in2, *out[p2.ID()]) {
		t.Fatal("Record not updated in existing key", in1, in2, out)
	}
	err = vm.Update(p2.ID(), k, in1)
	if err != nil {
		t.Fatal(err)
	}
	out = vm.Get(k)
	if !ir.IsEqual(in, *out[p2.ID()]) {
		t.Fatal("Record not updated in existing key", in1, in2, out)
	}
}
