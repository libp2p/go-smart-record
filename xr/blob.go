package xr

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"

	ipld "github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/node/mixins"
)

type Blob struct {
	Bytes []byte
}

func (b Blob) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "0x%s", hex.EncodeToString(b.Bytes)) // TODO: We can do better. E.g. wrap on 80-column boundary.
	return err
}

func (b Blob) EncodeJSON() (interface{}, error) {
	return struct {
		Type  marshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: BlobType, Value: b.Bytes}, nil
}

func IsEqualBlob(x, y Blob) bool {
	return bytes.Compare(x.Bytes, y.Bytes) == 0
}

func decodeBlob(s map[string]interface{}) (Node, error) {
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("decoding typ is not Blob")
	}
	// Unmarshaller inteprets []byte as string, we need to decode base64
	sDec, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, err
	}
	return Blob{sDec}, nil
}

// -- Node interface methods -->

func (Blob) Kind() ipld.Kind {
	return ipld.Kind_Bytes
}
func (Blob) LookupByString(string) (ipld.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupByString("")
}
func (Blob) LookupByNode(key ipld.Node) (ipld.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupByNode(nil)
}
func (Blob) LookupByIndex(idx int64) (ipld.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupByIndex(0)
}
func (Blob) LookupBySegment(seg ipld.PathSegment) (ipld.Node, error) {
	return mixins.Bytes{TypeName: "bytes"}.LookupBySegment(seg)
}
func (Blob) MapIterator() ipld.MapIterator {
	return nil
}
func (Blob) ListIterator() ipld.ListIterator {
	return nil
}
func (Blob) Length() int64 {
	return -1
}
func (Blob) IsAbsent() bool {
	return false
}
func (Blob) IsNull() bool {
	return false
}
func (Blob) AsBool() (bool, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsBool()
}
func (Blob) AsInt() (int64, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsInt()
}
func (Blob) AsFloat() (float64, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsFloat()
}
func (Blob) AsString() (string, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsString()
}
func (n Blob) AsBytes() ([]byte, error) {
	return []byte(n.Bytes), nil
}
func (Blob) AsLink() (ipld.Link, error) {
	return mixins.Bytes{TypeName: "bytes"}.AsLink()
}
func (Blob) Prototype() ipld.NodePrototype {
	return Prototype__Blob{}
}

// IPLD node interface implementation
var (
	_ ipld.Node          = Blob{}
	_ ipld.NodePrototype = Prototype__Blob{}
	_ ipld.NodeBuilder   = &blob__Builder{}
	_ ipld.NodeAssembler = &blob__Assembler{}
)

// -- NodePrototype -->

type Prototype__Blob struct{}

func (Prototype__Blob) NewBuilder() ipld.NodeBuilder {
	var w Blob
	return &blob__Builder{blob__Assembler{w: w}}
}

// -- NodeBuilder -->

type blob__Builder struct {
	blob__Assembler
}

func (nb *blob__Builder) Build() ipld.Node {
	return nb.w
}
func (nb *blob__Builder) Reset() {
	var w Blob
	*nb = blob__Builder{blob__Assembler{w: w}}
}

// -- NodeAssembler -->

type blob__Assembler struct {
	w Blob
}

func (blob__Assembler) BeginMap(sizeHint int64) (ipld.MapAssembler, error) {
	return mixins.BytesAssembler{TypeName: "bytes"}.BeginMap(0)
}
func (blob__Assembler) BeginList(sizeHint int64) (ipld.ListAssembler, error) {
	return mixins.BytesAssembler{TypeName: "bytes"}.BeginList(0)
}
func (blob__Assembler) AssignNull() error {
	return mixins.BytesAssembler{TypeName: "bytes"}.AssignNull()
}
func (blob__Assembler) AssignBool(bool) error {
	return mixins.BytesAssembler{TypeName: "bytes"}.AssignBool(false)
}
func (blob__Assembler) AssignInt(int64) error {
	return mixins.BytesAssembler{TypeName: "bytes"}.AssignInt(0)
}
func (blob__Assembler) AssignFloat(float64) error {
	return mixins.BytesAssembler{TypeName: "bytes"}.AssignFloat(0)
}
func (blob__Assembler) AssignString(string) error {
	return mixins.BytesAssembler{TypeName: "bytes"}.AssignString("")
}
func (na *blob__Assembler) AssignBytes(v []byte) error {
	na.w = Blob{Bytes: v}
	return nil
}
func (blob__Assembler) AssignLink(ipld.Link) error {
	return mixins.BytesAssembler{TypeName: "bytes"}.AssignLink(nil)
}
func (na *blob__Assembler) AssignNode(v ipld.Node) error {
	if v2, err := v.AsBytes(); err != nil {
		return err
	} else {
		na.w = Blob{Bytes: v2}
		return nil
	}
}
func (blob__Assembler) Prototype() ipld.NodePrototype {
	return Prototype__Blob{}
}
