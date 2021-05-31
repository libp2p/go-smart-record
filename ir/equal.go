package ir

import xr "github.com/libp2p/go-routing-language/syntax"

func IsEqual(x, y Node) bool {
	return xr.IsEqual(x.Disassemble(), y.Disassemble())
}
