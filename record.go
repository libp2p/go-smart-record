package sr

import (
	"io"
)

type Record struct {
	// Key of the record in DHT-space (XOR-space).
	// This is NOT intended to be a CID or a PEER ID. These are not in XOR-space.
	// The key is a function (usually hash) of the application level keys, e.g. CID or PEER ID.
	// Key corresponds to the XML attribute "key".
	Key string

	// Dict holds the named children of the record.
	Dict
}

func (r Record) WritePretty(w io.Writer) error {
	panic("XXX")
}

func MergeRecords(x, y *Record) Node {
	panic("XXX")
}
