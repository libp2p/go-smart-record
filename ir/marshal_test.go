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
			{String{"bar1"}, String{"baz"}},
			{Int{big.NewInt(567)}, String{"baz"}},
			{String{"bar2"}, Int{big.NewInt(567)}},
			{String{"bar3"}, Blob{[]byte("asdf")}},
			{Blob{[]byte("asdf")}, Int{big.NewInt(567)}},
			{String{"bar4"}, Dict{
				Tag: "foo2",
				Pairs: Pairs{
					{Bool{true}, Int{big.NewInt(567)}},
				},
			}},
		},
	}

	// Encode
	var b bytes.Buffer
	err := Marshal(&b, n)
	if err != nil {
		t.Fatal(err)
	}
	byteData := b.Bytes()

	// Decode
	r := bytes.NewReader(byteData)
	out := Dict{}
	err = Unmarshal(r, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, out) {
		fmt.Println(IsEqual(n, out))
		fmt.Println("== IN ==")
		var w bytes.Buffer
		n.WritePretty(&w)
		fmt.Println(w.String())

		fmt.Println("== OUT ==")
		w.Reset()
		out.WritePretty(&w)
		fmt.Println(w.String())
		t.Fatal("Error unmarshalling Dict")
	}
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

	n = Bool{false}
	err = Marshal(&b, n)
	if err != nil {
		t.Fatal(err)
	}
	o = Bool{}
	byteData = b.Bytes()
	r = bytes.NewReader(byteData)
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
	var b bytes.Buffer
	err := Marshal(&b, n)
	byteData := b.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	o := Int{}
	r := bytes.NewReader(byteData)
	err = Unmarshal(r, &o)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling Int", n, o)
	}

	// We must use 64 precision to perform the right comparison.
	// UnmarshalText generates a 64 precision float.
	// Check: https://github.com/golang/go/issues/45309
	f := Float{big.NewFloat(123.123).SetPrec(64)}
	of := Float{}
	err = Marshal(&b, f)
	if err != nil {
		t.Fatal(err)
	}
	fbyteData := b.Bytes()
	r = bytes.NewReader(fbyteData)
	err = Unmarshal(r, &of)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(f, of) {
		t.Fatal("Error unmarshalling Float", f, of)
	}

}

func TestMarshalPairs(t *testing.T) {
	n := Pairs{
		{Blob{[]byte("asdf")}, Int{big.NewInt(567)}},
		{String{"bar"}, String{"baz"}},
		{String{"bar2"}, String{"bar2123"}},
		{Bool{true}, Int{big.NewInt(567)}},
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
	if !AreSamePairs(n, no) {
		t.Fatal("Error marshalling pairs")
	}
}

func TestMarshalDict(t *testing.T) {
	n := Dict{
		Tag: "foo2",
		Pairs: Pairs{
			{Bool{true}, Int{big.NewInt(567)}},
		},
	}
	// Encode
	var b bytes.Buffer
	err := Marshal(&b, n)
	if err != nil {
		t.Fatal(err)
	}
	byteData := b.Bytes()

	// Decode
	r := bytes.NewReader(byteData)
	out := Dict{}
	err = Unmarshal(r, &out)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, out) {
		var w bytes.Buffer
		out.WritePretty(&w)
		fmt.Println(w.String())
		t.Fatal("Error unmarshalling Dict")
	}

}
