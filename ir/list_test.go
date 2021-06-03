package ir

import (
	"testing"
)

func TestUpdateListSameTag(t *testing.T) {
	s1 := &List{
		Elements: Nodes{
			&String{"x", nil},
			&String{"z", nil},
		},
	}
	s2 := &List{
		Elements: Nodes{
			&String{"x", nil},
			&String{"w", nil},
		},
	}
	exp := &List{
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
