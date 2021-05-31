package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
)

// Predicate models a function invocation with named and positional arguments, corresponding to the syntax:
//   tag(a1, a2, ...; n1=v1, n2=v2, ...)
type Predicate struct {
	Tag         string
	Positional  Nodes
	Named       Pairs // the keys in each pair must be unique wrt IsEqual
	metadataCtx *metadataContext
}

func (p *Predicate) Disassemble() xr.Node {
	x := xr.Predicate{
		Tag:        p.Tag,
		Positional: make(xr.Nodes, len(p.Positional)),
		Named:      make(xr.Pairs, len(p.Named)),
	}
	for i, e := range p.Positional {
		x.Positional[i] = e.Disassemble()
	}
	for i, p := range p.Named {
		x.Named[i] = xr.Pair{Key: p.Key.Disassemble(), Value: p.Value.Disassemble()}
	}
	return x
}

func (p *Predicate) Metadata() MetadataInfo {
	return p.metadataCtx.getMetadata()
}

func (p *Predicate) UpdateWith(ctx UpdateContext, with Node) error {
	w, ok := with.(*Predicate)
	if !ok {
		return fmt.Errorf("cannot update with a non-predicate")
	}
	// Update metadata
	p.metadataCtx.update(w.metadataCtx)
	return nil
}
