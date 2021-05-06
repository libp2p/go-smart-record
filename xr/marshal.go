package xr

import (
	"bytes"

	cbor "github.com/ipld/go-ipld-prime/codec/dagcbor"
	json "github.com/ipld/go-ipld-prime/codec/dagjson"
	xrIpld "github.com/libp2p/go-smart-record/xr/ipld"
)

// Marshal syntactic representation
func Marshal(n Node) ([]byte, error) {
	in, err := n.toNode_IPLD()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Encode(in, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal syntactic representation
func Unmarshal(r []byte) (Node, error) {
	n := xrIpld.Type.Node_IPLD.NewBuilder()
	err := json.Decode(n, bytes.NewReader(r))
	if err != nil {
		return nil, err
	}
	return FromIPLD(n.Build())
}

// Encode Serializes syntactic nodes in CBOR using its IPLD capabilities
func Encode(n Node) ([]byte, error) {
	in, err := n.toNode_IPLD()
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = cbor.Encode(in, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode de-serializes syntactic nodes in CBOR using its IPLD capabilities
func Decode(r []byte) (Node, error) {
	n := xrIpld.Type.Node_IPLD.NewBuilder()
	err := cbor.Decode(n, bytes.NewReader(r))
	if err != nil {
		return nil, err
	}
	return FromIPLD(n.Build())
}
