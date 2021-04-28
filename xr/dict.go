package xr

import (
	"fmt"
	"io"
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

// jsonPair is used to encode Pairs with JSON
type jsonPair struct {
	Key   interface{}
	Value interface{}
}

func (d Dict) EncodeJSON() (interface{}, error) {
	r := struct {
		Type  marshalType   `json:"type"`
		Tag   string        `json:"tag"`
		Pairs []interface{} `json:"pairs"`
	}{Type: DictType, Tag: d.Tag, Pairs: []interface{}{}}

	for _, p := range d.Pairs {
		k, err := p.Key.EncodeJSON()
		if err != nil {
			return nil, err
		}
		v, err := p.Value.EncodeJSON()
		if err != nil {
			return nil, err
		}
		r.Pairs = append(r.Pairs, jsonPair{
			Key:   k,
			Value: v,
		})
	}
	return r, nil

}

func decodeDict(s map[string]interface{}) (Node, error) {
	r := Dict{
		Tag:   s["tag"].(string),
		Pairs: []Pair{},
	}

	pairs := s["pairs"].([]interface{})
	for _, pi := range pairs {
		p, ok := pi.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("pair is wrong type")
		}
		// Get pair values
		pk, ok := p["Key"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("key in pair is wrong type")
		}
		pv, ok := p["Value"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("value in pair is wrong type")
		}
		// Decode them
		sk, err := decodeNode(pk)
		if err != nil {
			return nil, err
		}
		sv, err := decodeNode(pv)
		if err != nil {
			return nil, err
		}
		r.Pairs = append(r.Pairs,
			Pair{
				Key:   sk,
				Value: sv,
			})
	}
	return r, nil
}

func IsEqualDict(x, y Dict) bool {
	if x.Tag != y.Tag {
		return false
	}
	return AreSamePairs(x.Pairs, y.Pairs)
}
