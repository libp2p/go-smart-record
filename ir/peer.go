package ir

import (
	"io"
)

type Peer struct {
	ID string
	// User holds user fields.
	User Dict
}

func (p Peer) Disassemble() Dict {
	return p.User.CopySetTag("peer", String{"id"}, String{p.ID})
}

func (p Peer) WritePretty(w io.Writer) error {
	return p.Disassemble().WritePretty(w)
}

func (p Peer) MergeWith(ctx MergeContext, x Node) Node {
	panic("XXX")
}
