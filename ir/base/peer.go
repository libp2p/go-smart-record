package base

import (
	"fmt"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
)

type Peer struct {
	ID string
	// User holds user fields.
	User *ir.Dict
}

func (p *Peer) Disassemble() xr.Node {
	return p.User.Disassemble().(xr.Dict).CopySetTag("peer",
		xr.String{"id"}, xr.String{p.ID})
}

func (p *Peer) Metadata() ir.MetadataInfo {
	return p.User.Metadata()
}

func (p *Peer) UpdateWith(ctx ir.UpdateContext, with ir.Node) error {
	_, ok := with.(*Peer)
	if !ok {
		return fmt.Errorf("cannot update with a non-peer")
	}
	return nil
}
