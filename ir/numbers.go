package ir

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
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

func unmarshalInt(p []byte, n *Int) error {
	if string(p) == "null" {
		return nil
	}
	var z big.Int
	_, ok := z.SetString(strings.Replace(string(p), "\"", "", -1), 10)
	if !ok {
		return fmt.Errorf("not a valid big integer int: %s", p)
	}
	n.Int = &z
	return nil
}

func unmarshalFloat(p []byte, n *Float) error {
	if string(p) == "null" {
		return nil
	}
	var z big.Float
	_, ok := z.SetString(strings.Replace(string(p), "\"", "", -1))
	if !ok {
		return fmt.Errorf("not a valid big integer float: %s", p)
	}
	n.Float = &z
	return nil
}
func (n *Float) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		err = unmarshalFloat(b, n)

		if err != nil {
			return err
		}
	} else {

		if _, ok := objMap["type"]; ok {
			err = unmarshalFloat(*objMap["value"], n)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (n *Int) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		err = unmarshalInt(b, n)

		if err != nil {
			return err
		}
	} else {

		if _, ok := objMap["type"]; ok {
			err = unmarshalInt(*objMap["value"], n)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (n Int) MarshalJSON() (b []byte, e error) {
	c := struct {
		Type  MarshalType `json:"type"`
		Value string      `json:"value"`
	}{Type: IntType, Value: n.String()}
	return json.Marshal(&c)
}

func (n Float) MarshalJSON() (b []byte, e error) {
	c := struct {
		Type  MarshalType `json:"type"`
		Value string      `json:"value"`
	}{Type: FloatType, Value: n.String()}
	return json.Marshal(&c)
}
