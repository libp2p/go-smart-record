package ir

import (
	"fmt"

	xr "github.com/libp2p/go-routing-language/syntax"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

// Pair holds a key/value pair.
type Pair struct {
	Key   Node `json:"key"`
	Value Node `json:"value"`
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

// MergePairs returns the union (wrt keys) of the two lists of pairs.
// Ties are broken in favor of y, the right argument.
func MergePairs(x, y Pairs) Pairs {
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
	Pairs       Pairs // keys must be unique wrt IsEqual
	metadataCtx *meta.Meta
}

func (d *Dict) Disassemble() xr.Node {
	x := xr.Dict{Pairs: make(xr.Pairs, len(d.Pairs))}
	for i, p := range d.Pairs {
		x.Pairs[i] = xr.Pair{Key: p.Key.Disassemble(), Value: p.Value.Disassemble()}
	}
	return x
}

func (d *Dict) Metadata() meta.MetadataInfo {
	return d.metadataCtx.GetMeta()
}

func (d Dict) Len() int {
	return len(d.Pairs)
}

func (d Dict) Copy() Dict {
	c := d
	p := make(Pairs, len(c.Pairs))
	copy(p, c.Pairs)
	c.Pairs = p
	// Also copy metadata if it exists
	m := d.metadataCtx.Copy()
	c.metadataCtx = &m
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
	return old.Value
}

func (d Dict) Get(key Node) Node {
	for _, p := range d.Pairs {
		if IsEqual(p.Key, key) {
			return p.Value
		}
	}
	return nil
}

func (d *Dict) UpdateWith(ctx UpdateContext, with Node) error {
	wd, ok := with.(*Dict)
	if !ok {
		return fmt.Errorf("cannot update with a non-dict")
	}
	for _, p := range wd.Pairs {
		if i := d.Pairs.IndexOf(p.Key); i < 0 {
			d.Pairs = append(d.Pairs, p)
		} else {
			if err := d.Pairs[i].Value.UpdateWith(ctx, p.Value); err != nil {
				return fmt.Errorf("cannout update value (%v)", err)
			}
		}
	}
	// Update metadata
	d.metadataCtx.Update(wd.metadataCtx)
	return nil
}
