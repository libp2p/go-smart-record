package ir

import (
	"encoding/json"
	"fmt"
)

type marshalType string

// List of syntactic types supported
const (
	StringType = "String"
	BlobType   = "Blob"
	IntType    = "Int"
	FloatType  = "Float"
	BoolType   = "Bool"
	DictType   = "Dict"
	SetType    = "Set"
)

// decodeNode does the unmarshalling of the type
// once the wrapper has been process and the Node type has
// been identified.
func decodeNode(v interface{}) (Node, error) {
	s, ok := v.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("bad decoding format")
	}

	switch s["type"] {
	case StringType:
		return decodeString(s)
	case BlobType:
		return decodeBlob(s)
	case IntType:
		return decodeInt(s)
	case FloatType:
		return decodeFloat(s)
	case BoolType:
		return decodeBool(s)
	case DictType:
		return decodeDict(s)
	case SetType:
		return decodeSet(s)
	}
	return nil, fmt.Errorf("Wrong type")
}

// Marshal syntactic representation
func Marshal(n Node) ([]byte, error) {
	c, err := n.EncodeJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(c)
}

// Unmarshal syntactic representation
func Unmarshal(r []byte) (Node, error) {
	var v map[string]interface{}
	json.Unmarshal(r, &v)
	return decodeNode(v)
}
