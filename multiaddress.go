package sr

import (
	"io"
)

type Multiaddress struct {
	Multiaddress string
	Dict
}

func (m Multiaddress) WritePretty(w io.Writer) error {
	panic("XXX")
}
