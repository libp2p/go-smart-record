package ir

import (
	"fmt"
	"io"
)

// String is a node representing a string literal.
type String struct {
	Value string
}

func (s String) Disassemble() Node {
	return s
}

func (s String) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "%q", s.Value)
	return err
}

func IsEqualString(x, y String) bool {
	return x.Value == y.Value
}

func (s String) EncodeJSON() (interface{}, error) {
	return struct {
		Type  marshalType `json:"type"`
		Value string      `json:"value"`
	}{Type: StringType, Value: s.Value}, nil
}

func decodeString(s map[string]interface{}) (Node, error) {
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("decoded value not String")
	}
	return String{r}, nil
}

func (s String) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	w, ok := with.(String)
	if !ok {
		return nil, fmt.Errorf("cannot update with a non-string")
	}
	return w, nil
}
