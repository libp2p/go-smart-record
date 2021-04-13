package base

import (
	"io"

	"github.com/libp2p/go-smart-record/ir"
)

type Peer struct {
	ID string
	// User holds user fields.
	User ir.Dict
}

func (p Peer) EncodeJSON() (interface{}, error) {
	return p.Disassemble().EncodeJSON()
}

func (p Peer) Disassemble() ir.Dict {
	return p.User.CopySetTag("peer", ir.String{"id"}, ir.String{p.ID})
}

func (p Peer) WritePretty(w io.Writer) error {
	return p.Disassemble().WritePretty(w)
}

func (p Peer) MergeWith(ctx ir.MergeContext, x ir.Node) (ir.Node, error) {
	panic("XXX")
}

func (p Peer) Encoding() ir.Encoder {
	return nil
}
