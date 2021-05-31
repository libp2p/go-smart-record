package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
)

// String is a node representing a string literal.
type String struct {
	Value       string
	metadataCtx *metadataContext
}

func (s *String) Disassemble() xr.Node {
	return xr.String{Value: s.Value}
}

func (s *String) Metadata() MetadataInfo {
	return s.metadataCtx.getMetadata()
}

func (s *String) UpdateWith(ctx UpdateContext, with Node) error {
	w, ok := with.(*String)
	if !ok {
		return fmt.Errorf("cannot update with a non-string")
	}
	// Update metadata
	s.metadataCtx.update(w.metadataCtx)
	return nil
}
