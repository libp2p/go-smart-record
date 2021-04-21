package base

import (
	"fmt"
	"io"

	"github.com/libp2p/go-smart-record/ir"
)

//NOTE: Including everything in the same file for discussion purposes.
// This should be organized in its corresponding files when implementing it.
type SmartNode interface {
	ir.Node
	Disassemble() ir.Node // returns only syntactic nodes
}

// SmartString is the smart node for the String type.
// The difference is that additional metadata is included in the assembly process
type SmartString struct {
	Value    ir.String
	Metadata metadataContext
}

// List of metadata attributes supported.
type metadataContext struct {
	ttl uint64
}

// TTL sets in metadata
func TTL(value uint64) Metadata {
	return func(m *metadataContext) error {
		m.ttl = value
		return nil
	}
}

// Option type for smart records
type Metadata func(*metadataContext) error

// Applies metadata items to a metadataContext
func (m *metadataContext) apply(items ...Metadata) error {
	for i, item := range items {
		if err := item(m); err != nil {
			return fmt.Errorf("error applying metadata value: %s", i, err)
		}
	}
	return nil
}

// The Assembler includes metadata to the syntactic node and assembles a smartNode.
// SmartNode includes additional metadata
// This belongs to (asm SequenceAssembler)
func Assemble(ctx ir.AssemblerContext, src ir.Node, metadata ...Metadata) (SmartNode, error) {
	s, ok := src.(ir.String)
	if !ok {
		return nil, fmt.Errorf("not a string")
	}
	var m metadataContext

	if err := m.apply(metadata...); err != nil {
		// Assembly fails if the wrong metadata is passed to the context.
		return nil, err
	}
	return SmartString{s, m}, nil
}

func (s SmartString) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%q", s.Value)
	return err
}

func (s SmartString) EncodeJSON() (interface{}, error) {
	return s.Value.Disassemble().EncodeJSON()
}

// Disassemble returns the ir.String value without metadata
func (s SmartString) Disassemble() ir.Node {
	return s.Value
}

// UpdateWith updates metadataContext with new info.
func (s SmartString) UpdateWith(ctx ir.UpdateContext, with ir.Node, metadata ...Metadata) (SmartNode, error) {
	w, ok := with.(String)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-string")
	}
	// Update current metadata with updated values
	if err := s.m.apply(metadata...); err != nil {
		// Assembly fails if the wrong metadata is passed to the context.
		return nil, err
	}
	return w, nil
}
