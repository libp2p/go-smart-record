package ir

import "io"

type Multiaddress struct {
	Multiaddress string
	// User holds user fields.
	User Dict
}

func (m Multiaddress) Dict() Dict {
	return m.User.CopySetTag("multiaddress", String{m.Multiaddress}, String{m.Multiaddress})
}

func (m Multiaddress) WritePretty(w io.Writer) error {
	return m.Dict().WritePretty(w)
}

func (m Multiaddress) MergeWith(ctx MergeContext, x Node) Node {
	panic("XXX")
}
