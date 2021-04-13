package base

import (
	"io"

	"github.com/libp2p/go-smart-record/ir"
)

type Signed struct {
	By        Peer
	Statement ir.Node
	Signature ir.Blob
	// User holds user fields.
	User ir.Dict
}

func (s Signed) WritePretty(w io.Writer) error {
	return s.Disassemble().WritePretty(w)
}

func (s Signed) EncodeJSON() (interface{}, error) {
	return s.Disassemble().EncodeJSON()
}

func (s Signed) Disassemble() ir.Dict {
	return ir.Dict{
		Tag: "verify",
		Pairs: ir.MergePairsRight(
			s.User.Pairs,
			ir.Pairs{
				{Key: ir.String{"by"}, Value: s.By},
				{Key: ir.String{"statement"}, Value: s.Statement},
				{Key: ir.String{"signature"}, Value: s.Signature},
			},
		),
	}
}

func (s Signed) MergeWith(ctx ir.MergeContext, x ir.Node) (ir.Node, error) {
	panic("XXX")
}
