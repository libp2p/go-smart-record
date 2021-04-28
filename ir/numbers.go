package ir

import (
	"fmt"
	"io"
	"math/big"

	"github.com/libp2p/go-smart-record/xr"
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

func (n Int) Disassemble() xr.Node {
	return xr.Int{Int: n.Int}
}

func (n Int) TypeIsNumber() {}

func (n Int) WritePretty(w io.Writer) (err error) {
	_, err = w.Write([]byte(n.Int.String()))
	return err
}

func (n Int) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	wn, ok := with.(Int)
	if !ok {
		return nil, fmt.Errorf("cannot update with different primitive type")
	}
	return wn, nil
}

type Float struct {
	*big.Float
}

func (n Float) Disassemble() xr.Node {
	return xr.Float{Float: n.Float}
}

func (n Float) TypeIsNumber() {}

func (n Float) UpdateWith(ctx UpdateContext, with Node) (Node, error) {
	wn, ok := with.(Float)
	if !ok {
		return nil, fmt.Errorf("cannot update with different primitive type")
	}
	return wn, nil
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
