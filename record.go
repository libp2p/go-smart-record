package sr

// Record is a smart tag, representing data associated with a key.
type Record struct {
	// Key of the record in DHT-space (XOR-space).
	// This is NOT intended to be a CID or a PEER ID. Those are not in XOR-space.
	// The key is a function (usually hash) of the application level keys, e.g. CID or PEER ID.
	Key string

	// Dict holds the non-key fields of the record.
	Dict
}

func (r Record) AsDict() Dict {
	return r.Dict.CopySet(String{"key"}, String{r.Key})
}

func MergeRecords(x, y *Record) Node {
	panic("XXX")
}
