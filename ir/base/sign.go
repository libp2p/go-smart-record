package base

import (
	"fmt"
	"io"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
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

func (s Signed) Disassemble() xr.Node {
	return ir.Dict{
		Tag: "verify",
		Pairs: ir.MergePairs(
			s.User.Pairs,
			ir.Pairs{
				{Key: ir.String{Value: "by"}, Value: s.By},
				{Key: ir.String{Value: "statement"}, Value: s.Statement},
				{Key: ir.String{Value: "signature"}, Value: s.Signature},
			},
		),
	}.Disassemble()
}

func (s Signed) Metadata() ir.MetadataInfo {
	return s.User.Metadata()
}

func (s Signed) UpdateWith(ctx ir.UpdateContext, with ir.Node) (ir.Node, error) {
	w, ok := with.(Signed)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-signed")
	}
	return w, nil
}
