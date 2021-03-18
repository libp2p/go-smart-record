package ir

import (
	"encoding/json"
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

func (s String) MarshalJSON() (b []byte, e error) {
	// Temporal type to avoid recursion
	type tmp String
	ts := tmp(s)

	c := struct {
		Type  MarshalType `json:"type"`
		Value tmp         `json:"value"`
	}{Type: StringType, Value: ts}
	return json.Marshal(&c)
}

func (s *String) UnmarshalJSON(b []byte) error {
	type tmp String
	ts := tmp(*s)

	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	if _, ok := objMap["type"]; !ok {
		err = json.Unmarshal(b, &ts)
	} else {
		err = json.Unmarshal(*objMap["value"], &ts)
	}
	if err != nil {
		return err
	}
	*s = String(ts)
	return nil
}
