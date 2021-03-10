package ir

import (
	"io"
)

type Peer struct {
	ID string
	// User holds user fields.
	User Dict
}

func (p Peer) Dict() Dict {
	return p.User.CopySetTag("peer", String{"id"}, String{p.ID})
}

func (p Peer) WritePretty(w io.Writer) error {
	return p.Dict().WritePretty(w)
}
