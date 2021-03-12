package ir

import (
	"io"

	"github.com/ipfs/go-cid"
)

type Cid struct {
	// Cid string
	// User holds user fields.
	User Dict
}

func (c Cid) Dict() Dict {
	outP := Pairs{}
	for _, p := range c.User.Pairs {
		// If key is not equal to value can't parse CID or not string
		if p.Key != p.Value {
			// TODO: We need an error type
			p.Value = String{"Err"}
		}
		switch p.Key.(type) {
		case String:
		default:
			p.Value = String{"Err"}
		}
		// Decode CID
		_, err := cid.Decode(p.Key.(String).Value)
		if err != nil {
			p.Value = String{"Err"}
		}
		outP = append(outP, p)
	}
	return Dict{Tag: c.User.Tag, Pairs: outP}
}

func (c Cid) WritePretty(w io.Writer) error {
	return c.Dict().WritePretty(w)
}
