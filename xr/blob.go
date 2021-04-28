package xr

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
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

func (b Blob) EncodeJSON() (interface{}, error) {
	return struct {
		Type  marshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: BlobType, Value: b.Bytes}, nil
}

func IsEqualBlob(x, y Blob) bool {
	return bytes.Compare(x.Bytes, y.Bytes) == 0
}

func decodeBlob(s map[string]interface{}) (Node, error) {
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("decoding typ is not Blob")
	}
	// Unmarshaller inteprets []byte as string, we need to decode base64
	sDec, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, err
	}
	return Blob{sDec}, nil
}
