// Package base implements the basic set of smart tags supported by a smart record.
package base

import "github.com/libp2p/go-smart-record/ir"

// BaseGrammar is an assembler for the base vocabulary of smart tags supported by a record.
var BaseGrammar = ir.SequenceAssembler{
	// insert the assemblers of smart tags here
	CidAssembler{},
	// if no smart tag parses the input, keep it as is (in the form of syntactic nodes)
	ir.SyntacticGrammar,
}
