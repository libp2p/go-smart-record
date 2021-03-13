package ir

import (
	"fmt"
	"io"

	peer "github.com/libp2p/go-libp2p-peer"
)

type Peer struct {
	ID string
	// User holds user fields.
	User Dict
}

func (p Peer) Dict() Dict {
	return p.User.CopySetTag("peer", String{"id"}, String{p.ID})
}

func (p Peer) WritePretty(w io.Writer) error {
	return p.Dict().WritePretty(w)
}

func ParsePeerFromDict(d Dict) (Pairs, error) {
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
		_, err := peer.IDB58Decode(p.Key.(String).Value)
		if err != nil {
			return nil, err
		}
		outP = append(outP, p)
	}
	return outP, nil
}
