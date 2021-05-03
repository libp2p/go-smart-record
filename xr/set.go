package xr

import (
	"fmt"
	"io"

	"github.com/ipld/go-ipld-prime"
	xrIpld "github.com/libp2p/go-smart-record/xr/ipld"
)

// Set is a set of (uniquely) elements.
type Set struct {
	Tag      string
	Elements Nodes
}

func (s Set) Copy() Set {
	e := make(Nodes, len(s.Elements))
	copy(e, s.Elements)
	return Set{
		Tag:      s.Tag,
		Elements: e,
	}
}

func (s Set) Len() int {
	return len(s.Elements)
}

func (s Set) WritePretty(w io.Writer) error {
	if _, err := w.Write([]byte(s.Tag)); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'{'}); err != nil {
		return err
	}
	u := IndentWriter(w)
	if _, err := u.Write([]byte{'\n'}); err != nil {
		return err
	}
	for i, p := range s.Elements {
		if err := p.WritePretty(u); err != nil {
			return err
		}
		if i+1 == len(s.Elements) {
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		} else {
			if _, err := u.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	if _, err := w.Write([]byte{'}'}); err != nil {
		return err
	}
	return nil
}

func (s Set) EncodeJSON() (interface{}, error) {
	r := struct {
		Type     marshalType   `json:"type"`
		Tag      string        `json:"tag"`
		Elements []interface{} `json:"elements"`
	}{Type: SetType, Tag: s.Tag, Elements: []interface{}{}}

	for _, n := range s.Elements {
		no, err := n.EncodeJSON()
		if err != nil {
			return nil, err
		}
		r.Elements = append(r.Elements, no)
	}
	return r, nil

}

func decodeSet(s map[string]interface{}) (Node, error) {
	r := Set{
		Tag:      s["tag"].(string),
		Elements: []Node{},
	}
	nodes, ok := s["elements"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("bad Nodes decoding format")
	}
	for _, n := range nodes {
		pv, ok := n.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("node in set element is wrong type")
		}
		nv, err := decodeNode(pv)
		if err != nil {
			return nil, err
		}
		r.Elements = append(r.Elements, nv)
	}
	return r, nil
}

func IsEqualSet(x, y Set) bool {
	if x.Tag != y.Tag {
		return false
	}
	return AreSameNodes(x.Elements, y.Elements)
}

// ToIPLD converts xr.Node into its corresponding IPLD Node type
func (s Set) ToIPLD() (ipld.Node, error) {
	// NOTE: Consider adding multierr throughout this whole function
	// Initialize Dict
	sbuild := xrIpld.Type.Set_IPLD.NewBuilder()
	ma, err := sbuild.BeginMap(-1)
	if err != nil {
		return nil, err
	}
	// Assign tag
	tasm, err := ma.AssembleEntry("Tag")
	if err != nil {
		return nil, err
	}
	err = tasm.AssignString(s.Tag)
	if err != nil {
		return nil, err
	}

	// Build elements
	lbuild := xrIpld.Type.Nodes_IPLD.NewBuilder()
	// NOTE: We can assign here directly the size of Pairs instead of -1
	la, err := lbuild.BeginList(-1)
	if err != nil {
		return nil, err
	}
	// For each pair
	for _, e := range s.Elements {

		// Add element to the list of nodes
		n, err := e.toNode_IPLD()
		if err != nil {
			return nil, err
		}
		// la.AssembleValue is Node_IPLD Assembler. Need to assemble a node
		if err := la.AssembleValue().AssignNode(n); err != nil {
			return nil, fmt.Errorf("Error assembling value: %s", err)
		}
	}
	// Finish list building
	if err := la.Finish(); err != nil {
		return nil, err
	}
	// Assign elements to set
	psasm, err := ma.AssembleEntry("Elements")
	if err != nil {
		return nil, err
	}
	err = psasm.AssignNode(lbuild.Build())
	if err != nil {
		return nil, err
	}
	// Finish elements building
	if err := ma.Finish(); err != nil {
		return nil, err
	}
	return sbuild.Build(), nil
}

// toNode_IPLD convert into IPLD Node of dynamic type NODE_IPLD
func (s Set) toNode_IPLD() (ipld.Node, error) {
	t := xrIpld.Type.Node_IPLD.NewBuilder()
	ma, err := t.BeginMap(-1)
	asm, err := ma.AssembleEntry("Set_IPLD")
	if err != nil {
		return nil, err
	}
	nd, err := s.ToIPLD()
	if err != nil {
		return nil, err
	}
	err = asm.AssignNode(nd)
	if err != nil {
		return nil, err
	}
	if err := ma.Finish(); err != nil {
		return nil, err
	}
	return t.Build(), nil
}
