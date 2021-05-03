package xr

import (
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
	xrIpld "github.com/libp2p/go-smart-record/xr/ipld"
)

type Bool struct {
	Value bool
}

func (b Bool) WritePretty(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, "%v", b.Value)
	return err
}

func IsEqualBool(x, y Bool) bool {
	return x.Value == y.Value
}

func (b Bool) EncodeJSON() (interface{}, error) {
	return struct {
		Type  marshalType `json:"type"`
		Value bool        `json:"value"`
	}{Type: BoolType, Value: b.Value}, nil
}

func decodeBool(s map[string]interface{}) (Node, error) {
	r, ok := s["value"].(bool)
	if !ok {
		return nil, fmt.Errorf("decoded value not Bool")
	}
	return Bool{r}, nil
}

// ToIPLD converts xr.Node into its corresponding IPLD Node type
func (b Bool) ToIPLD() (ipld.Node, error) {
	t := xrIpld.Type.Bool_IPLD.NewBuilder()
	err := t.AssignBool(b.Value)
	if err != nil {
		return nil, err
	}
	return t.Build(), nil
}

// toNode_IPLD convert into IPLD Node of dynamic type NODE_IPLD
func (b Bool) toNode_IPLD() (ipld.Node, error) {
	t := xrIpld.Type.Node_IPLD.NewBuilder()
	ma, err := t.BeginMap(-1)
	asm, err := ma.AssembleEntry("Bool_IPLD")
	if err != nil {
		return nil, err
	}
	err = asm.AssignBool(b.Value)
	if err != nil {
		return nil, err
	}
	if err := ma.Finish(); err != nil {
		return nil, err
	}
	return t.Build(), nil
}
