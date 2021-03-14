package ir

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
)

// String is a node representing a string literal.
type String struct {
	Value string
}

func (s String) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%q", s.Value)
	return err
}

func IsEqualString(x, y String) bool {
	return x.Value == y.Value
}

type Int64 struct {
	Value int64
}

func (s Int64) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%v", s.Value)
	return err
}

func IsEqualInt64(x, y Int64) bool {
	return x.Value == y.Value
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

func IsEqualNumber(x, y Number) bool {
	switch {
	case x.Int != nil && y.Int != nil:
		return x.Int.Cmp(y.Int) == 0
	case x.Float != nil && y.Float != nil:
		return x.Float.Cmp(y.Float) == 0
	case x.Rat != nil && y.Rat != nil:
		return x.Rat.Cmp(y.Rat) == 0
	}
	return false
}

type Blob struct {
	Bytes []byte
}

func (b Blob) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "0x%s", hex.EncodeToString(b.Bytes))
	return err
}

func IsEqualBlob(x, y Blob) bool {
	return bytes.Compare(x.Bytes, y.Bytes) == 0
}
