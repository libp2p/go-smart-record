package ir

import (
	"fmt"

	"github.com/libp2p/go-smart-record/xr"
)

type Bool struct {
	Value       bool
	metadataCtx *metadataContext
}

func (b *Bool) Disassemble() xr.Node {
	return xr.Bool{Value: b.Value}
}

func (b *Bool) Metadata() MetadataInfo {
	return b.metadataCtx.getMetadata()
}

func (b *Bool) UpdateWith(ctx UpdateContext, with Node) error {
	w, ok := with.(*Bool)
	if !ok {
		return fmt.Errorf("cannot update with a non-bool")
	}
	// Update metadata
	b.metadataCtx.update(w.metadataCtx)
	return nil
}
