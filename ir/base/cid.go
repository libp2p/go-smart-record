package base

import (
	"io"

	"github.com/libp2p/go-smart-record/ir"
)

// Cid is a smart node, representing a valid CID.
type Cid struct {
	Cid string // TODO: This should be of type cid.Cid

	// User holds user fields.
	User ir.Dict
}

func (c Cid) EncodeJSON() (interface{}, error) {
	return c.Disassemble().EncodeJSON()
}

func (c Cid) Disassemble() ir.Dict {
	return c.User.CopySetTag("cid", ir.String{c.Cid}, ir.String{c.Cid})
}

func (c Cid) WritePretty(w io.Writer) error {
	return c.Disassemble().WritePretty(w)
}

func (c Cid) MergeWith(ctx ir.MergeContext, x ir.Node) (ir.Node, error) {
	panic("XXX")
}
