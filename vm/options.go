package vm

import (
	"fmt"
	"time"
)

// Protocol ID
const (
	// gcPeriod determines the granularity of GC by the VM.
	gcPeriod = 60 * time.Second
)

// Options is a structure containing all the options for VM
type vmConfig struct {
	gcPeriod time.Duration
}

// Option type
type VMOption func(*vmConfig) error

// defaults are the default vm option. This option will be automatically
// prepended to any options you pass to the constructor.
var defaults = func(o *vmConfig) error {
	o.gcPeriod = gcPeriod
	return nil
}

// apply applies the given options to this Option
func (c *vmConfig) apply(opts ...VMOption) error {
	for i, opt := range opts {
		if err := opt(c); err != nil {
			return fmt.Errorf("smart record server option %d failed: %s", i, err)
		}
	}
	return nil
}

// GCPeriod garbage collection period for VM
func GCPeriod(p time.Duration) VMOption {
	return func(c *vmConfig) error {
		c.gcPeriod = p
		return nil
	}
}
