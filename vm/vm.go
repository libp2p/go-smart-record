package vm

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	xr "github.com/libp2p/go-routing-language/syntax"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
	meta "github.com/libp2p/go-smart-record/ir/metadata"
)

// RecordEntry determines the structure of data stored in a record.
// Each peer has a private dataspace to store semantic dicts.
type recordEntry map[peer.ID]*ir.Dict

// RecordValue determines the structure used by the VM to interact with
// the outside world, outputing disassembled syntactic dicts.
type RecordValue map[peer.ID]*xr.Dict

// Machine captures the public interface of a smart record virtual machine.
type Machine interface {
	Update(writer peer.ID, k string, update xr.Dict, metadata ...meta.Metadata) error // Updates the dictionary in the writer's private space.
	Get(k string) RecordValue                                                         // Get the full Record in a key
	// NOTE: No query operation will be supported until we figure out selectors
	// Query(key string, selector Selector) (RecordValue, error)
	Close() error
}

// VM implements the Machine interface and keeps the map of records in its state.
type vm struct {
	ctx  context.Context
	proc goprocess.Process
	lk   sync.RWMutex // Lock to enable multiple access
	host host.Host

	updateCtx ir.UpdateContext // UpdateContext the VM uses to resolve conflicts
	// ds  ds.Datastore    // TODO: Add a datastore instead of using map[string] for the VM state
	keys map[string]*recordEntry // State of the VM storing the map of records.
	asm  ir.AssemblerContext     // Assemble to use in the VM.

	// NOTE: When performance matters in the future, implement incremental garbage collection,
	// which runs on every operation and uses a priority queue to know (in O(1) time)
	// if anything needs garbage collection.
	// (When there are bursts of uneven traffic, no choice of garbage collection interval helps.)
	// We can add it in a GCType option.
	gcPeriod time.Duration // Period of the gc process
}

// NewVM creates a new smart record Machine
func NewVM(ctx context.Context, h host.Host, updateCtx ir.UpdateContext, asm ir.AssemblerContext, options ...VMOption) (Machine, error) {
	return newVM(ctx, h, updateCtx, asm, options...)
}

// newVM instantiates a new VM with an updateContext and an assembler
func newVM(ctx context.Context, h host.Host, updateCtx ir.UpdateContext, asm ir.AssemblerContext, options ...VMOption) (*vm, error) {
	var cfg vmConfig
	if err := cfg.apply(append([]VMOption{defaults}, options...)...); err != nil {
		return nil, err
	}
	v := &vm{
		ctx:       ctx,
		host:      h,
		updateCtx: updateCtx,
		keys:      make(map[string]*recordEntry),
		asm:       asm,
		gcPeriod:  cfg.gcPeriod,
	}

	// Initialize process so routines are ended with context
	v.proc = goprocessctx.WithContext(ctx)
	// Start garbage collection process
	// NOTE: Add an option for gcType?
	v.proc.Go(v.gcLoop)
	return v, nil
}

// Get the whole record stored in a key
func (v *vm) Get(k string) RecordValue {
	v.lk.RLock()
	defer v.lk.RUnlock()
	// If nothing in key
	if v.keys[k] == nil {
		return RecordValue{}
	}

	// Disassembles all nodes in record
	out := make(map[peer.ID]*xr.Dict)
	for pk, v := range *v.keys[k] {
		d := v.Disassemble()
		do, ok := d.(xr.Dict)
		// Do not return nodes which are not ir.Dict
		// after dissassembling
		if ok {
			out[pk] = &do
		}
	}
	return out
}

// Update the dictionary in the writer's private space
// NOTE: We currently store an assembled version of the record.
// We may need to disassemble and serialize before storage
// if we choose to use a datastore for persistance.
func (v *vm) Update(writer peer.ID, k string, update xr.Dict, metadata ...meta.Metadata) error {
	v.lk.Lock()
	defer v.lk.Unlock()

	// Start assemble process with the parent VM assemblerContext
	ds, err := v.asm.Grammar.Assemble(v.asm, update, metadata...)
	if err != nil {
		return err
	}

	// Check if the result of the assembler is of type Dict
	d, ok := ds.(*ir.Dict)
	if !ok {
		return fmt.Errorf("assembler didn't generate a dict")
	}

	// Trigger reachibility verifications.
	// NOTE: What about considering once we have more than one smart tag
	// a general method to define the lifecycle of smart tags
	// and when to trigger them?
	base.TriggerReachable(d, v.host)

	// Directly store d if there is nothing in the key
	if v.keys[k] == nil {
		v.keys[k] = &recordEntry{writer: d}
		return nil
	} else {
		// If no data in peer
		if (*v.keys[k])[writer] == nil {
			(*v.keys[k])[writer] = d
		} else {
			// Update existing dict with the stored one if there's already
			// something in the peer's key
			err := ir.Update(v.ctx, (*v.keys[k])[writer], d)
			if err != nil {
				return nil
			}
		}
		return nil
	}
}

// Close calls Process Close.
func (v *vm) Close() error {
	return v.proc.Close()
}
