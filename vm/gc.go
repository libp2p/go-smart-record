package vm

import (
	"time"

	"github.com/jbenet/goprocess"
	"github.com/libp2p/go-smart-record/ir"
)

func (v *vm) gcLoop(proc goprocess.Process) {
	for {
		msgSyncTicker := time.NewTicker(v.gcPeriod)
		select {
		case <-msgSyncTicker.C:
			// Stopping ticker while garbage collecting.
			msgSyncTicker.Stop()
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
			// Run garbage collection
			if gcDict(entry) {
				// Delete that entry if dict for peer expired.
				delete(*r, p)
			}
		}
	}
}

func gcNode(n ir.Node) bool {
	switch n1 := n.(type) {
	case *ir.Dict:
		return gcDict(n1)
	case *ir.Set:
		return gcSet(n1)
	default:
		return isTTLExpired(n1)
	}
}

func gcDict(d *ir.Dict) bool {
	// Check if we can remove Dict if all children have expired.
	gcFlag := isTTLExpired(d)
	// For each pair.
	for k := len(d.Pairs) - 1; k >= 0; k-- {
		// Check if pair has expired and garbage collect.
		gcP := gcNode(d.Pairs[k].Key) && gcNode(d.Pairs[k].Value)
		if gcP {
			// Remove pair if both expired
			d.Remove(d.Pairs[k].Key)
		}
		// Accummulate the result for the child in dict.
		gcFlag = gcFlag && gcP
	}
	return gcFlag
}

func gcSet(s *ir.Set) bool {
	// Check if we can remove Dict if all children have expired.
	gcFlag := isTTLExpired(s)
	// For each element
	for k := len(s.Elements) - 1; k >= 0; k-- {
		// Check if element gas expired
		gcP := gcNode(s.Elements[k])
		if gcP {
			// Remove pair if both expired
			s.Elements = append(s.Elements[:k], s.Elements[k+1:]...)
		}
		// Accummulate the result for the child in set
		gcFlag = gcFlag && gcP
	}
	return gcFlag
}

func isTTLExpired(n ir.Node) bool {
	if uint64(time.Now().Unix()) > n.Metadata().ExpirationTime {
		return true
	}
	return false
}
