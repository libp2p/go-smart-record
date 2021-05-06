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

func (p Peer) Disassemble() xr.Node {
	return p.User.Disassemble().(xr.Dict).CopySetTag("peer",
		xr.String{"id"}, xr.String{p.ID})
}

func (p Peer) Metadata() ir.MetadataInfo {
	return p.User.Metadata()
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
