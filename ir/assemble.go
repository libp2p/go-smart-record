package ir

import "fmt"

// AssemblerContext holds general contextual data for the stage of the assembly process.
// It provides a standard mechanism for assemblers to pass context to subordinate assemblers
// that are called recursively.
// NOTE: The right general long-term design for AssemblerContext is to make it an interface.
// This is currently not necassitated by our uses, so such improvements are deferred for when needed.
type AssemblerContext struct {
	Grammar Assembler
	Keys    map[string]interface{}
}

// Assembler is an object that can "parse" a syntactic tree (given as a dictionary)
// and produces a semantic representation of what is parsed.
// Usually what is produced will be a smart tag.
// However, an Assembler is not restricted to produce semantic nodes (like smart tags).
// It can also produce syntactic nodes (like dictionaries).
// In that sense, an Assembler can be used to implement any transformation or
// even just a verification operation that returns the input unchanged.
type Assembler interface {
	Assemble(ctx AssemblerContext, src Dict) (Node, error)
}

// SequenceAssembler is, in common parlance, a parser combinator. Or, in our nomenclature, an "assembler combinator".
// SequenceAssembler tries to assemble the input, using each of its subordinate assemblers in turn until one of them succeeds.
type SequenceAssembler []Assembler

func (asm SequenceAssembler) Assemble(ctx AssemblerContext, src Dict) (Node, error) {
	for _, a := range asm {
		out, err := a.Assemble(ctx, src)
		if err == nil {
			return out, nil
		}
	}
	return nil, fmt.Errorf("no assembler in the sequence recognized the input")
}
