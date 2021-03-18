package ir

import (
	"encoding/json"
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

func (p *Pair) UnmarshalJSON(b []byte) error {

	var objMap map[string]map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}
	var kType, vType int
	json.Unmarshal(*objMap["key"]["type"], &kType)
	json.Unmarshal(*objMap["value"]["type"], &vType)
	p.Key, err = UnmarshalType(MarshalType(kType), *objMap["key"]["value"])
	if err != nil {
		return err
	}
	p.Value, err = UnmarshalType(MarshalType(vType), *objMap["value"]["value"])
	if err != nil {
		return err
	}
	return nil
}

/*
func (p *Pair) UnmarshalJSON2(b []byte) error {
	// Get types for key and value
	var objMap map[string]map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}
	var k, v int
	json.Unmarshal(*objMap["key"]["type"], &k)
	json.Unmarshal(*objMap["value"]["type"], &v)

	// TODO: Here we check the type and then we use the unmarshal directly
	// We may need to assign types.
	// Unmarshal according to type
	if k == StringType {
		// var key String
		// err := json.Unmarshal(*objMap["key"]["value"], &key)
		// if err != nil {
		//         return err
		// }
		p.Key, err = UnmarshalType(StringType, *objMap["key"]["value"])
		if err != nil {
			return err
		}
	}

	if v == StringType {
		var value String
		err := json.Unmarshal(*objMap["value"]["value"], &value)
		if err != nil {
			return err
		}
		p.Value = value
	}

	if k == BlobType {
		var key Blob
		err := json.Unmarshal(*objMap["key"]["value"], &key)
		if err != nil {
			return err
		}
		p.Key = key
	}

	if v == BlobType {
		var value Blob
		err := json.Unmarshal(*objMap["value"]["value"], &value)
		if err != nil {
			return err
		}
		p.Value = value
	}

	return nil
}*/

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

func AreEqualPairs(x, y Pairs) bool {
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

func (d Dict) WritePretty(w io.Writer) error {
	if _, err := w.Write([]byte(d.Tag)); err != nil {
		return err
	}
	if _, err := w.Write([]byte{'{'}); err != nil {
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

func (d Dict) MarshalJSON() (b []byte, e error) {
	// Temporal type to avoid recursion
	type tmp Dict
	ts := tmp(d)

	return json.Marshal(&ts)
}

func (d *Dict) UnmarshalJSON(data []byte) error {
	// Temporal type to avoid recursion
	type tmp Dict
	var ts tmp

	err := json.Unmarshal(data, &ts)
	if err != nil {
		return err
	}
	*d = Dict(ts)
	return nil
}

func IsEqualDict(x, y Dict) bool {
	if x.Tag != y.Tag {
		return false
	}
	return AreEqualPairs(x.Pairs, y.Pairs)
}

func MergeDict(ctx MergeContext, x, y Dict) Node {
	if x.Tag != y.Tag {
		return ctx.MergeConflict(x, y)
	}
	panic("XXX")
}
