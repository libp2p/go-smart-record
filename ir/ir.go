package ir

import (
	"io"
)

type Node interface {
	WritePretty(w io.Writer) error
	//UnmarshalJSON(data []byte) error
	MarshalJSON() (b []byte, e error)
}
