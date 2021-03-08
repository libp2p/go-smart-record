package sr

import (
	"io"
)

type Peer struct {
	ID string
	Dict
}

func (p Peer) WritePretty(w io.Writer, level int) error {
	panic("XXX")
}
