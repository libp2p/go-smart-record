package sr

import (
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

type String struct {
	Value string
}

func (s String) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%q", s.Value)
	return err
}

type Number struct {
	*big.Int
	*big.Float
	*big.Rat
}

func (n Number) WritePretty(w io.Writer) (err error) {
	switch {
	case n.Int != nil:
		_, err = w.Write([]byte(n.Int.String()))
	case n.Float != nil:
		_, err = w.Write([]byte(n.Float.String()))
	case n.Rat != nil:
		_, err = w.Write([]byte(n.Rat.String()))
	}
	return err
}

type Blob struct {
	Bytes []byte
}

func (b Blob) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "0x%s", hex.EncodeToString(b.Bytes))
	return err
}
