package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

type Bool struct {
	Value       bool
	metadataCtx *meta.Meta
}

func (b *Bool) Disassemble() xr.Node {
	return xr.Bool{Value: b.Value}
}

func (b *Bool) Metadata() meta.MetadataInfo {
	return b.metadataCtx.GetMeta()
}

func (b *Bool) UpdateWith(ctx UpdateContext, with Node) error {
	w, ok := with.(*Bool)
	if !ok {
		return fmt.Errorf("cannot update with a non-bool")
	}
	// Update metadata
	b.metadataCtx.Update(w.metadataCtx)
	return nil
}
