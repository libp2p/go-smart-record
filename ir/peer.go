package ir

import (
	"io"

	peer "github.com/libp2p/go-libp2p-peer"
)

type Peer struct {
	//ID string
	// User holds user fields.
	User Dict
}

func (px Peer) WritePretty(w io.Writer) error {
	return px.Dict().WritePretty(w)
}

func (px Peer) Dict() Dict {
	outP := Pairs{}
	for _, p := range px.User.Pairs {
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
		//_, err := peer.IDFromString(p.Key.(String).Value)
		_, err := peer.IDB58Decode(p.Key.(String).Value)
		if err != nil {
			p.Value = String{"Err"}
		}
		outP = append(outP, p)
	}
	return Dict{Tag: px.User.Tag, Pairs: outP}
}
