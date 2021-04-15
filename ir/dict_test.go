package ir

import (
	"testing"
)

func TestUpdateDictDiffTag(t *testing.T) {
	d1 := Dict{
		Tag:   "aaa",
		Pairs: Pairs{},
	}
	d2 := Dict{
		Tag:   "bbb",
		Pairs: Pairs{},
	}
	mctx := DefaultUpdateContext{}
	if _, err := Update(mctx, d1, d2); err != nil {
		t.Errorf("update (%v)", err)
	}
}

func TestUpdateDictDisjointPairs(t *testing.T) {
	d1 := Dict{
		Tag:   "aaa",
		Pairs: Pairs{{String{"x"}, NewInt64(1)}},
	}
	d2 := Dict{
		Tag:   "aaa",
		Pairs: Pairs{{String{"y"}, NewInt64(1)}},
	}
	exp := Dict{
		Tag: "aaa",
		Pairs: Pairs{
			{String{"x"}, NewInt64(1)},
			{String{"y"}, NewInt64(1)},
		},
	}
	mctx := DefaultUpdateContext{}
	m, err := Update(mctx, d1, d2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(m, exp) {
		t.Errorf("expecting %v, got %v", exp, m)
	}
}

func TestUpdateDictOverlappingPairs(t *testing.T) {
	d1 := Dict{
		Tag: "aaa",
		Pairs: Pairs{
			{String{"x"}, NewInt64(1)},
			{String{"z"}, NewInt64(1)},
		},
	}
	d2 := Dict{
		Tag: "aaa",
		Pairs: Pairs{
			{String{"x"}, NewInt64(1)},
			{String{"w"}, NewInt64(1)},
		},
	}
	exp := Dict{
		Tag: "aaa",
		Pairs: Pairs{
			{String{"x"}, NewInt64(1)},
			{String{"z"}, NewInt64(1)},
			{String{"w"}, NewInt64(1)},
		},
	}
	mctx := DefaultUpdateContext{}
	m, err := Update(mctx, d1, d2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(m, exp) {
		t.Errorf("expecting %v, got %v", exp, m)
	}
}
