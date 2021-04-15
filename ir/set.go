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

func (s Set) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	ws, ok := with.(Set)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-set")
	}
	// if ws.Tag != s.Tag {
	// 	return nil, fmt.Errorf("cannot change set tag")
	// }
	u := s.Copy()
	u.Tag = ws.Tag
	for _, e := range ws.Elements {
		if i := u.Elements.IndexOf(e); i < 0 {
			u.Elements = append(u.Elements, e)
		}
	}
	return u, nil
}
