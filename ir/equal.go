package ir

import "github.com/libp2p/go-smart-record/xr"

func IsEqual(x, y Node) bool {
	return xr.IsEqual(x.Disassemble(), y.Disassemble())
}
