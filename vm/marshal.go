package vm

import (
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-smart-record/ir"
)

// MarshalRecordValue serializes RecordValue to send it through libp2p protocol.
func MarshalRecordValue(r RecordValue) ([]byte, error) {
	out := make(map[string][]byte)
	for k, v := range r {
		n, err := ir.Marshal(v)
		if err != nil {
			return nil, err
		}
		out[k.String()] = n
	}
	return json.Marshal(out)
}

// UnmarshalRecordValue unmarshals a serialized representation of RecordValue
func UnmarshalRecordValue(b []byte) (RecordValue, error) {
	unm := make(map[string][]byte)
	err := json.Unmarshal(b, &unm)
	if err != nil {
		return nil, err
	}
	out := make(map[peer.ID]*ir.Dict)
	for k, v := range unm {
		n, err := ir.Unmarshal(v)
		if err != nil {
			return nil, err
		}
		no, ok := n.(ir.Dict)
		if !ok {
			return nil, fmt.Errorf("no dict type unmarshalling RecordValue item")
		}
		pid, err := peer.IDFromString(k)
		if err != nil {
			return nil, err
		}
		out[pid] = &no
	}
	return out, nil
}
