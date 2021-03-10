package sr

import (
	"io"
)

type Node interface {
	WritePretty(w io.Writer) error
}

// Smart represents a "smart" node.
type Smart interface {
	Node
	AsDict() Dict
}
