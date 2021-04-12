package ir

import (
	"fmt"
	"io"
)

// String is a node representing a string literal.
type String struct {
	Value string
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
