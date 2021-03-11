package ir

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPretty(t *testing.T) {
	n := Dict{
		Tag: "foo",
		Pairs: Pairs{
			{String{"bar"}, String{"baz"}},
			{String{"bar2"}, String{"baz"}},
		},
	}
	var w bytes.Buffer
	n.WritePretty(&w)
	fmt.Println(w.String())
}