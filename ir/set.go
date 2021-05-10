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

func (s Set) Disassemble() xr.Node {
	x := xr.Set{Tag: s.Tag, Elements: make(xr.Nodes, len(s.Elements))}
	for i, e := range s.Elements {
		x.Elements[i] = e.Disassemble()
	}
	return x
}

func (s Set) Metadata() MetadataInfo {
	return s.metadataCtx.getMetadata()
}

// MergeElements returns the union of the two sets
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

func (s Set) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	ws, ok := with.(Set)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-set")
	}
	u := s.Copy()
	u.Tag = ws.Tag
	for _, e := range ws.Elements {
		if i := u.Elements.IndexOf(e); i < 0 {
			u.Elements = append(u.Elements, e)
		}
	}
	// Update metadata
	u.metadataCtx.update(ws.metadataCtx)
	return u, nil
}
