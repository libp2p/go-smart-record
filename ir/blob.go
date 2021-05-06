package ir

import (
	"fmt"

	"github.com/libp2p/go-smart-record/xr"
)

type Blob struct {
	Bytes       []byte
	metadataCtx *metadataContext
}

func (b Blob) Disassemble() xr.Node {
	return xr.Blob{Bytes: b.Bytes}
}

func (b Blob) Metadata() MetadataInfo {
	return b.metadataCtx.getMetadata()
}

func (b Blob) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	w, ok := with.(Blob)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-blob")
	}
	// Update metadata
	b.metadataCtx.update(w.metadataCtx)
	return w, nil
}
