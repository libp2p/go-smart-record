package sr

import (
	"io"
)

// Dict is a set of uniquely-named child nodes.
type Dict struct {
	Pairs []Pair // keys are unique
}

type Pair struct {
	Key   Node
	Value Node
}

func (d Dict) WritePretty(w io.Writer, level int) error {
	panic("XXX")
}

func MergeDicts(x, y *Dict) Node {
	panic("XXX")
}
