package xr

import (
	"bytes"
	"math/big"
	"testing"

	cbor "github.com/ipld/go-ipld-prime/codec/dagcbor"
)

func TestBlobNode(t *testing.T) {
	n := Blob{[]byte("testing")}
	var buf bytes.Buffer
	err := cbor.Encode(n, &buf)
	if err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()

	no := Prototype__Blob{}.NewBuilder()
	err = cbor.Decode(no, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	//TODO: Remove casting, use xr.Node
	if !IsEqual(n, no.Build().(Blob)) {
		t.Fatal("Marshalled Blob nodes not equal")
	}

}

func TestBoolNode(t *testing.T) {
	n := Bool{true}
	var buf bytes.Buffer
	err := cbor.Encode(n, &buf)
	if err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()

	no := Prototype__Bool{}.NewBuilder()
	err = cbor.Decode(no, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	//TODO: Remove casting, use xr.Node
	if !IsEqual(n, no.Build().(Bool)) {
		t.Fatal("Marshalled Blob nodes not equal")
	}

}

func TestStringNode(t *testing.T) {
	n := String{"testing"}
	var buf bytes.Buffer
	err := cbor.Encode(n, &buf)
	if err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()

	no := Prototype__String{}.NewBuilder()
	err = cbor.Decode(no, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	//TODO: Remove casting, use xr.Node
	if !IsEqual(n, no.Build().(String)) {
		t.Fatal("Marshalled String nodes not equal")
	}
}

func TestIntNode(t *testing.T) {
	n := NewInt64(123)
	var buf bytes.Buffer
	err := cbor.Encode(n, &buf)
	if err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()

	no := Prototype__Int{}.NewBuilder()
	err = cbor.Decode(no, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	//TODO: Remove casting, use xr.Node
	if !IsEqual(n, no.Build().(Int)) {
		t.Fatal("Marshalled Int nodes not equal")
	}
}

func TestFloatNode(t *testing.T) {
	n := Float{big.NewFloat(123.332).SetPrec(64)}
	var buf bytes.Buffer
	err := cbor.Encode(n, &buf)
	if err != nil {
		t.Fatal(err)
	}
	b := buf.Bytes()

	no := Prototype__Float{}.NewBuilder()
	err = cbor.Decode(no, bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	//TODO: Remove casting, use xr.Node
	if !IsEqual(n, no.Build().(Float)) {
		t.Fatal("Marshalled Float nodes not equal")
	}
}
