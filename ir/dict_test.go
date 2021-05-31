package ir

import (
	"testing"
)

func TestUpdateDictDiffTag(t *testing.T) {
	d1 := &Dict{
		Pairs: Pairs{},
	}
	d2 := &Dict{
		Pairs: Pairs{},
	}
	mctx := DefaultUpdateContext{}
	if err := Update(mctx, d1, d2); err != nil {
		t.Errorf("update (%v)", err)
	}
}

func TestUpdateDictDisjointPairs(t *testing.T) {
	d1 := &Dict{
		Pairs: Pairs{{&String{"x", nil}, NewInt64(1)}},
	}
	d2 := &Dict{
		Pairs: Pairs{{&String{"y", nil}, NewInt64(1)}},
	}
	exp := &Dict{
		Pairs: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"y", nil}, NewInt64(1)},
		},
	}
	mctx := DefaultUpdateContext{}
	err := Update(mctx, d1, d2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(d1, exp) {
		t.Errorf("expecting %v, got %v", exp, d1)
	}
}

func TestUpdateDictOverlappingPairs(t *testing.T) {
	d1 := &Dict{
		Pairs: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"z", nil}, NewInt64(1)},
		},
	}
	d2 := &Dict{
		Pairs: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"w", nil}, NewInt64(1)},
		},
	}
	exp := &Dict{
		Pairs: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"z", nil}, NewInt64(1)},
			{&String{"w", nil}, NewInt64(1)},
		},
	}
	mctx := DefaultUpdateContext{}
	err := Update(mctx, d1, d2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(d1, exp) {
		t.Errorf("expecting %v, got %v", exp, d1)
	}
}
