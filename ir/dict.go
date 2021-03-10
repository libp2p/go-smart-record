package ir

import (
	"io"
)

// Pair holds a key/value pair.
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

// Pairs is a list of pairs.
type Pairs []Pair

func (ps Pairs) IndexOf(key Node) int {
	for i, p := range ps {
		if IsEqual(p.Key, key) {
			return i
		}
	}
	return -1
}

// MergePairsRight returns the union (wrt keys) of the two lists of pairs.
// Ties are broken in favor of y, the right argument.
func MergePairsRight(x, y Pairs) Pairs {
	z := make(Pairs, len(x), len(x)+len(y))
	copy(z, x)
	for _, p := range y {
		if i := z.IndexOf(p.Key); i < 0 {
			z = append(z, p)
		} else {
			z[i] = p
		}
	}
	return z
}

// Dict is a set of uniquely-named child nodes.
type Dict struct {
	Tag   string
	Pairs Pairs // keys must be unique wrt IsEqual
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
	return d.CopySetTag(d.Tag, key, value)
}

func (d Dict) CopySetTag(tag string, key Node, value Node) Dict {
	c := Dict{Tag: tag, Pairs: make(Pairs, 0, len(d.Pairs)+1)}
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
	panic("XXX")
}

func MergeDicts(x, y *Dict) Node {
	panic("XXX")
}
