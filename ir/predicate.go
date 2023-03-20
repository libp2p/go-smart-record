package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

// Predicate models a function invocation with named and positional arguments, corresponding to the syntax:
//
//	tag(a1, a2, ...; n1=v1, n2=v2, ...)
type Predicate struct {
	Tag         string
	Positional  Nodes
	Named       Pairs // the keys in each pair must be unique wrt IsEqual
	metadataCtx *meta.Meta
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

func (p *Predicate) Metadata() meta.MetadataInfo {
	return p.metadataCtx.Get()
}

func (p *Predicate) GetNamed(key Node) Node {
	for _, ps := range p.Named {
		if IsEqual(ps.Key, key) {
			return ps.Value
		}
	}
	return nil
}

func (p *Predicate) UpdateWith(ctx UpdateContext, with Node) error {
	wp, ok := with.(*Predicate)
	if !ok {
		return fmt.Errorf("cannot update with a non-predicate")
	}

	// Check equal tag
	if wp.Tag != p.Tag {
		return fmt.Errorf("predicate tags are not equal")
	}

	// Update positional
	for _, e := range wp.Positional {
		if i := p.Positional.IndexOf(e); i < 0 {
			p.Positional = append(p.Positional, e)
		}
	}

	// Update named
	for _, ps := range wp.Named {
		if i := p.Named.IndexOf(ps.Key); i < 0 {
			p.Named = append(p.Named, ps)
		} else {
			if err := p.Named[i].Value.UpdateWith(ctx, ps.Value); err != nil {
				return fmt.Errorf("cannout update value (%v)", err)
			}
		}
	}

	// Update metadata
	p.metadataCtx.Update(wp.metadataCtx)
	return nil
}
