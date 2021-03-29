package ir

import (
	"fmt"
	"io"
)

type Bool struct {
	Value bool
}

func (b Bool) WritePretty(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, "%v", b.Value)
	return err
}

func IsEqualBool(x, y Bool) bool {
	return x.Value == y.Value
}
