package xr

import (
	"fmt"
	"math/big"

	"github.com/ipld/go-ipld-prime"
	xrIpld "github.com/libp2p/go-smart-record/xr/ipld"
)

// ipldTypeTags used in Node_IPLD type
var ipldTypeTags = []string{
	"String_IPLD",
	"Blob_IPLD",
	"Float_IPLD",
	"Int_IPLD",
	"Bool_IPLD",
	"Dict_IPLD",
	"Set_IPLD",
}

// FromIPLD transforms an IPLD Node into its xr.Node representation
func FromIPLD(n ipld.Node) (Node, error) {
	switch n1 := n.(type) {
	case xrIpld.Blob_IPLD:
		b, err := n1.AsBytes()
		if err != nil {
			return nil, err
		}
		return Blob{b}, nil

	case xrIpld.Bool_IPLD:
		b, err := n1.AsBool()
		if err != nil {
			return nil, err
		}
		return Bool{b}, nil

	case xrIpld.String_IPLD:
		b, err := n1.AsString()
		if err != nil {
			return nil, err
		}
		return String{b}, nil

	case xrIpld.Int_IPLD:
		b, err := n1.AsInt()
		if err != nil {
			return nil, err
		}
		return Int{big.NewInt(b)}, nil

	case xrIpld.Float_IPLD:
		b, err := n1.AsFloat()
		if err != nil {
			return nil, err
		}
		return Float{big.NewFloat(b).SetPrec(64)}, nil

	case xrIpld.Set_IPLD:
		return fromIPLDToSet(n1)

	case xrIpld.Dict_IPLD:
		return fromIPLDToDict(n1)

	case xrIpld.Node_IPLD:
		for _, k := range ipldTypeTags {
			// Check which type is Node_IPLD to convert into the right IPLD Node
			nt, err := n1.LookupByString(k)
			if err == nil {
				return FromIPLD(nt)
			}
		}
		return nil, fmt.Errorf("Node_IPLD has no valid type inside")
	}

	return nil, fmt.Errorf("IPLD type for xr.Node not found. Can't convert.")
}

// Creates a Set in XR from Set_IPLD
func fromIPLDToSet(n xrIpld.Set_IPLD) (Set, error) {
	// Get Tag
	tag, err := n.FieldTag().AsString()
	if err != nil {
		return Set{}, err
	}

	// Get elements
	els := make([]Node, 0)
	li := n.FieldElements().Iterator()
	for !li.Done() {
		_, enode := li.Next()
		n, err := FromIPLD(enode)
		if err != nil {
			return Set{}, err
		}
		// Append element
		els = append(els, n)
	}

	return Set{Tag: tag, Elements: els}, nil
}

// Create Dict in XR from Dict_IPLD
func fromIPLDToDict(n xrIpld.Dict_IPLD) (Dict, error) {
	// Get Tag
	tag, err := n.FieldTag().AsString()
	if err != nil {
		return Dict{}, err
	}

	// Get pairs
	pairs := make([]Pair, 0)
	li := n.FieldPairs().Iterator()
	for !li.Done() {
		_, enode := li.Next()
		// Get key and convert to xr.Node
		ikey := enode.FieldKey()
		k, err := FromIPLD(ikey)
		if err != nil {
			return Dict{}, err
		}
		// Get value and convert to xr.Node
		ivalue := enode.FieldValue()
		v, err := FromIPLD(ivalue)
		if err != nil {
			return Dict{}, err
		}
		// Append pair
		pairs = append(pairs, Pair{Key: k, Value: v})
	}
	return Dict{Tag: tag, Pairs: pairs}, nil

}
