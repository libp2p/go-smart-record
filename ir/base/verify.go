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
	User *ir.Dict
}

func (v *Verify) Disassemble() xr.Node {
	return v.User.Disassemble().(xr.Dict).CopySetTag("verify",
		xr.String{"statement"}, v.Statement.Disassemble())
}

type Verified struct {
	By        *Peer
	Statement ir.Node
	Signature *ir.Blob
	// User holds user fields.
	User *ir.Dict
}

func (v *Verified) Disassemble() xr.Node {
	return (&ir.Dict{
		Tag: "verify",
		Pairs: ir.MergePairs(
			v.User.Pairs,
			ir.Pairs{
				{Key: &ir.String{Value: "by"}, Value: v.By},
				{Key: &ir.String{Value: "statement"}, Value: v.Statement},
				{Key: &ir.String{Value: "signature"}, Value: v.Signature},
			},
		),
	}).Disassemble()
}

func (v *Verified) Metadata() ir.MetadataInfo {
	return v.User.Metadata()
}
func (v *Verified) WritePretty(w io.Writer) error {
	return v.Disassemble().WritePretty(w)
}

func (v *Verified) UpdateWith(ctx ir.UpdateContext, with ir.Node) error {
	_, ok := with.(*Signed)
	if !ok {
		return fmt.Errorf("cannot update with a non-verified")
	}
	return nil
}
