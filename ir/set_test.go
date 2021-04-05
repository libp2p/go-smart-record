package ir

import (
	"testing"
)

func TestMergeSetDiffTag(t *testing.T) {
	s1 := Set{
		Tag:      "aaa",
		Elements: Nodes{},
	}
	s2 := Set{
		Tag:      "bbb",
		Elements: Nodes{},
	}
	mctx := DefaultMergeContext{}
	if _, err := Merge(mctx, s1, s2); err == nil {
		t.Errorf("expecting a merge conflict")
	}
}

func TestMergeSetSameTag(t *testing.T) {
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
	mctx := DefaultMergeContext{}
	m, err := Merge(mctx, s1, s2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(m, exp) {
		t.Errorf("expecting %v, got %v", exp, m)
	}
}
