package xr

import (
	"fmt"
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
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

// IPLD node interface implementation
var (
	_ ipld.Node          = Dict{}
	_ ipld.NodePrototype = Prototype__Dict{}
	_ ipld.NodeBuilder   = &dict__Builder{}
	_ ipld.NodeAssembler = &dict__Assembler{}
)

// -- Node interface methods -->

func (Dict) Kind() ipld.Kind {
	return ipld.Kind_Map
}
func (n Dict) LookupByString(key string) (ipld.Node, error) {
	v := n.Get(String{key})
	if v == nil {
		return nil, ipld.ErrNotExists{Segment: ipld.PathSegmentOfString(key)}
	}
	return v, nil
}
func (n Dict) LookupByNode(key ipld.Node) (ipld.Node, error) {
	k, ok := key.(Node)
	if !ok {
		fmt.Errorf("Lookup key is not a xr.Node")
	}
	v := n.Get(k)
	if v == nil {
		return nil, fmt.Errorf("No Node with that key")
	}
	return v, nil
}
func (n Dict) LookupByIndex(idx int64) (ipld.Node, error) {
	// NOTE: Consider not supporting orders.
	if idx >= int64(len(n.Pairs)) {
		return nil, fmt.Errorf("index out of bounds")
	}

	return Dict{Pairs: []Pair{n.Pairs[idx]}}, nil
}
func (n Dict) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return n.LookupByString(seg.String())
}
func (n Dict) MapIterator() ipld.MapIterator {
	return &dict_MapIterator{&n, 0}
}
func (Dict) ListIterator() ipld.ListIterator {
	return nil
}
func (n Dict) Length() int64 {
	return int64(len(n.Pairs))
}
func (Dict) IsAbsent() bool {
	return false
}
func (Dict) IsNull() bool {
	return false
}
func (Dict) AsBool() (bool, error) {
	return mixins.Map{TypeName: "map"}.AsBool()
}
func (Dict) AsInt() (int64, error) {
	return mixins.Map{TypeName: "map"}.AsInt()
}
func (Dict) AsFloat() (float64, error) {
	return mixins.Map{TypeName: "map"}.AsFloat()
}
func (Dict) AsString() (string, error) {
	return mixins.Map{TypeName: "map"}.AsString()
}
func (Dict) AsBytes() ([]byte, error) {
	return mixins.Map{TypeName: "map"}.AsBytes()
}
func (Dict) AsLink() (ipld.Link, error) {
	return mixins.Map{TypeName: "map"}.AsLink()
}
func (Dict) Prototype() ipld.NodePrototype {
	return Prototype__Dict{}
}

// -- DictIterator
type dict_MapIterator struct {
	n   *Dict
	idx int
}

func (itr *dict_MapIterator) Next() (k ipld.Node, v ipld.Node, _ error) {
	if itr.Done() {
		return nil, nil, ipld.ErrIteratorOverread{}
	}
	p := itr.n.Pairs[itr.idx]

	k = p.Key
	v = p.Value
	itr.idx++
	return
}
func (itr *dict_MapIterator) Done() bool {
	return itr.idx >= len(itr.n.t)
}

// -- NodePrototype -->

type Prototype__Dict struct{}

func (Prototype__Dict) NewBuilder() ipld.NodeBuilder {
	var w Dict
	return &dict__Builder{dict__Assembler{w: w}}
}

// -- NodeBuilder -->

type dict__Builder struct {
	dict__Assembler
}

func (nb *dict__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *dict__Builder) Reset() {
	var w Dict
	*nb = dict__Builder{string__Assembler{w: Dict{}}}
}

// -- NodeAssembler -->

type dict__Assembler struct {
	w  Dict
	ka dict__KeyAssembler
	va dict__ValueAssembler

	state maState
}
type dict__KeyAssembler struct {
	ma *dict__Assembler
}
type dict__ValueAssembler struct {
	ma *dict__Assembler
}

// maState is an enum of the state machine for a map assembler.
// (this might be something to export reusably, but it's also very much an impl detail that need not be seen, so, dubious.)
type maState uint8

const (
	maState_initial     maState = iota // also the 'expect key or finish' state
	maState_midKey                     // waiting for a 'finished' state in the KeyAssembler.
	maState_expectValue                // 'AssembleValue' is the only valid next step
	maState_midValue                   // waiting for a 'finished' state in the ValueAssembler.
	maState_finished                   // 'w' will also be nil, but this is a politer statement
)

func (na *dict__Assembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	if sizeHint < 0 {
		sizeHint = 0
	}
	// Allocate storage space.
	na.w.Pairs = make([]Pair, 0, sizeHint)
	// That's it; return self as the MapAssembler.  We already have all the right methods on this structure.
	return na, nil
}
func (dict__Assembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	return mixins.MapAssembler{TypeName: "map"}.BeginList(0)
}
func (dict__Assembler) AssignNull() error {
	return mixins.MapAssembler{TypeName: "map"}.AssignNull()
}
func (dict__Assembler) AssignBool(bool) error {
	return mixins.MapAssembler{TypeName: "map"}.AssignBool(false)
}
func (dict__Assembler) AssignInt(int64) error {
	return mixins.MapAssembler{TypeName: "map"}.AssignInt(0)
}
func (dict__Assembler) AssignFloat(float64) error {
	return mixins.MapAssembler{TypeName: "map"}.AssignFloat(0)
}
func (dict__Assembler) AssignString(string) error {
	return mixins.MapAssembler{TypeName: "map"}.AssignString("")
}
func (dict__Assembler) AssignBytes([]byte) error {
	return mixins.MapAssembler{TypeName: "map"}.AssignBytes(nil)
}
func (dict__Assembler) AssignLink(ipld.Link) error {
	return mixins.MapAssembler{TypeName: "map"}.AssignLink(nil)
}
func (na *dict__Assembler) AssignNode(v ipld.Node) error {
	// Sanity check assembler state.
	//  Update of state to 'finished' comes later; where exactly depends on if shortcuts apply.
	if na.state != maState_initial {
		panic("misuse")
	}
	// Copy the content.
	if v2, ok := v.(Dict); ok { // if our own type: shortcut.
		// Copy the structure by value.
		//  This means we'll have pointers into the same internal maps and slices;
		//   this is okay, because the Node type promises it's immutable, and we are going to instantly finish ourselves to also maintain that.
		// FIXME: the shortcut behaves differently than the long way: it discards any existing progress.  Doesn't violate immut, but is odd.
		na.w = v2
		na.state = maState_finished
		return nil
	}
	// If the above shortcut didn't work, resort to a generic copy.
	//  We call AssignNode for all the child values, giving them a chance to hit shortcuts even if we didn't.
	if v.Kind() != ipld.Kind_Map {
		return ipld.ErrWrongKind{TypeName: "map", MethodName: "AssignNode", AppropriateKind: ipld.KindSet_JustMap, ActualKind: v.Kind()}
	}
	itr := v.MapIterator()
	for !itr.Done() {
		k, v, err := itr.Next()
		if err != nil {
			return err
		}
		if err := na.AssembleKey().AssignNode(k); err != nil {
			return err
		}
		if err := na.AssembleValue().AssignNode(v); err != nil {
			return err
		}
	}
	return na.Finish()
}
func (dict__Assembler) Prototype() ipld.NodePrototype {
	return Prototype__Dict{}
}

func (ma *dict__Assembler) AssembleKey() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_midKey
	// Make key assembler valid by giving it pointer back to whole 'ma'; yield it.
	ma.ka.ma = ma
	return &ma.ka
}

func (ma *dict__Assembler) AssembleValue() ipld.NodeAssembler {
	// Sanity check, then update, assembler state.
	if ma.state != maState_expectValue {
		panic("misuse")
	}
	ma.state = maState_midValue
	// Make value assembler valid by giving it pointer back to whole 'ma'; yield it.
	ma.va.ma = ma
	return &ma.va
}

func (ma *dict__Assembler) Finish() error {
	// Sanity check, then update, assembler state.
	if ma.state != maState_initial {
		panic("misuse")
	}
	ma.state = maState_finished
	// validators could run and report errors promptly, if this type had any.
	return nil
}
func (dict__Assembler) KeyPrototype() ipld.NodePrototype {
	return Prototype__Any{}
}
func (dict__Assembler) ValuePrototype(_ string) ipld.NodePrototype {
	return Prototype__Any{}
}

// -- MapAssembler.KeyAssembler -->

func (plainMap__KeyAssembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	return mixins.StringAssembler{TypeName: "string"}.BeginMap(0)
}
func (plainMap__KeyAssembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	return mixins.StringAssembler{TypeName: "string"}.BeginList(0)
}
func (plainMap__KeyAssembler) AssignNull() error {
	return mixins.StringAssembler{TypeName: "string"}.AssignNull()
}
func (plainMap__KeyAssembler) AssignBool(bool) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignBool(false)
}
func (plainMap__KeyAssembler) AssignInt(int64) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignInt(0)
}
func (plainMap__KeyAssembler) AssignFloat(float64) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignFloat(0)
}
func (mka *plainMap__KeyAssembler) AssignString(v string) error {
	// Check for dup keys; error if so.
	//  (And, backtrack state to accepting keys again so we don't get eternally wedged here.)
	_, exists := mka.ma.w.m[v]
	if exists {
		mka.ma.state = maState_initial
		mka.ma = nil // invalidate self to prevent further incorrect use.
		return ipld.ErrRepeatedMapKey{Key: plainString(v)}
	}
	// Assign the key into the end of the entry table;
	//  we'll be doing map insertions after we get the value in hand.
	//  (There's no need to delegate to another assembler for the key type,
	//   because we're just at Data Model level here, which only regards plain strings.)
	mka.ma.w.t = append(mka.ma.w.t, plainMap__Entry{})
	mka.ma.w.t[len(mka.ma.w.t)-1].k = plainString(v)
	// Update parent assembler state: clear to proceed.
	mka.ma.state = maState_expectValue
	mka.ma = nil // invalidate self to prevent further incorrect use.
	return nil
}
func (plainMap__KeyAssembler) AssignBytes([]byte) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignBytes(nil)
}
func (plainMap__KeyAssembler) AssignLink(ipld.Link) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignLink(nil)
}
func (mka *plainMap__KeyAssembler) AssignNode(v ipld.Node) error {
	vs, err := v.AsString()
	if err != nil {
		return fmt.Errorf("cannot assign non-string node into map key assembler") // FIXME:errors: this doesn't quite fit in ErrWrongKind cleanly; new error type?
	}
	return mka.AssignString(vs)
}
func (plainMap__KeyAssembler) Prototype() ipld.NodePrototype {
	return Prototype__String{}
}

// -- MapAssembler.ValueAssembler -->

func (mva *plainMap__ValueAssembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	ma := plainMap__ValueAssemblerMap{}
	ma.ca.w = &plainMap{}
	ma.p = mva.ma
	_, err := ma.ca.BeginMap(sizeHint)
	return &ma, err
}
func (mva *plainMap__ValueAssembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	la := plainMap__ValueAssemblerList{}
	la.ca.w = &plainList{}
	la.p = mva.ma
	_, err := la.ca.BeginList(sizeHint)
	return &la, err
}
func (mva *plainMap__ValueAssembler) AssignNull() error {
	return mva.AssignNode(ipld.Null)
}
func (mva *plainMap__ValueAssembler) AssignBool(v bool) error {
	vb := plainBool(v)
	return mva.AssignNode(&vb)
}
func (mva *plainMap__ValueAssembler) AssignInt(v int64) error {
	vb := plainInt(v)
	return mva.AssignNode(&vb)
}
func (mva *plainMap__ValueAssembler) AssignFloat(v float64) error {
	vb := plainFloat(v)
	return mva.AssignNode(&vb)
}
func (mva *plainMap__ValueAssembler) AssignString(v string) error {
	vb := plainString(v)
	return mva.AssignNode(&vb)
}
func (mva *plainMap__ValueAssembler) AssignBytes(v []byte) error {
	vb := plainBytes(v)
	return mva.AssignNode(&vb)
}
func (mva *plainMap__ValueAssembler) AssignLink(v ipld.Link) error {
	vb := plainLink{v}
	return mva.AssignNode(&vb)
}
func (mva *plainMap__ValueAssembler) AssignNode(v ipld.Node) error {
	l := len(mva.ma.w.t) - 1
	mva.ma.w.t[l].v = v
	mva.ma.w.m[string(mva.ma.w.t[l].k)] = v
	mva.ma.state = maState_initial
	mva.ma = nil // invalidate self to prevent further incorrect use.
	return nil
}
func (plainMap__ValueAssembler) Prototype() ipld.NodePrototype {
	return Prototype__Any{}
}
