package sr

import (
	"io"
)

type Cid struct {
	Cid string
	Dict
}

func (c Cid) WritePretty(w io.Writer) error {
	panic("XXX")
}
