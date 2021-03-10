package sr

import (
	"io"
)

// Dict is a set of uniquely-named child nodes.
type Dict struct {
	Tag   string
	Pairs []Pair // maintain: keys are unique
}

type Pair struct {
	Key   Node
	Value Node
}

func (p Pair) WritePretty(w io.Writer) error {
	if !IsEqual(p.Key, p.Value) {
		if err := p.Key.WritePretty(w); err != nil {
			return err
		}
		if _, err := w.Write([]byte(" : ")); err != nil {
			return err
		}
	}
	if err := p.Value.WritePretty(IndentWriter(w)); err != nil {
		return err
	}
	return nil
}

func (d Dict) WritePretty(w io.Writer) error {
	if _, err := w.Write([]byte(d.Tag)); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'{'}); err != nil {
		return err
	}
	for _, p := range d.Pairs {
		if err := p.WritePretty(w); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	if _, err := w.Write([]byte{'}'}); err != nil {
		return err
	}
	return nil
}

func (d Dict) CopySet(key Node, value Node) Dict {
	c := Dict{Tag: d.Tag, Pairs: make([]Pair, 0, len(d.Pairs)+1)}
	found := false
	for _, p := range d.Pairs {
		if IsEqual(key, p.Key) {
			c.Pairs = append(c.Pairs, Pair{key, value})
			found = true
		} else {
			c.Pairs = append(c.Pairs, p)
		}
	}
	if !found {
		c.Pairs = append(c.Pairs, Pair{key, value})
	}
	return c
}

func (d Dict) Get(key Node) Node {
	for _, p := range d.Pairs {
		if IsEqual(p.Key, key) {
			return p.Value
		}
	}
	return nil
}

func IsEqualDict(x, y Dict) bool {
	XXX
}

func MergeDicts(x, y *Dict) Node {
	panic("XXX")
}
