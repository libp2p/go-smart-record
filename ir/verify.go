package ir

import (
	"io"
)

type Verify struct {
	Statement Node
	// User holds user fields.
	User Dict
}

func (v Verify) Disassemble() Dict {
	return v.User.CopySetTag("verify", String{"statement"}, v.Statement)
}

func (v Verify) WritePretty(w io.Writer) error {
	return v.Disassemble().WritePretty(w)
}

type Verified struct {
	By        Peer
	Statement Node
	Signature Blob
	// User holds user fields.
	User Dict
}

func (v Verified) Disassemble() Dict {
	return Dict{
		Tag: "verify",
		Pairs: MergePairsRight(
			v.User.Pairs,
			Pairs{
				{Key: String{"by"}, Value: v.By},
				{Key: String{"statement"}, Value: v.Statement},
				{Key: String{"signature"}, Value: v.Signature},
			},
		),
	}
}

func (v Verified) WritePretty(w io.Writer) error {
	return v.Disassemble().WritePretty(w)
}
