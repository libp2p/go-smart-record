package ir

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"
)

func TestPretty(t *testing.T) {
	n := Dict{
		Tag: "foo",
		Pairs: Pairs{
			{String{"bar"}, String{"baz"}},
			{String{"bar2"}, Int{big.NewInt(567)}},
		},
	}
	var w bytes.Buffer
	n.WritePretty(&w)
	fmt.Println(w.String())
}
