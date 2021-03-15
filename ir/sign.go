package ir

import (
	"io"
)

type Signed struct {
	By        Peer
	Statement Node
	Signature Blob
	// User holds user fields.
	User Dict
}

func (s Signed) WritePretty(w io.Writer) error {
	return s.Disassemble().WritePretty(w)
}

func (s Signed) Disassemble() Dict {
	return Dict{
		Tag: "verify",
		Pairs: MergePairsRight(
			s.User.Pairs,
			Pairs{
				{Key: String{"by"}, Value: s.By},
				{Key: String{"statement"}, Value: s.Statement},
				{Key: String{"signature"}, Value: s.Signature},
			},
		),
	}
}

func (s Signed) MergeWith(ctx MergeContext, x Node) Node {
	panic("XXX")
}
