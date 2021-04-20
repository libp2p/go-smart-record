// Package vm defines and implements a Virtual Machine for smart records.
package vm

import (
	"sync"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-smart-record/ir"
)

// RecordValue determines the structure of data stored in a record.
// Each peer has a private dataspace to store dicts.
type RecordValue map[peer.ID]*ir.Dict

// Machine captures the public interface of a smart record virtual machine.
type Machine interface {
	Update(writer peer.ID, k string, update ir.Dict) error // Updates the dictionary in the writer's private space.
	Get(k string) RecordValue                              // Get the full Record in a key
	// NOTE: No query operation will be supported until we figure out selectors
	// Query(key string, selector Selector) (RecordValue, error)
}

// VM implements the Machine interface and keeps the map of records in its state.
type vm struct {
	ctx ir.UpdateContext // UpdateContext the VM uses to resolve conflicts
	//ds  ds.Datastore    // TODO: Add a datastore instead of using map[string] for the VM state
	keys map[string]*RecordValue // State of the VM storing the map of records.
	asm  ir.Assembler            // Assemble to use in the VM.
	lk   sync.RWMutex            // Lock to enable multiple access
}

// NewVM creates a new smart record Machine
func NewVM(ctx ir.UpdateContext, asm ir.Assembler) Machine {
	return newVM(ctx, asm)
}

//newVM instantiates a new VM with an updateContext and an assembler
func newVM(ctx ir.UpdateContext, asm ir.Assembler) *vm {
	return &vm{
		ctx:  ctx,
		keys: make(map[string]*RecordValue),
		asm:  asm,
	}
}

// Get the whole record stored in a key
func (v *vm) Get(k string) RecordValue {
	// TODO: Implementation
	return RecordValue{}
}

// Update the dictionary in the writer's private space
func (v *vm) Update(writer peer.ID, k string, update ir.Dict) error {
	// TODO: Implementation
	return nil
}
