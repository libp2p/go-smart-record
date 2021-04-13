package base

import (
	"io"

	"github.com/libp2p/go-smart-record/ir"
)

type Verify struct {
	Statement ir.Node
	// User holds user fields.
	User ir.Dict
}

func (v Verify) EncodeJSON() (interface{}, error) {
	return v.Disassemble().EncodeJSON()
}

func (v Verify) Disassemble() ir.Dict {
	return v.User.CopySetTag("verify", ir.String{"statement"}, v.Statement)
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

func (v Verified) Disassemble() ir.Dict {
	return ir.Dict{
		Tag: "verify",
		Pairs: ir.MergePairsRight(
			v.User.Pairs,
			ir.Pairs{
				{Key: ir.String{"by"}, Value: v.By},
				{Key: ir.String{"statement"}, Value: v.Statement},
				{Key: ir.String{"signature"}, Value: v.Signature},
			},
		),
	}
}

func (v Verified) WritePretty(w io.Writer) error {
	return v.Disassemble().WritePretty(w)
}
