package ir

import (
	"fmt"
	"io"
)

// Set is a set of (uniquely) elements.
type Set struct {
	Tag      string
	Elements Nodes
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

func MergeSet(ctx MergeContext, x, y Set) (Node, error) {
	if x.Tag != y.Tag {
		return ctx.MergeConflict(x, y)
	}
	x, y = orderSetByLength(x, y) // x is smaller, y is larger
	m := Set{
		Tag:      x.Tag,
		Elements: make(Nodes, len(y.Elements), len(x.Elements)+len(y.Elements)),
	}
	copy(m.Elements, y.Elements)
	for _, p := range x.Elements {
		if i := m.Elements.IndexOf(p); i < 0 {
			m.Elements = append(m.Elements, p)
		}
	}
	return m, nil
}

func orderSetByLength(x, y Set) (shorter, longer Set) {
	if x.Len() <= y.Len() {
		return x, y
	}
	return y, x
}
