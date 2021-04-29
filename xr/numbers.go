package xr

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/big"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

type Number interface {
	TypeIsNumber()
}

type Int struct {
	*big.Int
}

func NewInt64(v int64) Int {
	return Int{big.NewInt(v)}
}

func (n Int) TypeIsNumber() {}

func (n Int) WritePretty(w io.Writer) (err error) {
	_, err = w.Write([]byte(n.Int.String()))
	return err
}

type Float struct {
	*big.Float
}

func (n Float) TypeIsNumber() {}

func (n Float) WritePretty(w io.Writer) (err error) {
	_, err = w.Write([]byte(n.Float.String()))
	return err
}

func IsEqualNumber(x, y Number) bool {
	switch x1 := x.(type) {
	case Int:
		switch y1 := y.(type) {
		case Int:
			return x1.Int.Cmp(y1.Int) == 0
		case Float:
			return false
		}
	case Float:
		switch y1 := y.(type) {
		case Int:
			return false
		case Float:
			return x1.Float.Cmp(y1.Float) == 0
		}
	}
	panic("bug: unknown number type")
}

func (n Int) EncodeJSON() (interface{}, error) {
	bn, err := n.MarshalText()
	if err != nil {
		return nil, err
	}
	return struct {
		Type  marshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: IntType, Value: bn}, nil
}

func (n Float) EncodeJSON() (interface{}, error) {
	bn, err := n.MarshalText()
	if err != nil {
		return nil, err
	}
	return struct {
		Type  marshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: FloatType, Value: bn}, nil
}

func decodeInt(s map[string]interface{}) (Node, error) {
	z := new(big.Int)
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("wrong int decoding type")
	}
	// Unmarshaller inteprets []byte as string, we need to decode base64
	sDec, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, err
	}
	err = z.UnmarshalText(sDec)
	if err != nil {
		return nil, err
	}
	return Int{z}, nil
}

func decodeFloat(s map[string]interface{}) (Node, error) {
	z := new(big.Float)
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("wrong float decoding type")
	}
	// Unmarshaller inteprets []byte as string, we need to decode base64
	sDec, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, err
	}
	err = z.UnmarshalText(sDec)
	if err != nil {
		return nil, err
	}
	return Float{z}, nil
}

// Int IPLD node interface implementation
var (
	_ ipld.Node          = Int{}
	_ ipld.NodePrototype = Prototype__Int{}
	_ ipld.NodeBuilder   = &int__Builder{}
	_ ipld.NodeAssembler = &int__Assembler{}
)

// -- Node interface methods -->

func (Int) Kind() ipld.Kind {
	return ipld.Kind_Int
}
func (Int) LookupByString(string) (ipld.Node, error) {
	return mixins.Int{TypeName: "int"}.LookupByString("")
}
func (Int) LookupByNode(key ipld.Node) (ipld.Node, error) {
	return mixins.Int{TypeName: "int"}.LookupByNode(nil)
}
func (Int) LookupByIndex(idx int64) (ipld.Node, error) {
	return mixins.Int{TypeName: "int"}.LookupByIndex(0)
}
func (Int) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.Int{TypeName: "int"}.LookupBySegment(seg)
}
func (Int) MapIterator() ipld.MapIterator {
	return nil
}
func (Int) ListIterator() ipld.ListIterator {
	return nil
}
func (Int) Length() int64 {
	return -1
}
func (Int) IsAbsent() bool {
	return false
}
func (Int) IsNull() bool {
	return false
}
func (Int) AsBool() (bool, error) {
	return mixins.Int{TypeName: "int"}.AsBool()
}
func (n Int) AsInt() (int64, error) {
	return int64(n.Int.Int64()), nil
}
func (Int) AsFloat() (float64, error) {
	return mixins.Int{TypeName: "int"}.AsFloat()
}
func (Int) AsString() (string, error) {
	return mixins.Int{TypeName: "int"}.AsString()
}
func (Int) AsBytes() ([]byte, error) {
	return mixins.Int{TypeName: "int"}.AsBytes()
}
func (Int) AsLink() (ipld.Link, error) {
	return mixins.Int{TypeName: "int"}.AsLink()
}
func (Int) Prototype() ipld.NodePrototype {
	return Prototype__Int{}
}

// -- NodePrototype -->

type Prototype__Int struct{}

func (Prototype__Int) NewBuilder() ipld.NodeBuilder {
	var w Int
	return &int__Builder{int__Assembler{w: w}}
}

// -- NodeBuilder -->

type int__Builder struct {
	int__Assembler
}

func (nb *int__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *int__Builder) Reset() {
	var w Int
	*nb = int__Builder{int__Assembler{w: w}}
}

// -- NodeAssembler -->

type int__Assembler struct {
	w Int
}

func (int__Assembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	return mixins.IntAssembler{TypeName: "int"}.BeginMap(0)
}
func (int__Assembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	return mixins.IntAssembler{TypeName: "int"}.BeginList(0)
}
func (int__Assembler) AssignNull() error {
	return mixins.IntAssembler{TypeName: "int"}.AssignNull()
}
func (int__Assembler) AssignBool(bool) error {
	return mixins.IntAssembler{TypeName: "int"}.AssignBool(false)
}
func (na *int__Assembler) AssignInt(v int64) error {
	na.w = NewInt64(v)
	return nil
}
func (int__Assembler) AssignFloat(float64) error {
	return mixins.IntAssembler{TypeName: "int"}.AssignFloat(0)
}
func (int__Assembler) AssignString(string) error {
	return mixins.IntAssembler{TypeName: "int"}.AssignString("")
}
func (int__Assembler) AssignBytes([]byte) error {
	return mixins.IntAssembler{TypeName: "int"}.AssignBytes(nil)
}
func (int__Assembler) AssignLink(ipld.Link) error {
	return mixins.IntAssembler{TypeName: "int"}.AssignLink(nil)
}
func (na *int__Assembler) AssignNode(v ipld.Node) error {
	if v2, err := v.AsInt(); err != nil {
		return err
	} else {
		na.w = NewInt64(v2)
		return nil
	}
}
func (int__Assembler) Prototype() ipld.NodePrototype {
	return Prototype__Int{}
}

// Float IPLD node interface implementation
var (
	_ ipld.Node          = Float{}
	_ ipld.NodePrototype = Prototype__Float{}
	_ ipld.NodeBuilder   = &float__Builder{}
	_ ipld.NodeAssembler = &float__Assembler{}
)

func (Float) Kind() ipld.Kind {
	return ipld.Kind_Float
}
func (Float) LookupByString(string) (ipld.Node, error) {
	return mixins.Float{TypeName: "float"}.LookupByString("")
}
func (Float) LookupByNode(key ipld.Node) (ipld.Node, error) {
	return mixins.Float{TypeName: "float"}.LookupByNode(nil)
}
func (Float) LookupByIndex(idx int64) (ipld.Node, error) {
	return mixins.Float{TypeName: "float"}.LookupByIndex(0)
}
func (Float) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.Float{TypeName: "float"}.LookupBySegment(seg)
}
func (Float) MapIterator() ipld.MapIterator {
	return nil
}
func (Float) ListIterator() ipld.ListIterator {
	return nil
}
func (Float) Length() int64 {
	return -1
}
func (Float) IsAbsent() bool {
	return false
}
func (Float) IsNull() bool {
	return false
}
func (Float) AsBool() (bool, error) {
	return mixins.Float{TypeName: "float"}.AsBool()
}
func (Float) AsInt() (int64, error) {
	return mixins.Float{TypeName: "float"}.AsInt()
}
func (n Float) AsFloat() (float64, error) {
	// NOTE: We are disregarding the accuracy for now. This may lead to issues.
	f, _ := n.Float64()
	return float64(f), nil
}
func (Float) AsString() (string, error) {
	return mixins.Float{TypeName: "float"}.AsString()
}
func (Float) AsBytes() ([]byte, error) {
	return mixins.Float{TypeName: "float"}.AsBytes()
}
func (Float) AsLink() (ipld.Link, error) {
	return mixins.Float{TypeName: "float"}.AsLink()
}
func (Float) Prototype() ipld.NodePrototype {
	return Prototype__Float{}
}

// -- NodePrototype -->

type Prototype__Float struct{}

func (Prototype__Float) NewBuilder() ipld.NodeBuilder {
	var w Float
	return &float__Builder{float__Assembler{w: w}}
}

// -- NodeBuilder -->

type float__Builder struct {
	float__Assembler
}

func (nb *float__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *float__Builder) Reset() {
	var w Float
	*nb = float__Builder{float__Assembler{w: w}}
}

// -- NodeAssembler -->

type float__Assembler struct {
	w Float
}

func (float__Assembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	return mixins.FloatAssembler{TypeName: "float"}.BeginMap(0)
}
func (float__Assembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	return mixins.FloatAssembler{TypeName: "float"}.BeginList(0)
}
func (float__Assembler) AssignNull() error {
	return mixins.FloatAssembler{TypeName: "float"}.AssignNull()
}
func (float__Assembler) AssignBool(bool) error {
	return mixins.FloatAssembler{TypeName: "float"}.AssignBool(false)
}
func (float__Assembler) AssignInt(int64) error {
	return mixins.FloatAssembler{TypeName: "float"}.AssignInt(0)
}
func (na *float__Assembler) AssignFloat(v float64) error {
	na.w = Float{big.NewFloat(v).SetPrec(64)}
	return nil
}
func (float__Assembler) AssignString(string) error {
	return mixins.FloatAssembler{TypeName: "float"}.AssignString("")
}
func (float__Assembler) AssignBytes([]byte) error {
	return mixins.FloatAssembler{TypeName: "float"}.AssignBytes(nil)
}
func (float__Assembler) AssignLink(ipld.Link) error {
	return mixins.FloatAssembler{TypeName: "float"}.AssignLink(nil)
}
func (na *float__Assembler) AssignNode(v ipld.Node) error {
	if v2, err := v.AsFloat(); err != nil {
		return err
	} else {
		na.w = Float{big.NewFloat(v2).SetPrec(64)}
		return nil
	}
}
func (float__Assembler) Prototype() ipld.NodePrototype {
	return Prototype__Float{}
}
