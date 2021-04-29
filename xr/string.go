package xr

import (
	"fmt"
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

// String is a node representing a string literal.
type String struct {
	Value string
}

func (s String) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%q", s.Value)
	return err
}

func IsEqualString(x, y String) bool {
	return x.Value == y.Value
}

func (s String) EncodeJSON() (interface{}, error) {
	return struct {
		Type  marshalType `json:"type"`
		Value string      `json:"value"`
	}{Type: StringType, Value: s.Value}, nil
}

func decodeString(s map[string]interface{}) (Node, error) {
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("decoded value not String")
	}
	return String{r}, nil
}

// IPLD node interface implementation
var (
	_ ipld.Node          = String{}
	_ ipld.NodePrototype = Prototype__String{}
	_ ipld.NodeBuilder   = &string__Builder{}
	_ ipld.NodeAssembler = &string__Assembler{}
)

// -- Node interface methods -->

func (String) Kind() ipld.Kind {
	return ipld.Kind_String
}
func (String) LookupByString(string) (ipld.Node, error) {
	return mixins.String{TypeName: "string"}.LookupByString("")
}
func (String) LookupByNode(key ipld.Node) (ipld.Node, error) {
	return mixins.String{TypeName: "string"}.LookupByNode(nil)
}
func (String) LookupByIndex(idx int64) (ipld.Node, error) {
	return mixins.String{TypeName: "string"}.LookupByIndex(0)
}
func (String) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.String{TypeName: "string"}.LookupBySegment(seg)
}
func (String) MapIterator() ipld.MapIterator {
	return nil
}
func (String) ListIterator() ipld.ListIterator {
	return nil
}
func (String) Length() int64 {
	return -1
}
func (String) IsAbsent() bool {
	return false
}
func (String) IsNull() bool {
	return false
}
func (String) AsBool() (bool, error) {
	return mixins.String{TypeName: "string"}.AsBool()
}
func (String) AsInt() (int64, error) {
	return mixins.String{TypeName: "string"}.AsInt()
}
func (String) AsFloat() (float64, error) {
	return mixins.String{TypeName: "string"}.AsFloat()
}
func (x String) AsString() (string, error) {
	return string(x.Value), nil
}
func (String) AsBytes() ([]byte, error) {
	return mixins.String{TypeName: "string"}.AsBytes()
}
func (String) AsLink() (ipld.Link, error) {
	return mixins.String{TypeName: "string"}.AsLink()
}
func (String) Prototype() ipld.NodePrototype {
	return Prototype__String{}
}

// -- NodePrototype -->

type Prototype__String struct{}

func (Prototype__String) NewBuilder() ipld.NodeBuilder {
	var w String
	return &string__Builder{string__Assembler{w: w}}
}

// -- NodeBuilder -->

type string__Builder struct {
	string__Assembler
}

func (nb *string__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *string__Builder) Reset() {
	var w String
	*nb = string__Builder{string__Assembler{w: w}}
}

// -- NodeAssembler -->

type string__Assembler struct {
	w String
}

func (string__Assembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	return mixins.StringAssembler{TypeName: "string"}.BeginMap(0)
}
func (string__Assembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	return mixins.StringAssembler{TypeName: "string"}.BeginList(0)
}
func (string__Assembler) AssignNull() error {
	return mixins.StringAssembler{TypeName: "string"}.AssignNull()
}
func (string__Assembler) AssignBool(bool) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignBool(false)
}
func (string__Assembler) AssignInt(int64) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignInt(0)
}
func (string__Assembler) AssignFloat(float64) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignFloat(0)
}
func (na *string__Assembler) AssignString(v string) error {
	na.w = String{Value: v}
	return nil
}
func (string__Assembler) AssignBytes([]byte) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignBytes(nil)
}
func (string__Assembler) AssignLink(ipld.Link) error {
	return mixins.StringAssembler{TypeName: "string"}.AssignLink(nil)
}
func (na *string__Assembler) AssignNode(v ipld.Node) error {
	if v2, err := v.AsString(); err != nil {
		return err
	} else {
		na.w = String{Value: v2}
		return nil
	}
}
func (string__Assembler) Prototype() ipld.NodePrototype {
	return Prototype__String{}
}
