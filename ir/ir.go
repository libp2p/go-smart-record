package ir

import (
	"io"
)

type Node interface {
	WritePretty(w io.Writer) error
}
