package xr

import (
	"io"

	"github.com/ipld/go-ipld-prime"
)

type Node interface {
	WritePretty(w io.Writer) error    // Pretty writes the node
	EncodeJSON() (interface{}, error) // Custom JSON encoder. Not IPLD compatible
	ToIPLD() (ipld.Node, error)       // Converts xr.Node into its corresponding IPLD Node type
	toNode_IPLD() (ipld.Node, error)  // Convert into IPLD Node of dynamic type NODE_IPLD

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
