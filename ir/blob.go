package ir

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

type Blob struct {
	Bytes []byte
}

func (b Blob) WritePretty(w io.Writer) error {
	_, err := fmt.Fprintf(w, "0x%s", hex.EncodeToString(b.Bytes)) // TODO: We can do better. E.g. wrap on 80-column boundary.
	return err
}

func (b Blob) MarshalJSON() (bdata []byte, e error) {
	// Temporal type to avoid recursion
	type tmp Blob
	ts := tmp(b)

	c := struct {
		Type  MarshalType `json:"type"`
		Value tmp         `json:"value"`
	}{Type: BlobType, Value: ts}
	return json.Marshal(&c)
}

func IsEqualBlob(x, y Blob) bool {
	return bytes.Compare(x.Bytes, y.Bytes) == 0
}

func (b *Blob) UnmarshalJSON(bdata []byte) error {
	type tmp Blob
	ts := tmp(*b)

	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(bdata, &objMap)
	if err != nil {
		return err
	}

	if _, ok := objMap["type"]; !ok {
		err = json.Unmarshal(bdata, &ts)
	} else {
		err = json.Unmarshal(*objMap["value"], &ts)
	}
	if err != nil {
		return err
	}
	*b = Blob(ts)
	return nil
}
