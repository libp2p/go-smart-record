package ir

import (
	"io"
)

type Cid struct {
	Cid string
	// User holds user fields.
	User Dict
}

func (c Cid) Dict() Dict {
	return c.User.CopySetTag("cid", String{c.Cid}, String{c.Cid})
}

func (c Cid) WritePretty(w io.Writer) error {
	return c.Dict().WritePretty(w)
}
