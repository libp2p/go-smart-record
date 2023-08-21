package vm

import (
	"encoding/json"
	"fmt"

	"github.com/libp2p/go-libp2p/core/peer"
	xr "github.com/libp2p/go-routing-language/syntax"
)

// MarshalRecordValue serializes RecordValue to send it through libp2p protocol.
func MarshalRecordValue(r RecordValue) ([]byte, error) {
	out := make(map[string][]byte)
	for k, v := range r {
		n, err := xr.MarshalJSON(v)
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
	out := make(map[peer.ID]*xr.Dict)
	for k, p := range unm {
		n, err := xr.UnmarshalJSON(p)
		if err != nil {
			return nil, err
		}
		no, ok := n.(xr.Dict)
		if !ok {
			return nil, fmt.Errorf("no dict type unmarshalling RecordValue item")
		}
		pid, err := peer.Decode(k)
		if err != nil {
			return nil, err
		}
		out[pid] = &no
	}
	return out, nil
}
