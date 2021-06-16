package ir

import "testing"

// NOTE: For basic types, an update operation overwrites
// with the update value as long as it is the same type.
// In order to overwrite for another type a Set operation
// will be added to nodes.
func TestUpdateBasics(t *testing.T) {
	s1 := &String{Value: "test1"}
	s2 := &String{Value: "test2"}
	mctx := DefaultUpdateContext{}
	err := s1.UpdateWith(mctx, s2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(s1, s2) {
		t.Errorf("expecting %v, got %v", s2, s1)
	}

	b1 := &Bool{Value: true}
	b2 := &Bool{Value: false}
	err = b1.UpdateWith(mctx, b2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(b1, b2) {
		t.Errorf("expecting %v, got %v", s2, s1)
	}

	bt1 := &Bytes{Bytes: []byte("test1")}
	bt2 := &Bytes{Bytes: []byte("test2")}
	err = bt1.UpdateWith(mctx, bt2)
	if err != nil {
		t.Errorf("expecting no merge conflict, got %v", err)
	}
	if !IsEqual(bt1, bt2) {
		t.Errorf("expecting %v, got %v", bt2, bt1)
	}
}
