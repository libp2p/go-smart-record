package base

import (
	"io"

	"github.com/libp2p/go-smart-record/ir"
)

// Record is a smart tag, representing data associated with a key.
type Record struct {
	// Key of the record in DHT-space (XOR-space).
	// This is NOT intended to be a CID or a PEER ID. Those are not in XOR-space.
	// The key is a function (usually hash) of the application level keys, e.g. CID or PEER ID.
	Key string

	// User holds user fields.
	User ir.Dict
}

func (r Record) Disassemble() ir.Dict {
	return r.User.CopySetTag("record", ir.String{"key"}, ir.String{r.Key})
}

func (r Record) WritePretty(w io.Writer) error {
	return r.Disassemble().WritePretty(w)
}

func (r Record) MergeWith(ctx ir.MergeContext, x ir.Node) ir.Node {
	panic("XXX")
}
