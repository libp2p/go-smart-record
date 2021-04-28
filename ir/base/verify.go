package base

import (
	"fmt"
	"io"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
)

type Verify struct {
	Statement ir.Node
	// User holds user fields.
	User ir.Dict
}

func (v Verify) EncodeJSON() (interface{}, error) {
	return v.Disassemble().EncodeJSON()
}

func (v Verify) Disassemble() xr.Node {
	return v.User.CopySetTag("verify", ir.String{"statement"}, v.Statement).Disassemble()
}

func (v Verify) WritePretty(w io.Writer) error {
	return v.Disassemble().WritePretty(w)
}

type Verified struct {
	By        Peer
	Statement ir.Node
	Signature ir.Blob
	// User holds user fields.
	User ir.Dict
}

func (v Verified) EncodeJSON() (interface{}, error) {
	return v.Disassemble().EncodeJSON()
}

func (v Verified) Disassemble() xr.Node {
	return ir.Dict{
		Tag: "verify",
		Pairs: ir.MergePairs(
			v.User.Pairs,
			ir.Pairs{
				{Key: ir.String{"by"}, Value: v.By},
				{Key: ir.String{"statement"}, Value: v.Statement},
				{Key: ir.String{"signature"}, Value: v.Signature},
			},
		),
	}.Disassemble()
}

func (v Verified) WritePretty(w io.Writer) error {
	return v.Disassemble().WritePretty(w)
}

func (v Verified) UpdateWith(ctx ir.UpdateContext, with ir.Node) (ir.Node, error) {
	w, ok := with.(Signed)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-verified")
	}
	return w, nil
}
