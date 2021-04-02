package vm

import (
	"testing"

	"github.com/libp2p/go-smart-record/ir"
)

func TestEmptyUpdate(t *testing.T) {
	ctx := ir.DefaultMergeContext{}
	vm := NewVM(ctx)

	in := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}

	vm.Update("234", in)
	out := vm.Get("234")
	if !ir.IsEqualDict(in, out) {
		t.Fatal("Record not updated in empty key", in, out)
	}
}

// TODO: Add more tests for different merging scenarios.
func TestExistingUpdate(t *testing.T) {
	ctx := ir.DefaultMergeContext{}
	vm := NewVM(ctx)

	in1 := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}
	in2 := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
		},
	}
	in := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}

	err := vm.Update("234", in1)
	if err != nil {
		t.Fatal(err)
	}
	err = vm.Update("234", in2)
	if err != nil {
		t.Fatal(err)
	}
	out := vm.Get("234")
	if !ir.IsEqualDict(in, out) {
		t.Fatal("Record not updated in existing key", in, out)
	}
}
