package base

import (
	"fmt"
	"io"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
)

type Peer struct {
	ID string
	// User holds user fields.
	User ir.Dict
}

func (p Peer) EncodeJSON() (interface{}, error) {
	return p.Disassemble().EncodeJSON()
}

func (p Peer) Disassemble() xr.Node {
	return p.User.CopySetTag("peer", ir.String{"id"}, ir.String{p.ID}).Disassemble()
}

func (p Peer) WritePretty(w io.Writer) error {
	return p.Disassemble().WritePretty(w)
}

func (p Peer) UpdateWith(ctx ir.UpdateContext, with ir.Node) (ir.Node, error) {
	w, ok := with.(Peer)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-peer")
	}
	return w, nil
}
