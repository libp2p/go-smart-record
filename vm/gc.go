package vm

import (
	"fmt"
	"time"

	"github.com/jbenet/goprocess"
	"github.com/libp2p/go-smart-record/ir"
)

// gcPeriod determines the granularity of GC by the VM.
// NOTE: Make it a configurable parameter at initialization.
const gcPeriod = 2 * time.Second

func (v *vm) gcLoop(proc goprocess.Process) {
	// TODO: Add gcInterval as an option
	msgSyncTicker := time.NewTicker(gcPeriod)
	defer msgSyncTicker.Stop()
	for {
		select {
		case <-msgSyncTicker.C:
			fmt.Println("Garbage collect triggered")
			// NOTE: Locking all keys while garbage collecting may really
			// harm performance, specially if the gcPeriod is low. Consider
			// adding a lock per entry or other schemes to improve this.
			// It'd be useful to gather some metrics.
			v.lk.Lock()
			v.garbageCollect()
			v.lk.Unlock()

		case <-proc.Closing():
			return
		}
	}
}

func (v *vm) garbageCollect() {
	// For each record
	for _, r := range v.keys {
		// And the datastore of each peer
		for p, entry := range *r {
			if gc := gcDict(entry); gc {
				// Delete that entry if dict for peer expired.
				delete(*r, p)
			}
		}
	}
}

func gcNode(n ir.Node) bool {
	switch n1 := n.(type) {
	case ir.Dict:
		return gcDict(&n1)
	case ir.Set:
		return gcSet(&n1)
	default:
		return isTTLExpired(n1)
	}
}

func gcDict(d *ir.Dict) bool {
	// Check if we can remove Dict if all children have expired.
	gcFlag := isTTLExpired(d)
	pairs := make(ir.Pairs, 0)
	// For each pair.
	for _, p := range d.Pairs {
		// Check if pair has expired and garbage collect.
		gcP := gcNode(p.Key) && gcNode(p.Value)
		// If it hasn't keep the pair
		if !gcP {
			pairs = append(pairs, p)
		}
		// Accummulate the result for the child in dict.
		gcFlag = gcFlag && gcP
	}
	// Assign the pairs that haven't expired to dict.
	if !gcFlag {
		d.Pairs = pairs
	}
	// Return gc result for dict.
	return gcFlag
}

func gcSet(s *ir.Set) bool {
	// Check if we can remove Dict if all children have expired.
	gcFlag := isTTLExpired(s)
	els := make(ir.Nodes, 0)
	// For each elements.
	for _, e := range s.Elements {
		// Check if element gas expired
		gcP := gcNode(e)
		// If it hasn't keep the element
		if !gcP {
			els = append(els, e)
		}
		// Accummulate the result for the child in set
		gcFlag = gcFlag && gcP
	}
	// Assign the pairs that haven't expired to set.
	if !gcFlag {
		s.Elements = els
	}
	// Return gc result for dict.
	return gcFlag
}

func isTTLExpired(n ir.Node) bool {
	if uint64(time.Now().Unix()) > n.Metadata().TTL+n.Metadata().AssemblyTime {
		return true
	}
	return false
}
