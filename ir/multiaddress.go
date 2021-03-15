package ir

import "io"

type Multiaddress struct {
	Multiaddress string // TODO: This should be of type multiaddr.Multiaddr
	// User holds user fields.
	User Dict
}

func (m Multiaddress) Disassemble() Dict {
	return m.User.CopySetTag("multiaddress", String{m.Multiaddress}, String{m.Multiaddress})
}

func (m Multiaddress) WritePretty(w io.Writer) error {
	return m.Disassemble().WritePretty(w)
}

func (m Multiaddress) MergeWith(ctx MergeContext, x Node) Node {
	panic("XXX")
}
