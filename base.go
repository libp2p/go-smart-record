package sr

import (
	"io"
)

type Node interface {
	WritePretty(w io.Writer, level int) error
}
