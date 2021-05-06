package ir

import "github.com/libp2p/go-smart-record/xr"

type Node interface {
	Disassemble() xr.Node // returns only syntactic nodes
	UpdateWith(ctx UpdateContext, with Node) (Node, error)
	Metadata() MetadataInfo
}

type Nodes []Node

func (ns Nodes) IndexOf(element Node) int {
	for i, p := range ns {
		if IsEqual(p, element) {
			return i
		}
	}
	return -1
}

// AreSameNodes compairs to lists of key/values for set-wise equality (order independent).
func AreSameNodes(x, y Nodes) bool {
	if len(x) != len(y) {
		return false
	}
	for _, x := range x {
		if i := y.IndexOf(x); i < 0 {
			return false
		}
	}
	return true
}
