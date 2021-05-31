package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
)

type Bytes struct {
	Bytes       []byte
	metadataCtx *metadataContext
}

func (b *Bytes) Disassemble() xr.Node {
	return xr.Bytes{Bytes: b.Bytes}
}

func (b *Bytes) Metadata() MetadataInfo {
	return b.metadataCtx.getMetadata()
}

func (b *Bytes) UpdateWith(ctx UpdateContext, with Node) error {
	w, ok := with.(*Bytes)
	if !ok {
		return fmt.Errorf("cannot update with a non-Bytes")
	}
	// Update metadata
	b.metadataCtx.update(w.metadataCtx)
	return nil
}
