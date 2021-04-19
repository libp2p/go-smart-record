package base

import (
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-smart-record/ir"
)

// Cid is a smart node, representing a valid CID.
type Cid struct {
	Cid cid.Cid

	// User holds user fields.
	User ir.Dict
}

func (c Cid) EncodeJSON() (interface{}, error) {
	return c.Disassemble().EncodeJSON()
}

func (c Cid) Disassemble() ir.Node {
	return c.User.CopySetTag("cid", ir.String{"cid"}, ir.String{c.Cid.String()})
}

func (c Cid) WritePretty(w io.Writer) error {
	return c.Disassemble().WritePretty(w)
}

func (c Cid) UpdateWith(ctx ir.UpdateContext, with ir.Node) (ir.Node, error) {
	w, ok := with.(Cid)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-cid")
	}
	return w, nil
}

type CidAssembler struct{}

func (CidAssembler) Assemble(ctx ir.AssemblerContext, srcNode ir.Node) (ir.Node, error) {
	d, ok := srcNode.(ir.Dict)
	if !ok {
		return nil, fmt.Errorf("expecting dict")
	}
	if d.Tag != "cid" {
		return nil, fmt.Errorf("expecting tag cid")
	}
	if v := d.Get(ir.String{"cid"}); v == nil {
		return nil, fmt.Errorf("missing cid field")
	} else {
		s, ok := v.(ir.String)
		if !ok {
			return nil, fmt.Errorf("cid is not a string")
		}
		x, err := cid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("cid does not parse (%v)", err)
		}
		u := d.Copy()
		u.Tag = ""
		u.Remove(ir.String{"cid"})
		return Cid{
			Cid:  x,
			User: u,
		}, nil
	}
}
