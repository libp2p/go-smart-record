package ir

import (
	"bytes"
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
	b, err := Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	// Decode
	out, err := Unmarshal(b)
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
	b, err := Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	o, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}

	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling string", n, o)
	}
}

func TestMarshalBool(t *testing.T) {
	n := Bool{true}
	b, err := Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	o, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling string", n, o)
	}

	n = Bool{false}
	b, err = Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	o, err = Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling string", n, o)
	}
}

func TestMarshalBlob(t *testing.T) {
	n := Blob{[]byte("testing!")}
	b, err := Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	o, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, o) {
		t.Fatal("Error unmarshalling blob", n, o)
	}
}

func TestMarshalNumber(t *testing.T) {
	n := Int{big.NewInt(123)}
	b, err := Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	o, err := Unmarshal(b)
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
	b, err = Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	of, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(f, of) {
		t.Fatal("Error unmarshalling Float", f, of)
	}

}

func TestMarshalDict(t *testing.T) {
	n := Dict{
		Tag: "foo2",
		Pairs: Pairs{
			{Bool{true}, Int{big.NewInt(567)}},
		},
	}
	b, err := Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, out) {
		t.Fatal("Error unmarshalling Dict")
	}

}

func TestMarshalSet(t *testing.T) {
	n1 := Bool{true}
	n2 := String{"testing!"}
	n3 := Int{big.NewInt(567)}

	n := Set{
		Tag:      "foo2",
		Elements: []Node{n1, n2, n3},
	}

	b, err := Marshal(n)
	if err != nil {
		t.Fatal(err)
	}
	out, err := Unmarshal(b)
	if err != nil {
		t.Fatal(err)
	}
	if !IsEqual(n, out) {
		t.Fatal("Error unmarshalling Set")
	}

}
