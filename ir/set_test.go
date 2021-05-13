package ir

import (
	"testing"
)

func TestUpdateSetDiffTag(t *testing.T) {
	s1 := &Set{
		Tag:      "aaa",
		Elements: Nodes{},
	}
	s2 := &Set{
		Tag:      "bbb",
		Elements: Nodes{},
	}
	mctx := DefaultUpdateContext{}
	if err := Update(mctx, s1, s2); err != nil {
		t.Errorf("update (%v)", err)
	}
}

func TestUpdateSetSameTag(t *testing.T) {
	s1 := &Set{
		Tag: "aaa",
		Elements: Nodes{
			&String{"x", nil},
			&String{"z", nil},
		},
	}
	s2 := &Set{
		Tag: "aaa",
		Elements: Nodes{
			&String{"x", nil},
			&String{"w", nil},
		},
	}
	exp := &Set{
		Tag: "aaa",
		Elements: Nodes{
			&String{"x", nil},
			&String{"z", nil},
			&String{"w", nil},
		},
	}
	mctx := DefaultUpdateContext{}
	err := Update(mctx, s1, s2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(s1, exp) {
		t.Errorf("expecting %v, got %v", exp, s1)
	}
}
