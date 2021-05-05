package ir

import (
	"fmt"

	"github.com/libp2p/go-smart-record/xr"
)

type Bool struct {
	Value bool
}

func (b Bool) Disassemble() xr.Node {
	return xr.Bool{Value: b.Value}
}

func (b Bool) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	w, ok := with.(Bool)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-bool")
	}
	return w, nil
}
