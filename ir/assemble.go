package ir

import (
	"fmt"

	"github.com/libp2p/go-smart-record/xr"
)

// AssemblerContext holds general contextual data for the stage of the assembly process.
// It provides a standard mechanism for assemblers to pass context to subordinate assemblers
// that are called recursively.
// NOTE: The right general long-term design for AssemblerContext is to make it an interface.
// This is currently not necassitated by our uses, so such improvements are deferred for when needed.
type AssemblerContext struct {
	Grammar Assembler
	Keys    map[string]interface{}
}

func (ctx AssemblerContext) Assemble(src xr.Node) (Node, error) {
	return ctx.Grammar.Assemble(ctx, src)
}

// Assembler is an object that can "parse" a syntactic tree (given as a dictionary)
// and produces a semantic representation of what is parsed.
// Usually what is produced will be a smart tag.
// However, an Assembler is not restricted to produce semantic nodes (like smart tags).
// It can also produce syntactic nodes (like dictionaries).
// In that sense, an Assembler can be used to implement any transformation or
// even just a verification operation that returns the input unchanged.
type Assembler interface {
	Assemble(ctx AssemblerContext, src xr.Node) (Node, error)
}

// SequenceAssembler is, in common parlance, a parser combinator. Or, in our nomenclature, an "assembler combinator".
// SequenceAssembler tries to assemble the input, using each of its subordinate assemblers in turn until one of them succeeds.
type SequenceAssembler []Assembler

func (asm SequenceAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	for _, a := range asm {
		out, err := a.Assemble(ctx, src)
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
	BlobAssembler{},
	DictAssembler{},
	SetAssembler{},
}

type StringAssembler struct{}

func (asm StringAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	s, ok := src.(xr.String)
	if !ok {
		return nil, fmt.Errorf("not a string")
	}
	return String{Value: s.Value}, nil
}

type IntAssembler struct{}

func (asm IntAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	s, ok := src.(xr.Int)
	if !ok {
		return nil, fmt.Errorf("not an int")
	}
	return Int{Int: s.Int}, nil
}

type FloatAssembler struct{}

func (asm FloatAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	s, ok := src.(xr.Float)
	if !ok {
		return nil, fmt.Errorf("not a float")
	}
	return Float{Float: s.Float}, nil
}

type BoolAssembler struct{}

func (asm BoolAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	s, ok := src.(xr.Bool)
	if !ok {
		return nil, fmt.Errorf("not a bool")
	}
	return Bool{Value: s.Value}, nil
}

type BlobAssembler struct{}

func (asm BlobAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	s, ok := src.(xr.Blob)
	if !ok {
		return nil, fmt.Errorf("not a blob")
	}
	return Blob{Bytes: s.Bytes}, nil
}

type DictAssembler struct{}

func (asm DictAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	s, ok := src.(xr.Dict)
	if !ok {
		return nil, fmt.Errorf("not a dict")
	}
	d := Dict{
		Tag:   s.Tag,
		Pairs: make(Pairs, len(s.Pairs)),
	}
	for i, p := range s.Pairs {
		k, err := ctx.Assemble(p.Key)
		if err != nil {
			return nil, fmt.Errorf("key assembly (%v)", err)
		}
		v, err := ctx.Assemble(p.Value)
		if err != nil {
			return nil, fmt.Errorf("value assembly (%v)", err)
		}
		d.Pairs[i] = Pair{Key: k, Value: v}
	}
	return d, nil
}

type SetAssembler struct{}

func (asm SetAssembler) Assemble(ctx AssemblerContext, src xr.Node) (Node, error) {
	s, ok := src.(xr.Set)
	if !ok {
		return nil, fmt.Errorf("not a set")
	}
	d := Set{
		Tag:      s.Tag,
		Elements: make(Nodes, len(s.Elements)),
	}
	for i, e := range s.Elements {
		ae, err := ctx.Assemble(e)
		if err != nil {
			return nil, fmt.Errorf("element assembly (%v)", err)
		}
		d.Elements[i] = ae
	}
	return d, nil
}
