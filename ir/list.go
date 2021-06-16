package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

// List is a List of (uniquely) elements.
type List struct {
	Elements    Nodes
	metadataCtx *meta.Meta
}

func (s *List) Disassemble() xr.Node {
	x := xr.List{Elements: make(xr.Nodes, len(s.Elements))}
	for i, e := range s.Elements {
		x.Elements[i] = e.Disassemble()
	}
	return x
}

func (s *List) Metadata() meta.MetadataInfo {
	return s.metadataCtx.Get()
}

// MergeElements returns the union of the two Lists
func MergeElements(x, y Nodes) Nodes {
	z := make(Nodes, len(x), len(x)+len(y))
	copy(z, x)
	for _, el := range y {
		if i := x.IndexOf(el); i < 0 {
			z = append(z, el)
		}
	}
	return z
}

func (s List) Copy() List {
	e := make(Nodes, len(s.Elements))
	copy(e, s.Elements)
	// Also copy metadata if it exists
	m := s.metadataCtx.Copy()
	return List{
		Elements:    e,
		metadataCtx: &m,
	}
}

func (s List) Len() int {
	return len(s.Elements)
}

func (s *List) UpdateWith(ctx UpdateContext, with Node) error {
	ws, ok := with.(*List)
	if !ok {
		return fmt.Errorf("cannot update with a non-List")
	}
	for _, e := range ws.Elements {
		if i := s.Elements.IndexOf(e); i < 0 {
			s.Elements = append(s.Elements, e)
		}
	}
	// Update metadata
	s.metadataCtx.Update(ws.metadataCtx)
	return nil
}
