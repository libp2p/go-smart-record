package ir

import (
	"testing"
)

func TestUpdateSetDiffTag(t *testing.T) {
	s1 := Set{
		Tag:      "aaa",
		Elements: Nodes{},
	}
	s2 := Set{
		Tag:      "bbb",
		Elements: Nodes{},
	}
	mctx := DefaultUpdateContext{}
	if _, err := Update(mctx, s1, s2); err != nil {
		t.Errorf("update (%v)", err)
	}
}

func TestUpdateSetSameTag(t *testing.T) {
	s1 := Set{
		Tag: "aaa",
		Elements: Nodes{
			String{"x"},
			String{"z"},
		},
	}
	s2 := Set{
		Tag: "aaa",
		Elements: Nodes{
			String{"x"},
			String{"w"},
		},
	}
	exp := Set{
		Tag: "aaa",
		Elements: Nodes{
			String{"x"},
			String{"z"},
			String{"w"},
		},
	}
	mctx := DefaultUpdateContext{}
	m, err := Update(mctx, s1, s2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(m, exp) {
		t.Errorf("expecting %v, got %v", exp, m)
	}
}
