package ir

import (
	"fmt"
	"io"
)

type Bool struct {
	Value bool
}

func (b Bool) WritePretty(w io.Writer) (err error) {
	_, err = fmt.Fprintf(w, "%v", b.Value)
	return err
}

func (b Bool) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	w, ok := with.(Bool)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-bool")
	}
	return w, nil
}

func IsEqualBool(x, y Bool) bool {
	return x.Value == y.Value
}

func (b Bool) EncodeJSON() (interface{}, error) {
	return struct {
		Type  marshalType `json:"type"`
		Value bool        `json:"value"`
	}{Type: BoolType, Value: b.Value}, nil
}

func decodeBool(s map[string]interface{}) (Node, error) {
	r, ok := s["value"].(bool)
	if !ok {
		return nil, fmt.Errorf("decoded value not Bool")
	}
	return Bool{r}, nil
}
