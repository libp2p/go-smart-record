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
	if !ir.IsEqual(in, out) {
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
	if !ir.IsEqual(in, out) {
		t.Fatal("Record not updated in existing key", in, out)
	}
}

func TestQueryWrongSelector(t *testing.T) {
	// TODO: Make tests with incorrectly formed selectors, etc.
}

func TestQuery(t *testing.T) {
	ctx := ir.DefaultMergeContext{}
	vm := NewVM(ctx)

	in := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p1"}, Value: ir.String{Value: "ssss"}},
			ir.Pair{Key: ir.String{Value: "p2"}, Value: ir.Dict{
				Tag: "foo.2",
				Pairs: ir.Pairs{
					ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
				},
			},
			}},
	}
	in1 := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p2"}, Value: ir.Dict{
				Tag: "foo.2",
				Pairs: ir.Pairs{
					ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
				},
			},
			}},
	}
	in2 := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p1"}, Value: ir.String{Value: "ssss"}},
		},
	}

	selector1 := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p2"}, Value: ir.Dict{Tag: "foo.2"}},
		},
	}

	selector2 := ir.Dict{
		Tag: "foo",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p1"}},
		},
	}

	// Add to the key
	err := vm.Update("234", in)
	if err != nil {
		t.Fatal(err)
	}

	out, err := vm.Query("234", selector1)
	if err != nil {
		t.Fatal(err)
	}
	if !ir.IsEqual(in1, out) {
		t.Fatal("Error querying key", "in:", in1, "out:", out)
	}

	out, err = vm.Query("234", selector2)
	if err != nil {
		t.Fatal(err)
	}
	if !ir.IsEqual(in2, out) {
		t.Fatal("Error querying key", "in:", in2, "out:", out)
	}

	// Check empty selector returns nothing
	out, err = vm.Query("234", ir.Dict{})
	if err != nil {
		t.Fatal(err)
	}
	if !ir.IsEqual(out, ir.Dict{}) {
		t.Fatal("Error querying key", "in:", in2, "out:", out)
	}
}
