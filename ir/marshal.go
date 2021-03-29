package ir

import (
	"encoding/json"
	"fmt"
	"io"
)

type MarshalType int

const (
	StringType = iota
	BlobType
	PeerType
	IntType
	FloatType
	BoolType
)

func UnmarshalType(tp MarshalType, b []byte) (Node, error) {
	switch tp {
	case StringType:
		var n String
		err := json.Unmarshal(b, &n)
		if err != nil {
			return nil, err
		}
		return n, nil
	case BlobType:
		var n Blob
		err := json.Unmarshal(b, &n)
		if err != nil {
			return nil, err
		}
		return n, nil
	case IntType:
		var n Int
		err := json.Unmarshal(b, &n)
		if err != nil {
			return nil, err
		}
		return n, nil
	case FloatType:
		var n Float
		err := json.Unmarshal(b, &n)
		if err != nil {
			return nil, err
		}
		return n, nil
	case BoolType:
		var n Bool
		err := json.Unmarshal(b, &n)
		if err != nil {
			return nil, err
		}
		return n, nil
	}
	return nil, fmt.Errorf("Wrong type")
}

func Marshal(w io.Writer, in Node) error {
	enc := json.NewEncoder(w)
	return enc.Encode(in)
}

func Unmarshal(r io.Reader, out Node) error {
	dec := json.NewDecoder(r)
	for {
		if err := dec.Decode(out); err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
	}
}
