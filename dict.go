package sr

import (
	"io"
)

// Dict is a set of uniquely-named child nodes.
type Dict struct {
	Pairs []Pair // maintain: keys are unique
}

type Pair struct {
	Key   Node
	Value Node
}

func (p Pair) WritePretty(w io.Writer) error {
	if err := p.Key.WritePretty(w); err != nil {
		return err
	}
	if _, err := w.Write([]byte(" : ")); err != nil {
		return err
	}
	if err := p.Value.WritePretty(IndentWriter(w)); err != nil {
		return err
	}
	return nil
}

func (d Dict) WritePretty(w io.Writer) error {
	XXX
}

func MergeDicts(x, y *Dict) Node {
	panic("XXX")
}
