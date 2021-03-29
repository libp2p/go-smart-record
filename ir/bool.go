package ir

import (
	"encoding/json"
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

func IsEqualBool(x, y Bool) bool {
	return x.Value == y.Value
}

func (b Bool) MarshalJSON() (byteData []byte, e error) {
	// Temporal type to avoid recursion
	type tmp Bool
	ts := tmp(b)

	c := struct {
		Type  MarshalType `json:"type"`
		Value tmp         `json:"value"`
	}{Type: BoolType, Value: ts}
	return json.Marshal(&c)
}

func (b *Bool) UnmarshalJSON(byteData []byte) error {
	type tmp Bool
	ts := tmp(*b)

	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(byteData, &objMap)
	if err != nil {
		return err
	}

	if _, ok := objMap["type"]; !ok {
		err = json.Unmarshal(byteData, &ts)
	} else {
		err = json.Unmarshal(*objMap["value"], &ts)
	}
	if err != nil {
		return err
	}
	*b = Bool(ts)
	return nil
}
