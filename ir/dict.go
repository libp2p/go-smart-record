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

// AreSamePairs compairs to lists of key/values for set-wise equality (order independent).
func AreSamePairs(x, y Pairs) bool {
	if len(x) != len(y) {
		return false
	}
	for _, x := range x {
		if i := y.IndexOf(x.Key); i < 0 {
			return false
		} else {
			if !IsEqual(x.Value, y[i].Value) {
				return false
			}
		}
	}
	return true
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

// Dict is a set of uniquely-keyed values.
type Dict struct {
	Tag   string
	Pairs Pairs // keys must be unique wrt IsEqual
}

func (d Dict) Len() int {
	return len(d.Pairs)
}

func (d Dict) WritePretty(w io.Writer) error {
	if _, err := w.Write([]byte(d.Tag)); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'('}); err != nil {
		return err
	}
	u := IndentWriter(w)
	if _, err := u.Write([]byte{'\n'}); err != nil {
		return err
	}
	for i, p := range d.Pairs {
		if err := p.WritePretty(u); err != nil {
			return err
		}
		if i+1 == len(d.Pairs) {
			if _, err := w.Write([]byte("\n")); err != nil {
				return err
			}
		} else {
			if _, err := u.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	if _, err := w.Write([]byte{')'}); err != nil {
		return err
	}
	return nil
}

func (d Dict) Copy() Dict {
	c := d
	p := make(Pairs, len(c.Pairs))
	copy(p, c.Pairs)
	c.Pairs = p
	return c
}

func (d Dict) CopySet(key Node, value Node) Dict {
	return d.CopySetTag(d.Tag, key, value)
}

func (d Dict) CopySetTag(tag string, key Node, value Node) Dict {
	c := d.Copy()
	c.Tag = tag
	found := false
	for i, p := range c.Pairs {
		if IsEqual(key, p.Key) {
			c.Pairs[i] = Pair{key, value}
			found = true
			break
		}
	}
	if !found {
		c.Pairs = append(c.Pairs, Pair{key, value})
	}
	return c
}

func (d *Dict) Remove(key Node) Node {
	i := d.Pairs.IndexOf(key)
	if i < 0 {
		return nil
	}
	old := d.Pairs[i]
	n := len(d.Pairs)
	d.Pairs[i], d.Pairs[n-1] = d.Pairs[n-1], d.Pairs[i]
	d.Pairs = d.Pairs[:n-1]
	return old
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
	if x.Tag != y.Tag {
		return false
	}
	return AreSamePairs(x.Pairs, y.Pairs)
}

func MergeDict(ctx MergeContext, x, y Dict) (Node, error) {
	if x.Tag != y.Tag {
		return ctx.MergeConflict(x, y)
	}
	return MergeDictIgnoreTag(ctx, x, y)
}

func MergeDictIgnoreTag(ctx MergeContext, x, y Dict) (Node, error) {
	x, y = orderDictByLength(x, y) // x is smaller, y is larger
	m := Dict{
		Tag:   x.Tag,
		Pairs: make(Pairs, len(y.Pairs), len(x.Pairs)+len(y.Pairs)),
	}
	copy(m.Pairs, y.Pairs)
	for _, p := range x.Pairs {
		if i := m.Pairs.IndexOf(p.Key); i < 0 {
			m.Pairs = append(m.Pairs, p)
		} else {
			var err error
			m.Pairs[i].Value, err = Merge(ctx, p.Value, m.Pairs[i].Value)
			if err != nil {
				return nil, err
			}
		}
	}
	return m, nil
}

func orderDictByLength(x, y Dict) (shorter, longer Dict) {
	if x.Len() <= y.Len() {
		return x, y
	}
	return y, x
}
