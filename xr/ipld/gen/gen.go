//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"

	"github.com/ipld/go-ipld-prime/schema"
	gengo "github.com/ipld/go-ipld-prime/schema/gen/go"
)

func main() {
	ts := schema.TypeSystem{}
	ts.Init()
	adjCfg := &gengo.AdjunctCfg{
		CfgUnionMemlayout: map[schema.TypeName]string{
			"Node_IPLD": "interface", // Use a more pointer-heavy memory layout for this type; this is necessary because it is recursive.
		},
	}

	pkgName := "xr"

	// This needs to preclude to have Type Kinds available.
	ts.Accumulate(schema.SpawnString("String"))
	// ts.Accumulate(schema.SpawnBytes("Blob"))
	// ts.Accumulate(schema.SpawnInt("Int"))
	// ts.Accumulate(schema.SpawnFloat("Float"))
	// ts.Accumulate(schema.SpawnBool("Bool"))

	ts.Accumulate(schema.SpawnString("String_IPLD"))
	ts.Accumulate(schema.SpawnBytes("Blob_IPLD"))
	ts.Accumulate(schema.SpawnInt("Int_IPLD"))
	ts.Accumulate(schema.SpawnFloat("Float_IPLD"))
	ts.Accumulate(schema.SpawnBool("Bool_IPLD"))

	ts.Accumulate(schema.SpawnUnion("Node_IPLD",
		[]schema.TypeName{ // Note that these are somewhat redundant statements due to reasons having to do with how we defined the schema DMT.  The DSL is not so redundant.
			"String_IPLD",
			"Blob_IPLD",
			"Int_IPLD",
			"Float_IPLD",
			"Bool_IPLD",
			"Dict_IPLD",
			"Set_IPLD",
		},

		// Keys for the JSON representation of types
		schema.SpawnUnionRepresentationKeyed(map[string]schema.TypeName{
			"String": "String_IPLD",
			"Blob":   "Blob_IPLD",
			"Int":    "Int_IPLD",
			"Float":  "Float_IPLD",
			"Bool":   "Bool_IPLD",
			"Dict":   "Dict_IPLD",
			"Set":    "Set_IPLD",
		}),
	))

	ts.Accumulate(schema.SpawnStruct("Dict_IPLD",
		[]schema.StructField{
			// Notice the lack of field called "type" -- that is expressed in the keyed union, which wraps this, instead: therefore it's not necessary to repeat here.
			schema.SpawnStructField("Tag", "String", true, false),        // The bools here say "is optional; is not nullable".
			schema.SpawnStructField("Pairs", "Pairs_IPLD", false, false), // I think it may be possible to just use a map here.  (IPLD maps are order-preserving.)  But we'd want to discuss that; I'm not sure I know all desires on this structure.  (E.g., repeat keys?)
		},
		schema.SpawnStructRepresentationMap(nil),
	))
	ts.Accumulate(schema.SpawnList("Pairs_IPLD",
		"Pair_IPLD", false,
	))
	ts.Accumulate(schema.SpawnStruct("Pair_IPLD",
		[]schema.StructField{
			schema.SpawnStructField("Key", "Node_IPLD", false, false),
			schema.SpawnStructField("Value", "Node_IPLD", false, false),
		},
		schema.SpawnStructRepresentationMap(nil),
	))

	ts.Accumulate(schema.SpawnStruct("Set_IPLD",
		[]schema.StructField{
			// Notice the lack of field called "type" -- that is expressed in the keyed union, which wraps this, instead: therefore it's not necessary to repeat here.
			schema.SpawnStructField("Tag", "String", true, false), // The bools here say "is optional; is not nullable".
			schema.SpawnStructField("Elements", "Nodes_IPLD", false, false),
		},
		schema.SpawnStructRepresentationMap(nil),
	))
	ts.Accumulate(schema.SpawnList("Nodes_IPLD",
		"Node_IPLD", false,
	))

	if errs := ts.ValidateGraph(); errs != nil {
		for _, err := range errs {
			fmt.Printf("- %s\n", err)
		}
		os.Exit(1)
	}

	gengo.Generate(".", pkgName, ts, adjCfg)
}
