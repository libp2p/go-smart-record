package sr

import (
	"io"
)

type Cid struct {
	Cid string
	Dict
}

func (c Cid) WritePretty(w io.Writer, level int) error {
	panic("XXX")
}
