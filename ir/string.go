package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

// String is a node representing a string literal.
type String struct {
	Value       string
	metadataCtx *meta.Meta
}

func (s *String) Disassemble() xr.Node {
	return xr.String{Value: s.Value}
}

func (s *String) Metadata() meta.MetadataInfo {
	return s.metadataCtx.Get()
}

func (s *String) UpdateWith(ctx UpdateContext, with Node) error {
	w, ok := with.(*String)
	if !ok {
		return fmt.Errorf("cannot update with a non-string")
	}
	// Update value
	*s = *w
	// Update metadata
	s.metadataCtx.Update(w.metadataCtx)
	return nil
}
