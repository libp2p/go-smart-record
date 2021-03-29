package ir

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

func TestMarshale2e(t *testing.T) {
	n := Dict{
		Tag: "foo",
		Pairs: Pairs{
			{String{"bar"}, String{"baz"}},
			{Int{big.NewInt(567)}, String{"baz"}},
			{String{"bar22"}, Int{big.NewInt(567)}},
			{String{"bar2"}, Blob{[]byte("asdf")}},
			{Blob{[]byte("asdf")}, Int{big.NewInt(567)}},
			{Bool{true}, Int{big.NewInt(567)}},
		}}
	// Encode
	var b bytes.Buffer
	err := Marshal(&b, n)
	if err != nil {
		t.Fatal(err)
	}
	byteData := b.Bytes()

	r := bytes.NewReader(byteData)
	out := Dict{}
	err = Unmarshal(r, &out)
	if err != nil {
		t.Fatal(err)
	}

	var w bytes.Buffer
	n.WritePretty(&w)
	fmt.Println(w.String())
}

func TestMarshalString(t *testing.T) {
	n := String{"testing!"}
	var b bytes.Buffer
	err := Marshal(&b, n)
	if err != nil {
		t.Fatal(err)
	}
	o := String{}
	byteData := b.Bytes()
	r := bytes.NewReader(byteData)
	err = Unmarshal(r, &o)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling string", n, o)
	}
}

func TestMarshalBool(t *testing.T) {
	n := Bool{true}
	var b bytes.Buffer
	err := Marshal(&b, n)
	if err != nil {
		t.Fatal(err)
	}
	o := Bool{}
	byteData := b.Bytes()
	r := bytes.NewReader(byteData)
	err = Unmarshal(r, &o)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling string", n, o)
	}
}

func TestMarshalBlob(t *testing.T) {
	n := Blob{[]byte("testing!")}
	var b bytes.Buffer
	err := Marshal(&b, n)
	if err != nil {
		t.Fatal(err)
	}
	o := Blob{}
	byteData := b.Bytes()
	r := bytes.NewReader(byteData)
	err = Unmarshal(r, &o)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling blob", n, o)
	}
}

func TestMarshalNumber(t *testing.T) {
	n := Int{big.NewInt(123)}
	f := Float{big.NewFloat(123.123)}
	var b bytes.Buffer
	err := Marshal(&b, n)
	byteData := b.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	err = Marshal(&b, f)
	if err != nil {
		t.Fatal(err)
	}
	fbyteData := b.Bytes()
	o := Int{}
	of := Float{}
	r := bytes.NewReader(byteData)
	err = Unmarshal(r, &o)
	if err != nil {
		t.Fatal(err)
	}
	r = bytes.NewReader(fbyteData)
	err = Unmarshal(r, &of)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling Int", n, o)
	}
	if !IsEqual(f, of) {
		t.Fatal("Error unmarshalling Float", f, of)
	}

}

func TestMarshalPairs(t *testing.T) {
	n := Pairs{
		{String{"bar"}, String{"baz"}},
		{String{"bar2"}, String{"bar2123"}},
	}
	no := Pairs{}
	// Encode
	valueBytes, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(valueBytes, &no)
	if err != nil {
		panic(err)
	}
}
