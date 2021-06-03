package ir

import (
	"fmt"
	"math/big"

	xr "github.com/libp2p/go-routing-language/syntax"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

type Number interface {
	TypeIsNumber()
}

type Int struct {
	*big.Int
	metadataCtx *meta.Meta
}

func NewInt64(v int64, metadata ...meta.Metadata) *Int {
	// Assemble metadata provided and update assemblyTime
	m := meta.New()
	if err := m.Apply(metadata...); err != nil {
		return &Int{big.NewInt(v), m}
	}
	return &Int{big.NewInt(v), m}
}

func (n *Int) Disassemble() xr.Node {
	return xr.Int{Int: n.Int}
}

func (n *Int) Metadata() meta.MetadataInfo {
	return n.metadataCtx.Get()
}

func (n *Int) TypeIsNumber() {}

func (n *Int) UpdateWith(ctx UpdateContext, with Node) error {
	wn, ok := with.(*Int)
	if !ok {
		return fmt.Errorf("cannot update with different primitive type")
	}
	// Update metadata
	n.metadataCtx.Update(wn.metadataCtx)
	return nil
}

type Float struct {
	*big.Float
	metadataCtx *meta.Meta
}

func (n *Float) Disassemble() xr.Node {
	return xr.Float{Float: n.Float}
}

func (n *Float) Metadata() meta.MetadataInfo {
	return n.metadataCtx.Get()
}

func (n *Float) TypeIsNumber() {}

func (n *Float) UpdateWith(ctx UpdateContext, with Node) error {
	wn, ok := with.(*Float)
	if !ok {
		return fmt.Errorf("cannot update with different primitive type")
	}
	// Update metadata
	n.metadataCtx.Update(wn.metadataCtx)
	return nil
}

func IsEqualNumber(x, y Number) bool {
	switch x1 := x.(type) {
	case *Int:
		switch y1 := y.(type) {
		case *Int:
			return x1.Int.Cmp(y1.Int) == 0
		case *Float:
			return false
		}
	case *Float:
		switch y1 := y.(type) {
		case *Int:
			return false
		case *Float:
			return x1.Float.Cmp(y1.Float) == 0
		}
	}
	panic("bug: unknown number type")
}
