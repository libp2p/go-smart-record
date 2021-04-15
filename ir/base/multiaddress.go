package base

import (
	"io"

	"github.com/libp2p/go-smart-record/ir"
)

type Multiaddress struct {
	Multiaddress string // TODO: This should be of type multiaddr.Multiaddr
	// User holds user fields.
	User ir.Dict
}

func (m Multiaddress) EncodeJSON() (interface{}, error) {
	return m.Disassemble().EncodeJSON()
}

func (m Multiaddress) Disassemble() ir.Dict {
	return m.User.CopySetTag("multiaddress", ir.String{m.Multiaddress}, ir.String{m.Multiaddress})
}

func (m Multiaddress) WritePretty(w io.Writer) error {
	return m.Disassemble().WritePretty(w)
}
