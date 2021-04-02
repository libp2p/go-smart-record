package base

import (
	"fmt"
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

func (r Record) MergeWith(ctx ir.MergeContext, x ir.Node) (ir.Node, error) {
	q, ok := x.(Record)
	if !ok {
		return nil, fmt.Errorf("cannot merge record with a non-record")
	}
	if q.Key != r.Key {
		return nil, fmt.Errorf("cannot merge records with different keys")
	}
	u, err := ir.MergeDictIgnoreTag(ctx, r.User, q.User)
	if err != nil {
		return nil, fmt.Errorf("user keys unmergable (%v)", err)
	}
	w := u.(ir.Dict)
	w.Tag = ""
	return Record{
		Key:  r.Key,
		User: w,
	}, nil
}

type RecordAssembler struct{}

func (RecordAssembler) Assemble(ctx ir.AssemblerContext, src ir.Dict) (ir.Node, error) {
	if src.Tag != "record" {
		return nil, fmt.Errorf("not a record tag")
	}
	n := src.Get(ir.String{"key"})
	if n == nil {
		return nil, fmt.Errorf("record without a key")
	}
	k, ok := n.(ir.String)
	if !ok {
		return nil, fmt.Errorf("record key is not a string")
	}
	u := src.Copy()
	u.Tag = ""
	u.Remove(ir.String{"key"})
	return Record{
		Key:  k.Value,
		User: u,
	}, nil
}
