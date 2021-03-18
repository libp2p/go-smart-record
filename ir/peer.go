package ir

import (
	"fmt"
	"io"
)

type Peer struct {
	ID string
	// User holds user fields.
	User Dict
}

func (p Peer) Disassemble() Dict {
	return p.User.CopySetTag("peer", String{"id"}, String{p.ID})
}

func (p Peer) WritePretty(w io.Writer) error {
	return p.Disassemble().WritePretty(w)
}

func (p Peer) MergeWith(ctx MergeContext, x Node) Node {
	panic("XXX")
}

func (p Peer) MarshalJSON() (b []byte, e error) {
	// Temporal type to avoid recursion
	/*type tmp Blob
	ts := tmp(s)

	c := struct {
		Type  MarshalType `json:"type"`
		Value tmp         `json:"value"`
	}{Type: BlobType, Value: ts}
	return json.Marshal(&c)
	*/
	return nil, fmt.Errorf("Marshal for Peer not implemented")
}
