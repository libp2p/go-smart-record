package ir

import (
	"fmt"

	"github.com/libp2p/go-smart-record/xr"
)

// Set is a set of (uniquely) elements.
type Set struct {
	Tag         string
	Elements    Nodes
	metadataCtx *metadataContext
}

func (s *Set) Disassemble() xr.Node {
	x := xr.Set{Tag: s.Tag, Elements: make(xr.Nodes, len(s.Elements))}
	for i, e := range s.Elements {
		x.Elements[i] = e.Disassemble()
	}
	return x
}

func (s *Set) Metadata() MetadataInfo {
	return s.metadataCtx.getMetadata()
}

func (s Set) Copy() Set {
	e := make(Nodes, len(s.Elements))
	copy(e, s.Elements)
	// Also copy metadata if it exists
	m := s.metadataCtx.copy()
	return Set{
		Tag:         s.Tag,
		Elements:    e,
		metadataCtx: &m,
	}
}

func (s Set) Len() int {
	return len(s.Elements)
}

func (s *Set) UpdateWith(ctx UpdateContext, with Node) error {
	ws, ok := with.(*Set)
	if !ok {
		return fmt.Errorf("cannot update with a non-set")
	}
	s.Tag = ws.Tag
	for _, e := range ws.Elements {
		if i := s.Elements.IndexOf(e); i < 0 {
			s.Elements = append(s.Elements, e)
		}
	}
	// Update metadata
	s.metadataCtx.update(ws.metadataCtx)
	return nil
}
