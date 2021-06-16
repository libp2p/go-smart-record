package ir

import (
	"testing"
)

func TestUpdatePredicateDiffTag(t *testing.T) {
	d1 := &Predicate{
		Tag: "tag1",
	}
	d2 := &Predicate{
		Tag: "tag1",
	}
	mctx := DefaultUpdateContext{}
	if err := Update(mctx, d1, d2); err != nil {
		t.Errorf("update (%v)", err)
	}
}

func TestUpdatePredicateDisjoint(t *testing.T) {
	d1 := &Predicate{
		Tag:   "test",
		Named: Pairs{{&String{"x", nil}, NewInt64(1)}},
		Positional: Nodes{
			&String{"x", nil},
			&String{"z", nil},
		},
	}
	d2 := &Predicate{
		Tag:   "test",
		Named: Pairs{{&String{"y", nil}, NewInt64(1)}},
		Positional: Nodes{
			&String{"w", nil},
		},
	}
	exp := &Predicate{
		Tag: "test",
		Named: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"y", nil}, NewInt64(1)},
		},
		Positional: Nodes{
			&String{"x", nil},
			&String{"z", nil},
			&String{"w", nil},
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

func TestUpdatePredicateOverlapping(t *testing.T) {
	d1 := &Predicate{
		Tag: "test",
		Named: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"z", nil}, NewInt64(1)},
		},
		Positional: Nodes{
			&String{"x", nil},
			&String{"z", nil},
		},
	}
	d2 := &Predicate{
		Tag: "test",
		Named: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"w", nil}, NewInt64(1)},
		},
		Positional: Nodes{
			&String{"x", nil},
			&String{"w", nil},
		},
	}
	exp := &Predicate{
		Tag: "test",
		Named: Pairs{
			{&String{"x", nil}, NewInt64(1)},
			{&String{"z", nil}, NewInt64(1)},
			{&String{"w", nil}, NewInt64(1)},
		},
		Positional: Nodes{
			&String{"x", nil},
			&String{"z", nil},
			&String{"w", nil},
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
