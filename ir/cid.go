package ir

import (
	"fmt"
	"io"

	"github.com/ipfs/go-cid"
)

type Cid struct {
	Cid string
	// User holds user fields.
	User Dict
}

func (c Cid) Dict() Dict {
	return c.User.CopySetTag("cid", String{c.Cid}, String{c.Cid})
}

func (c Cid) WritePretty(w io.Writer) error {
	return c.Dict().WritePretty(w)
}

func ParseCidFromDict(d Dict) (Pairs, error) {
	outP := Pairs{}
	for _, p := range d.Pairs {
		// If key is not equal to value can't parse CID or not string
		if !IsEqual(p.Key, p.Value) {
			return nil, fmt.Errorf("Key an value must be equal for Cid type")
		}
		switch p.Key.(type) {
		case String:
		default:
			return nil, fmt.Errorf("String type expected in key")
		}
		// Decode CID
		_, err := cid.Decode(p.Key.(String).Value)
		if err != nil {
			return nil, err
		}
		outP = append(outP, p)
	}
	return outP, nil
}
