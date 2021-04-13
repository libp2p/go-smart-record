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

func (c Cid) Disassemble() ir.Dict {
	return c.User.CopySetTag("cid", ir.String{"cid"}, ir.String{c.Cid.String()})
}

func (c Cid) WritePretty(w io.Writer) error {
	return c.Disassemble().WritePretty(w)
}

func (c Cid) MergeWith(ctx ir.MergeContext, x ir.Node) (ir.Node, error) {
	xc, ok := x.(Cid)
	if !ok {
		return nil, fmt.Errorf("cannot merge cid with non-cid")
	}
	if !c.Cid.Equals(xc.Cid) {
		return nil, fmt.Errorf("cannot merge unequal cids")
	}
	u, err := ir.Merge(ctx, c.User, xc.User)
	if err != nil {
		return nil, fmt.Errorf("cannot merge cid user data")
	}
	return Cid{
		Cid:  c.Cid,
		User: u.(ir.Dict),
	}, nil
}

type CidAssembler struct{}

func (CidAssembler) Assemble(ctx ir.AssemblerContext, srcNode ir.Node) (ir.Node, error) {
	panic("XXX")
}
