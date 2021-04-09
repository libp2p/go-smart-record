// Package vm defines and implements a Virtual Machine for smart records.
package vm

import (
	"fmt"

	//ds "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-smart-record/ir"
)

// Machine captures the public interface of a smart record virtual machine.
// NOTE: Keys must be the same type as ir.Record Key (at least for now)
type Machine interface {
	Update(k string, d ir.Dict) error                  // Updates the record in key k with a dict
	Query(k string, selector ir.Dict) (ir.Dict, error) // Queries a record using a selector dict.
	Get(k string) ir.Dict                              // Gets the whole record stored in a key (debug purposes for now)

}

// VM implements the Machine interface and keeps the map of records in its state.
type VM struct {
	ctx ir.MergeContext // MergeContext the VM uses to resolve conflicts
	//ds  ds.Datastore    // TODO: Add a datastore instead of using map[string] for the VM state
	s   map[string]*ir.Dict // State of the VM storing the map of records.
	asm ir.Assembler        // Assemble to use in the VM.
}

func NewVM(ctx ir.MergeContext, asm ir.Assembler) *VM {
	return &VM{
		ctx: ctx,
		s:   make(map[string]*ir.Dict),
		asm: asm,
	}
}

// Update updates the record in key with the provided dict
func (v *VM) Update(k string, s ir.Dict) error {
	// Start assemble process with an empty root context
	ds, err := v.asm.Assemble(ir.AssemblerContext{}, s)
	if err != nil {
		return err
	}
	// Check if the result of the assembler is a dict ready to store.
	d, ok := ds.(ir.Dict)
	if !ok {
		return fmt.Errorf("assembler didn't generate a dict")
	}

	// Directly store d if there is nothing in the key
	if v.s[k] == nil {
		v.s[k] = &d
		return nil
	} else {
		// Merge existing dict with the stored one if there's already
		// something in the key
		n, err := ir.MergeDict(v.ctx, *v.s[k], d)
		if err != nil {
			return nil
		}
		*v.s[k] = n.(ir.Dict)
	}
	return nil
}

// Query receives a dict selector as input and traverses the dict in the key
// to return the corresponding values
func (v *VM) Query(k string, selector Selector) (ir.Dict, error) {
	src := v.s[k]
	if src == nil {
		return ir.Dict{}, fmt.Errorf("empty key in state")
	}

	// Traverse the selector and check if it exists in the stored dict for the key
	// If the path exists return it, if not do nothing
	return selector.Run(SelectorContext{}, src)

}

// Gets the whole record stored in a key
func (v *VM) Get(k string) ir.Dict {
	// TODO: Disassemble before returning the dict
	return *v.s[k]
}
