// Package vm defines and implements a Virtual Machine for smart records.
package vm

import (
	"fmt"
	"sync"

	peer "github.com/libp2p/go-libp2p-peer"
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
	v.lk.RLock()
	defer v.lk.RUnlock()
	if v.keys[k] == nil {
		return Record{Key: k, Value: RecordValue{}}
	}
	return *v.keys[k]
}

// Update the dictionary in the writer's private space
// NOTE: We currently store an assembled version of the record.
// We may need to disassemble and serialize before storage
// if we choose to use a datastore.
func (v *vm) Update(writer peer.ID, k string, update ir.Dict) error {
	v.lk.Lock()
	defer v.lk.Unlock()

	// Start assemble process with the parent VM assemblerContext
	ds, err := v.asm.Assemble(ir.AssemblerContext{Grammar: v.asm}, update)
	if err != nil {
		return err
	}

	// Check if the result of the assembler is of type Dict
	d, ok := ds.(ir.Dict)
	if !ok {
		return fmt.Errorf("assembler didn't generate a dict")
	}

	// Directly store d if there is nothing in the key
	if v.keys[k] == nil {
		v.keys[k] = &Record{Key: k, Value: make(map[peer.ID]*ir.Dict, 0)}
		v.keys[k].Value[writer] = &d
		return nil
	} else {
		// If no data in peer
		if v.keys[k].Value[writer] == nil {
			v.keys[k].Value[writer] = &d
		} else {
			// Update existing dict with the stored one if there's already
			// something in the peer's key
			n, err := ir.Update(v.ctx, *v.keys[k].Value[writer], d)
			if err != nil {
				return nil
			}
			*v.keys[k].Value[writer] = n.(ir.Dict)
		}

	}
	return nil
}
