package ir

import (
	"encoding/json"
	"io"
	"math/big"
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
	z := new(big.Int)
	err := z.UnmarshalText(p)
	if err != nil {
		return err
	}
	n.Int = z
	return nil
}

func unmarshalFloat(p []byte, n *Float) error {
	z := new(big.Float)
	err := z.UnmarshalText(p)
	if err != nil {
		return err
	}
	n.Float = z
	return nil
}

func (n *Float) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		bn := []byte{}
		err = json.Unmarshal(b, &bn)
		if err != nil {
			return err
		}
		err = unmarshalFloat(bn, n)

		if err != nil {
			return err
		}
	} else {

		if _, ok := objMap["type"]; ok {
			c := struct {
				Type  MarshalType `json:"type"`
				Value []byte      `json:"value"`
			}{}
			err := json.Unmarshal(b, &c)
			unmarshalFloat(c.Value, n)
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
		bn := []byte{}
		err = json.Unmarshal(b, &bn)
		if err != nil {
			return err
		}
		err = unmarshalInt(bn, n)

		if err != nil {
			return err
		}
	} else {

		if _, ok := objMap["type"]; ok {
			c := struct {
				Type  MarshalType `json:"type"`
				Value []byte      `json:"value"`
			}{}
			err := json.Unmarshal(b, &c)
			unmarshalInt(c.Value, n)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func (n Int) MarshalJSON() (b []byte, e error) {
	bn, err := n.MarshalText()
	if err != nil {
		return nil, err
	}
	c := struct {
		Type  MarshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: IntType, Value: bn}
	return json.Marshal(&c)
}

func (n Float) MarshalJSON() (b []byte, e error) {
	bn, err := n.MarshalText()
	if err != nil {
		return nil, err
	}
	c := struct {
		Type  MarshalType `json:"type"`
		Value []byte      `json:"value"`
	}{Type: IntType, Value: bn}
	return json.Marshal(&c)
}
