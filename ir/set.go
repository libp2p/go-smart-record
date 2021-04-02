package ir

import (
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
