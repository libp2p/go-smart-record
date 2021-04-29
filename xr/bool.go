package xr

import (
	"fmt"
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

type Bool struct {
	Value bool
}

func (b Bool) WritePretty(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, "%v", b.Value)
	return err
}

func IsEqualBool(x, y Bool) bool {
	return x.Value == y.Value
}

func (b Bool) EncodeJSON() (interface{}, error) {
	return struct {
		Type  marshalType `json:"type"`
		Value bool        `json:"value"`
	}{Type: BoolType, Value: b.Value}, nil
}

func decodeBool(s map[string]interface{}) (Node, error) {
	r, ok := s["value"].(bool)
	if !ok {
		return nil, fmt.Errorf("decoded value not Bool")
	}
	return Bool{r}, nil
}

// IPLD node interface implementation
var (
	_ ipld.Node          = Bool{}
	_ ipld.NodePrototype = Prototype__Bool{}
	_ ipld.NodeBuilder   = &bool__Builder{}
	_ ipld.NodeAssembler = &bool__Assembler{}
)

// -- Node interface methods -->

func (Bool) Kind() ipld.Kind {
	return ipld.Kind_Bool
}
func (Bool) LookupByString(string) (ipld.Node, error) {
	return mixins.Bool{TypeName: "bool"}.LookupByString("")
}
func (Bool) LookupByNode(key ipld.Node) (ipld.Node, error) {
	return mixins.Bool{TypeName: "bool"}.LookupByNode(nil)
}
func (Bool) LookupByIndex(idx int64) (ipld.Node, error) {
	return mixins.Bool{TypeName: "bool"}.LookupByIndex(0)
}
func (Bool) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.Bool{TypeName: "bool"}.LookupBySegment(seg)
}
func (Bool) MapIterator() ipld.MapIterator {
	return nil
}
func (Bool) ListIterator() ipld.ListIterator {
	return nil
}
func (Bool) Length() int64 {
	return -1
}
func (Bool) IsAbsent() bool {
	return false
}
func (Bool) IsNull() bool {
	return false
}
func (n Bool) AsBool() (bool, error) {
	return bool(n.Value), nil
}
func (Bool) AsInt() (int64, error) {
	return mixins.Bool{TypeName: "bool"}.AsInt()
}
func (Bool) AsFloat() (float64, error) {
	return mixins.Bool{TypeName: "bool"}.AsFloat()
}
func (Bool) AsString() (string, error) {
	return mixins.Bool{TypeName: "bool"}.AsString()
}
func (Bool) AsBytes() ([]byte, error) {
	return mixins.Bool{TypeName: "bool"}.AsBytes()
}
func (Bool) AsLink() (ipld.Link, error) {
	return mixins.Bool{TypeName: "bool"}.AsLink()
}
func (Bool) Prototype() ipld.NodePrototype {
	return Prototype__Bool{}
}

// -- NodePrototype -->

type Prototype__Bool struct{}

func (Prototype__Bool) NewBuilder() ipld.NodeBuilder {
	var w Bool
	return &bool__Builder{bool__Assembler{w: w}}
}

// -- NodeBuilder -->

type bool__Builder struct {
	bool__Assembler
}

func (nb *bool__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *bool__Builder) Reset() {
	var w Bool
	*nb = bool__Builder{bool__Assembler{w: w}}
}

// -- NodeAssembler -->

type bool__Assembler struct {
	w Bool
}

func (bool__Assembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	return mixins.BoolAssembler{TypeName: "bool"}.BeginMap(0)
}
func (bool__Assembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	return mixins.BoolAssembler{TypeName: "bool"}.BeginList(0)
}
func (bool__Assembler) AssignNull() error {
	return mixins.BoolAssembler{TypeName: "bool"}.AssignNull()
}
func (na *bool__Assembler) AssignBool(v bool) error {
	na.w = Bool{Value: v}
	return nil
}
func (bool__Assembler) AssignInt(int64) error {
	return mixins.BoolAssembler{TypeName: "bool"}.AssignInt(0)
}
func (bool__Assembler) AssignFloat(float64) error {
	return mixins.BoolAssembler{TypeName: "bool"}.AssignFloat(0)
}
func (bool__Assembler) AssignString(string) error {
	return mixins.BoolAssembler{TypeName: "bool"}.AssignString("")
}
func (bool__Assembler) AssignBytes([]byte) error {
	return mixins.BoolAssembler{TypeName: "bool"}.AssignBytes(nil)
}
func (bool__Assembler) AssignLink(ipld.Link) error {
	return mixins.BoolAssembler{TypeName: "bool"}.AssignLink(nil)
}
func (na *bool__Assembler) AssignNode(v ipld.Node) error {
	if v2, err := v.AsBool(); err != nil {
		return err
	} else {
		na.w = Bool{Value: v2}
		return nil
	}
}
func (bool__Assembler) Prototype() ipld.NodePrototype {
	return Prototype__Bool{}
}
