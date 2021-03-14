// Package vm defines and implements a Virtual Machine for smart records.
package vm

import (
	"github.com/libp2p/go-smart-record/ir"
)

// Machine captures the public interface of a virtual machine.
type Machine interface {
	Merge(ir.Dict) error
	Get() ir.Dict
}
