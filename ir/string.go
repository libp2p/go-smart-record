package ir

import (
	"fmt"

	"github.com/libp2p/go-smart-record/xr"
)

// String is a node representing a string literal.
type String struct {
	Value string
}

func (s String) Disassemble() xr.Node {
	return xr.String{Value: s.Value}
}

func (s String) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	w, ok := with.(String)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-string")
	}
	return w, nil
}
