package vm

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-smart-record/ir"
)

// RecordValue determines the structure of data stored in a record.
// Each peer has a private dataspace to store dicts.
type RecordValue map[peer.ID]*ir.Dict

type Record struct {
	Key   string      // Key of the record
	Value RecordValue // Root stores peers dictionaries in that key.
}

// Machine captures the public interface of a smart record virtual machine.
type Machine interface {
	Update(writer peer.ID, k string, update ir.Dict) error // Updates the dictionary in the writer's private space.
	Get(k string) Record                                   // Get the full Record in a key
	// NOTE: No query operation will be supported until we figure out selectors
	// Query(key string, selector Selector) (RecordValue, error)
}

// VM implements the Machine interface and keeps the map of records in its state.
type vm struct {
	ctx ir.UpdateContext // UpdateContext the VM uses to resolve conflicts
	//ds  ds.Datastore    // TODO: Add a datastore instead of using map[string] for the VM state
	keys map[string]*Record // State of the VM storing the map of records.
	asm  ir.Assembler       // Assemble to use in the VM.
	lk   sync.RWMutex       // Lock to enable multiple access
}

// NewVM creates a new smart record Machine
func NewVM(ctx ir.UpdateContext, asm ir.Assembler) Machine {
	return newVM(ctx, asm)
}

//newVM instantiates a new VM with an updateContext and an assembler
func newVM(ctx ir.UpdateContext, asm ir.Assembler) *vm {
	return &vm{
		ctx:  ctx,
		keys: make(map[string]*Record),
		asm:  asm,
	}
}

// Get the whole record stored in a key
func (v *vm) Get(k string) Record {
	// TODO: Implementation
	return Record{Key: "Sample", Value: RecordValue{}}
}

// Update the dictionary in the writer's private space
func (v *vm) Update(writer peer.ID, k string, update ir.Dict) error {
	// TODO: Implementation
	return nil
}

func MarshalRecordValue(r Record) ([]byte, error) {
	out := make(map[string][]byte)
	for k, v := range r.Value {
		n, err := ir.Marshal(v)
		if err != nil {
			return nil, err
		}
		out[k.String()] = n
	}
	return json.Marshal(out)
}

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
