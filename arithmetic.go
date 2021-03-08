package sr

import (
	"io"
	"math/big"
)

type String struct {
	Value string
}

func (s String) WritePretty(w io.Writer, level int) error {
	panic("XXX")
}

type Number struct {
	*big.Int
	*big.Float
	*big.Rat
}

func (n Number) WritePretty(w io.Writer, level int) error {
	panic("XXX")
}

type Blob struct {
	Bytes []byte
}

func (b Blob) WritePretty(w io.Writer, level int) error {
	panic("XXX")
}
