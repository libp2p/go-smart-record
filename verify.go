package sr

import (
	"io"
)

type Verify struct {
	Stmt Node
	Dict
}

func (v Verify) WritePretty(w io.Writer) error {
	panic("XXX")
}

type Verified struct {
	By   Peer
	Sign Blob
	Stmt Node
	Dict
}

func (v Verified) WritePretty(w io.Writer) error {
	panic("XXX")
}
