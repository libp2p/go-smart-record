package base

import (
	"io"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/xr"
)

type Multiaddress struct {
	Multiaddress string // TODO: This should be of type multiaddr.Multiaddr
	// User holds user fields.
	User ir.Dict
}

func (m Multiaddress) EncodeJSON() (interface{}, error) {
	return m.Disassemble().EncodeJSON()
}

func (m Multiaddress) Disassemble() xr.Node {
	return m.User.CopySetTag("multiaddress", ir.String{m.Multiaddress}, ir.String{m.Multiaddress}).Disassemble()
}

func (m Multiaddress) WritePretty(w io.Writer) error {
	return m.Disassemble().WritePretty(w)
}
