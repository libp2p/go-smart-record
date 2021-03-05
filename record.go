package cr

import (
	"encoding/xml"
)

type Record {
	// Update applies when the record structure is used to build an update command.
	Update

	// Key of the record in DHT-space (XOR-space).
	// This is NOT intended to be a CID or a PEER ID. These are not in XOR-space.
	// The key is a function (usually hash) of the application level keys, e.g. CID or PEER ID.
	// Key corresponds to the XML attribute "key".
	Key string

	// Dict holds the named children of the record.
	Dict
}

func (r Record) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return XXX
}

func MergeRecords(x, y *Record) Node {
	XXX
}
