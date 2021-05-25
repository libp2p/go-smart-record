package protocol

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
)

// Protocol ID
const (
	srProtocol protocol.ID = "/smart-record/0.0.1"
)

// Options is a structure containing all the options that can be used when constructing the smart records env
type serverConfig struct {
	//datastore          ds.Batching
	updateContext ir.UpdateContext
	assembler     ir.Assembler
	// NOTE: Add option for VM garbage collection if needed
	// gcPeriod      time.Duration
}

// Option type for smart records
type ServerOption func(*serverConfig) error

// defaults are the default smart record env options. This option will be automatically
// prepended to any options you pass to the constructor.
var serverDefaults = func(o *serverConfig) error {
	o.updateContext = ir.DefaultUpdateContext{}
	o.assembler = base.BaseGrammar

	return nil
}

// apply applies the given options to this Option
func (c *serverConfig) apply(opts ...ServerOption) error {
	for i, opt := range opts {
		if err := opt(c); err != nil {
			return fmt.Errorf("smart record server option %d failed: %s", i, err)
		}
	}
	return nil
}

// Assembler  configures the assembler to use in the smart record VM
func Assembler(asm ir.Assembler) ServerOption {
	return func(c *serverConfig) error {
		c.assembler = asm
		return nil
	}
}

// UpdateContext configures the context to use for updates in the smart record VM
func UpdateContext(uc ir.UpdateContext) ServerOption {
	return func(c *serverConfig) error {
		c.updateContext = uc
		return nil
	}
}

//NOTE: No options available for clients at this moment. Leaving this here as a placeholder.
// Options is a structure containing all the options that can be used when constructing the smart records env
type clientConfig struct {
}

// Option type for smart records
type ClientOption func(*clientConfig) error

// apply applies the given options to this Option
func (c *clientConfig) apply(opts ...ClientOption) error {
	for i, opt := range opts {
		if err := opt(c); err != nil {
			return fmt.Errorf("smart record client option %d failed: %s", i, err)
		}
	}
	return nil
}

var clientDefaults = func(o *clientConfig) error {
	return nil
}
