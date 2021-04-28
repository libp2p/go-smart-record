package ir

import (
	"fmt"

	"github.com/libp2p/go-smart-record/xr"
)

type Blob struct {
	Bytes []byte
}

func (b Blob) Disassemble() xr.Node {
	return xr.Blob{Bytes: b.Bytes}
}

func (b Blob) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	w, ok := with.(Blob)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-blob")
	}
	return w, nil
}
