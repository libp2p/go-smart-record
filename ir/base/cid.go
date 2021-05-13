package base

import (
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
)

// Cid is a smart node, representing a valid CID.
type Cid struct {
	Cid cid.Cid

	// User holds user fields.
	User *ir.Dict
}

func (c *Cid) Disassemble() xr.Node {
	return c.User.Disassemble().(xr.Dict).CopySetTag("cid",
		xr.String{"cid"}, xr.String{c.Cid.String()})
}

func (c *Cid) Metadata() ir.MetadataInfo {
	return c.User.Metadata()
}

func (c *Cid) UpdateWith(ctx ir.UpdateContext, with ir.Node) error {
	_, ok := with.(*Cid)
	if !ok {
		return fmt.Errorf("cannot update with a non-cid")
	}
	return nil
}

type CidAssembler struct{}

func (CidAssembler) Assemble(ctx ir.AssemblerContext, srcNode xr.Node, metadata ...ir.Metadata) (ir.Node, error) {
	d, ok := srcNode.(xr.Dict)
	if !ok {
		return nil, fmt.Errorf("expecting dict")
	}
	if d.Tag != "cid" {
		return nil, fmt.Errorf("expecting tag cid")
	}
	if v := d.Get(xr.String{"cid"}); v == nil {
		return nil, fmt.Errorf("missing cid field")
	} else {
		s, ok := v.(xr.String)
		if !ok {
			return nil, fmt.Errorf("cid is not a string")
		}
		x, err := cid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("cid does not parse (%v)", err)
		}
		u := d.Copy()
		u.Tag = ""
		u.Remove(xr.String{"cid"})

		asm := ir.DictAssembler{}
		uasm, err := asm.Assemble(ctx, d, metadata...)
		return &Cid{
			Cid:  x,
			User: uasm.(*ir.Dict),
		}, nil
	}
}
