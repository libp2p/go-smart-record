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

func (s Set) Copy() Set {
	e := make(Nodes, len(s.Elements))
	copy(e, s.Elements)
	return Set{
		Tag:         s.Tag,
		Elements:    e,
		metadataCtx: s.metadataCtx, // Also copy metadata.
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
