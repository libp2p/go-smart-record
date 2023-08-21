package ir

import (
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	xr "github.com/libp2p/go-routing-language/syntax"

	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

// AssemblerContext holds general contextual data for the stage of the assembly process.
// It provides a standard mechanism for assemblers to pass context to subordinate assemblers
// that are called recursively.
// NOTE: The right general long-term design for AssemblerContext is to make it an interface.
// This is currently not necassitated by our uses, so such improvements are deferred for when needed.
type AssemblerContext struct {
	Grammar Assembler
	Keys    map[string]interface{}
	Host    host.Host
}

func (ctx AssemblerContext) Assemble(src xr.Node, metadata ...meta.Metadata) (Node, error) {
	return ctx.Grammar.Assemble(ctx, src, metadata...)
}

// Assembler is an object that can "parse" a syntactic tree (given as a dictionary)
// and produces a semantic representation of what is parsed.
// Usually what is produced will be a smart tag.
// However, an Assembler is not restricted to produce semantic nodes (like smart tags).
// It can also produce syntactic nodes (like dictionaries).
// In that sense, an Assembler can be used to implement any transformation or
// even just a verification operation that returns the input unchanged.
// Assemblers can add metadata to the generated semantic nodes.
type Assembler interface {
	Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error)
}

// SequenceAssembler is, in common parlance, a parser combinator. Or, in our nomenclature, an "assembler combinator".
// SequenceAssembler tries to assemble the input, using each of its subordinate assemblers in turn until one of them succeeds.
// NOTE: With the current implementation all assembled nodes are updated with the same metadata. In the future
// additional tags could be specified so nodes can be assembled with different metadata each.
type SequenceAssembler []Assembler

func (asm SequenceAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	for _, a := range asm {
		out, err := a.Assemble(ctx, src, metadata...)
		if err == nil {
			return out, nil
		}
	}
	return nil, fmt.Errorf("no assembler in the sequence recognized the input")
}

var SyntacticGrammar = SequenceAssembler{
	StringAssembler{},
	IntAssembler{},
	FloatAssembler{},
	BoolAssembler{},
	BytesAssembler{},
	DictAssembler{},
	ListAssembler{},
	PredicateAssembler{},
}

type StringAssembler struct{}

func (asm StringAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.String)
	if !ok {
		return nil, fmt.Errorf("not a string")
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}

	return &String{Value: s.Value, metadataCtx: m}, nil
}

type IntAssembler struct{}

func (asm IntAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.Int)
	if !ok {
		return nil, fmt.Errorf("not an int")
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}

	return &Int{Int: s.Int, metadataCtx: m}, nil
}

type FloatAssembler struct{}

func (asm FloatAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.Float)
	if !ok {
		return nil, fmt.Errorf("not a float")
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}
	return &Float{Float: s.Float, metadataCtx: m}, nil
}

type BoolAssembler struct{}

func (asm BoolAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.Bool)
	if !ok {
		return nil, fmt.Errorf("not a bool")
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}

	return &Bool{Value: s.Value, metadataCtx: m}, nil
}

type BytesAssembler struct{}

func (asm BytesAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.Bytes)
	if !ok {
		return nil, fmt.Errorf("not bytes type")
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}

	return &Bytes{Bytes: s.Bytes, metadataCtx: m}, nil
}

type DictAssembler struct{}

func (asm DictAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.Dict)
	if !ok {
		return nil, fmt.Errorf("not a dict")
	}
	d := Dict{
		Pairs: make(Pairs, len(s.Pairs)),
	}
	for i, p := range s.Pairs {
		k, err := ctx.Assemble(p.Key, metadata...)
		if err != nil {
			return nil, fmt.Errorf("key assembly (%v)", err)
		}
		v, err := ctx.Assemble(p.Value, metadata...)
		if err != nil {
			return nil, fmt.Errorf("value assembly (%v)", err)
		}
		d.Pairs[i] = Pair{Key: k, Value: v}
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}
	d.metadataCtx = m

	return &d, nil
}

type ListAssembler struct{}

func (asm ListAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.List)
	if !ok {
		return nil, fmt.Errorf("not a List")
	}
	d := List{
		Elements: make(Nodes, len(s.Elements)),
	}
	for i, e := range s.Elements {
		ae, err := ctx.Assemble(e, metadata...)
		if err != nil {
			return nil, fmt.Errorf("element assembly (%v)", err)
		}
		d.Elements[i] = ae
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}
	d.metadataCtx = m

	return &d, nil
}

type PredicateAssembler struct{}

func (asm PredicateAssembler) Assemble(ctx AssemblerContext, src xr.Node, metadata ...meta.Metadata) (Node, error) {
	s, ok := src.(xr.Predicate)
	if !ok {
		return nil, fmt.Errorf("not a predicate")
	}
	d := Predicate{
		Positional: make(Nodes, len(s.Positional)),
		Named:      make(Pairs, len(s.Named)),
	}
	// Assemble positional
	for i, e := range s.Positional {
		ae, err := ctx.Assemble(e, metadata...)
		if err != nil {
			return nil, fmt.Errorf("positional assembly (%v)", err)
		}
		d.Positional[i] = ae
	}

	// Assemble named
	for i, p := range s.Named {
		k, err := ctx.Assemble(p.Key, metadata...)
		if err != nil {
			return nil, fmt.Errorf("key assembly (%v)", err)
		}
		v, err := ctx.Assemble(p.Value, metadata...)
		if err != nil {
			return nil, fmt.Errorf("value assembly (%v)", err)
		}
		d.Named[i] = Pair{Key: k, Value: v}
	}

	// Assemble metadata
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return nil, err
	}
	d.metadataCtx = m

	return &d, nil
}
