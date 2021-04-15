package base

import (
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

type CidAssembler struct{}

func (CidAssembler) Assemble(ctx ir.AssemblerContext, srcNode ir.Node) (ir.Node, error) {
	panic("XXX")
}
