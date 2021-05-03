package xr

import (
	"encoding/base64"
	"fmt"
	"io"
	"math/big"

	"github.com/ipld/go-ipld-prime"
	xrIpld "github.com/libp2p/go-smart-record/xr/ipld"
)

type Number interface {
	TypeIsNumber()
}

type Int struct {
	*big.Int
}

func NewInt64(v int64) Int {
	return Int{big.NewInt(v)}
}

func (n Int) TypeIsNumber() {}

func (n Int) WritePretty(w io.Writer) (err error) {
	_, err = w.Write([]byte(n.Int.String()))
	return err
}

type Float struct {
	*big.Float
}

func (n Float) TypeIsNumber() {}

func (n Float) WritePretty(w io.Writer) (err error) {
	_, err = w.Write([]byte(n.Float.String()))
	return err
}

func IsEqualNumber(x, y Number) bool {
	switch x1 := x.(type) {
	case Int:
		switch y1 := y.(type) {
		case Int:
			return x1.Int.Cmp(y1.Int) == 0
		case Float:
			return false
		}
	case Float:
		switch y1 := y.(type) {
		case Int:
			return false
		case Float:
			return x1.Float.Cmp(y1.Float) == 0
		}
	}
	panic("bug: unknown number type")
}

func (n Int) EncodeJSON() (interface{}, error) {
	bn, err := n.MarshalText()
	if err != nil {
		return nil, err
	}
	return struct {
		Type  marshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: IntType, Value: bn}, nil
}

func (n Float) EncodeJSON() (interface{}, error) {
	bn, err := n.MarshalText()
	if err != nil {
		return nil, err
	}
	return struct {
		Type  marshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: FloatType, Value: bn}, nil
}

func decodeInt(s map[string]interface{}) (Node, error) {
	z := new(big.Int)
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("wrong int decoding type")
	}
	// Unmarshaller inteprets []byte as string, we need to decode base64
	sDec, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, err
	}
	err = z.UnmarshalText(sDec)
	if err != nil {
		return nil, err
	}
	return Int{z}, nil
}

func decodeFloat(s map[string]interface{}) (Node, error) {
	z := new(big.Float)
	r, ok := s["value"].(string)
	if !ok {
		return nil, fmt.Errorf("wrong float decoding type")
	}
	// Unmarshaller inteprets []byte as string, we need to decode base64
	sDec, err := base64.StdEncoding.DecodeString(r)
	if err != nil {
		return nil, err
	}
	err = z.UnmarshalText(sDec)
	if err != nil {
		return nil, err
	}
	return Float{z}, nil
}

// ToIPLD converts xr.Node into its corresponding IPLD Node type
func (n Float) ToIPLD() (ipld.Node, error) {
	t := xrIpld.Type.Float_IPLD.NewBuilder()
	// NOTE: Disregarding accuracy
	f, _ := n.Float.Float64()
	err := t.AssignFloat(f)
	if err != nil {
		return nil, err
	}
	return t.Build(), nil
}

// ToIPLD converts xr.Node into its corresponding IPLD Node type
func (n Int) ToIPLD() (ipld.Node, error) {
	t := xrIpld.Type.Int_IPLD.NewBuilder()
	i := n.Int.Int64()
	err := t.AssignInt(i)
	if err != nil {
		return nil, err
	}
	return t.Build(), nil
}

// toNode_IPLD convert into IPLD Node of dynamic type NODE_IPLD
func (n Float) toNode_IPLD() (ipld.Node, error) {
	t := xrIpld.Type.Node_IPLD.NewBuilder()
	ma, err := t.BeginMap(-1)
	asm, err := ma.AssembleEntry("Float_IPLD")
	if err != nil {
		return nil, err
	}
	f, _ := n.Float.Float64()
	err = asm.AssignFloat(f)
	if err != nil {
		return nil, err
	}
	if err := ma.Finish(); err != nil {
		return nil, err
	}
	return t.Build(), nil
}

// toNode_IPLD convert into IPLD Node of dynamic type NODE_IPLD
func (n Int) toNode_IPLD() (ipld.Node, error) {
	t := xrIpld.Type.Node_IPLD.NewBuilder()
	ma, err := t.BeginMap(-1)
	asm, err := ma.AssembleEntry("Int_IPLD")
	if err != nil {
		return nil, err
	}
	err = asm.AssignInt(n.Int.Int64())
	if err != nil {
		return nil, err
	}
	if err := ma.Finish(); err != nil {
		return nil, err
	}
	return t.Build(), nil
}
