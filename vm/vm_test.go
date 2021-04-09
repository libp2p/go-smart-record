package vm

import (
	"testing"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
)

func TestEmptyUpdate(t *testing.T) {
	ctx := ir.DefaultMergeContext{}
	asm := base.BaseGrammar
	vm := NewVM(ctx, asm)

	ind := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "key"}, Value: ir.String{Value: "234"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}

	r := base.Record{Key: "234", User: ind}
	in := r.Disassemble()
	err := vm.Update(r.Key, in)
	if err != nil {
		t.Fatal(err)
	}
	out := vm.Get(r.Key)
	if !ir.IsEqual(in, out) {
		t.Fatal("Record not updated in empty key", in, out)
	}
}

// TODO: Add more tests for different merging scenarios.
func TestExistingUpdate(t *testing.T) {
	ctx := ir.DefaultMergeContext{}
	asm := base.BaseGrammar
	vm := NewVM(ctx, asm)

	ind1 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}
	ind2 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
		},
	}
	ind := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}

	r1 := base.Record{Key: "234", User: ind1}
	r2 := base.Record{Key: "234", User: ind2}
	r := base.Record{Key: "234", User: ind}

	in1 := r1.Disassemble()
	in2 := r2.Disassemble()
	in := r.Disassemble()

	err := vm.Update(r.Key, in1)
	if err != nil {
		t.Fatal(err)
	}
	err = vm.Update(r.Key, in2)
	if err != nil {
		t.Fatal(err)
	}
	out := vm.Get(r.Key)
	if !ir.IsEqual(in, out) {
		t.Fatal("Record not updated in existing key", in, out)
	}
}

func TestQuery(t *testing.T) {
	ctx := ir.DefaultMergeContext{}
	asm := base.BaseGrammar
	vm := NewVM(ctx, asm)

	ind := ir.Dict{
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
	ind1 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p2"}, Value: ir.Dict{
				Tag: "foo.2",
				Pairs: ir.Pairs{
					ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
				},
			},
			}},
	}
	ind2 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p1"}, Value: ir.String{Value: "ssss"}},
		},
	}

	r1 := base.Record{Key: "234", User: ind1}
	r2 := base.Record{Key: "234", User: ind2}
	r := base.Record{Key: "234", User: ind}

	in1 := r1.Disassemble()
	in2 := r2.Disassemble()
	in := r.Disassemble()

	dict1 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p2"}, Value: ir.Dict{Tag: "foo.2"}},
		},
	}

	dict2 := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "p1"}},
		},
	}
	selector1 := SyntacticDictSelector{dict1}
	selector2 := SyntacticDictSelector{dict2}

	// Add to the key
	err := vm.Update(r.Key, in)
	if err != nil {
		t.Fatal(err)
	}

	out, err := vm.Query(r.Key, selector1)
	if err != nil {
		t.Fatal(err)
	}
	if !ir.IsEqual(in1, out) {
		t.Fatal("Error querying key", "in:", in1, "out:", out)
	}

	out, err = vm.Query(r.Key, selector2)
	if err != nil {
		t.Fatal(err)
	}
	if !ir.IsEqual(in2, out) {
		t.Fatal("Error querying key", "in:", in2, "out:", out)
	}

	// Check empty selector returns empty record.
	out, err = vm.Query(r.Key, SyntacticDictSelector{})
	if err != nil {
		t.Fatal(err)
	}
	if !ir.IsEqual(out, base.Record{Key: "234"}.Disassemble()) {
		t.Fatal("Error querying key empty selector", "in:", in, "out:", out)
	}
}

func TestAssemblingSymmetry(t *testing.T) {
	// Checks the record assembly symmetry
	ind := ir.Dict{
		Tag: "fff",
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "key"}, Value: ir.String{Value: "234"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}
	asm := base.BaseGrammar
	r := base.Record{Key: "234", User: ind}
	in := r.Disassemble()
	as, err := asm.Assemble(ir.AssemblerContext{Grammar: base.BaseGrammar}, in)
	if err != nil {
		t.Fatal(err)
	}
	if !ir.IsEqual(as, r) {
		t.Fatal("e2e record assembly failed", "in:", r, "assembly:", as)
	}
}

// func TestQueryWrongSelectors(t *testing.T) {
// 	// TODO: Make tests with incorrectly formed selectors, etc.
// 	ctx := ir.DefaultMergeContext{}
// 	asm := base.BaseGrammar
// 	vm := NewVM(ctx, asm)

// 	sameKeys := ir.Dict{
// 		Tag: "foo",
// 		Pairs: ir.Pairs{
// 			ir.Pair{Key: ir.String{Value: "foo"}, Value: ir.String{Value: "ssss"}},
// 			ir.Pair{Key: ir.String{Value: "foo"}, Value: ir.Dict{
// 				Tag: "foo.2",
// 				Pairs: ir.Pairs{
// 					ir.Pair{Key: ir.String{Value: "asdf"}, Value: ir.String{Value: "asfd"}},
// 				},
// 			},
// 			}},
// 	}

// 	dict1 := ir.Dict{
// 		Tag: "foo",
// 		Pairs: ir.Pairs{
// 			ir.Pair{Key: ir.String{Value: "p2"}, Value: ir.Dict{Tag: "foo.2"}},
// 		},
// 	}

// 	selector1 := SyntacticDictSelector{dict1}

// 	// Add to the key
// 	err := vm.Update("234", sameKeys)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	out, err := vm.Query("234", selector1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !ir.IsEqual(ir.Dict{Tag: "foo"}, out) {
// 		t.Fatal("Error querying key", "in:", sameKeys, "out:", out)
// 	}

// }
