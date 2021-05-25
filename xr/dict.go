package xr

import (
	"io"

	"github.com/ipld/go-ipld-prime"
	xrIpld "github.com/libp2p/go-smart-record/xr/ipld"
)

// Pair holds a key/value pair.
type Pair struct {
	Key   Node `json:"key"`
	Value Node `json:"value"`
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

func (d Dict) ToIPLD() (ipld.Node, error) {
	// NOTE: Consider adding multierr throughout this whole function
	// Initialize Dict
	dbuild := xrIpld.Type.Dict_IPLD.NewBuilder()
	ma, err := dbuild.BeginMap(-1)
	if err != nil {
		return nil, err
	}
	// Assign tag
	tasm, err := ma.AssembleEntry("Tag")
	if err != nil {
		return nil, err
	}
	err = tasm.AssignString(d.Tag)
	if err != nil {
		return nil, err
	}

	// Build pairs
	lbuild := xrIpld.Type.Pairs_IPLD.NewBuilder()
	// NOTE: We can assign here directly the size of Pairs instead of -1
	la, err := lbuild.BeginList(-1)
	if err != nil {
		return nil, err
	}
	// For each pair
	for _, p := range d.Pairs {
		pbuild := xrIpld.Type.Pair_IPLD.NewBuilder()
		// NOTE: We can initialize the map with 2 instead of -1
		pa, err := pbuild.BeginMap(-1)
		if err != nil {
			return nil, err
		}
		k, err := p.Key.toNode_IPLD()
		if err != nil {
			return nil, err
		}

		// Create key to IPLDNode and assign
		kasm, err := pa.AssembleEntry("Key")
		if err != nil {
			return nil, err
		}
		err = kasm.AssignNode(k)
		if err != nil {
			return nil, err
		}
		// Create value to IPLDNode and assign
		v, err := p.Value.toNode_IPLD()
		if err != nil {
			return nil, err
		}
		vasm, err := pa.AssembleEntry("Value")
		if err != nil {
			return nil, err
		}
		err = vasm.AssignNode(v)
		if err != nil {
			return nil, err
		}
		// Finish pair building
		if err := pa.Finish(); err != nil {
			return nil, err
		}
		// Add pair to the list of pairs
		if err := la.AssembleValue().AssignNode(pbuild.Build()); err != nil {
			return nil, err
		}
	}
	// Finish building pairs
	if err := la.Finish(); err != nil {
		return nil, err
	}
	// Assign list of pairs to dict
	psasm, err := ma.AssembleEntry("Pairs")
	if err != nil {
		return nil, err
	}
	err = psasm.AssignNode(lbuild.Build())
	if err != nil {
		return nil, err
	}
	// Finish building dict
	if err := ma.Finish(); err != nil {
		return nil, err
	}
	return dbuild.Build(), nil
}

func (d Dict) Copy() Dict {
	c := d
	p := make(Pairs, len(c.Pairs))
	copy(p, c.Pairs)
	c.Pairs = p
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

func IsEqualDict(x, y Dict) bool {
	if x.Tag != y.Tag {
		return false
	}
	return AreSamePairs(x.Pairs, y.Pairs)
}

// toNode_IPLD convert into IPLD Node of dynamic type NODE_IPLD
func (d Dict) toNode_IPLD() (ipld.Node, error) {
	t := xrIpld.Type.Node_IPLD.NewBuilder()
	ma, err := t.BeginMap(-1)
	asm, err := ma.AssembleEntry("Dict_IPLD")
	if err != nil {
		return nil, err
	}
	nd, err := d.ToIPLD()
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
