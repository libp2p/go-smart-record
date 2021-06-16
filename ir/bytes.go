package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

type Bytes struct {
	Bytes       []byte
	metadataCtx *meta.Meta
}

func (b *Bytes) Disassemble() xr.Node {
	return xr.Bytes{Bytes: b.Bytes}
}

func (b *Bytes) Metadata() meta.MetadataInfo {
	return b.metadataCtx.Get()
}

func (b *Bytes) UpdateWith(ctx UpdateContext, with Node) error {
	w, ok := with.(*Bytes)
	if !ok {
		return fmt.Errorf("cannot update with a non-Bytes")
	}
	// Update value
	*b = *w
	// Update metadata
	b.metadataCtx.Update(w.metadataCtx)
	return nil
}
