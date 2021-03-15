package ir

import (
	"bytes"
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

func IsEqualBlob(x, y Blob) bool {
	return bytes.Compare(x.Bytes, y.Bytes) == 0
}
