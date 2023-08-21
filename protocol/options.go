package protocol

import (
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
)

// Protocol ID
const (
	srid          protocol.ID = "/smart-record/0.0.1"
	DefaultPrefix protocol.ID = "/ipfs"
)

// Options is a structure containing all the options that can be used when constructing the smart records env
type serverConfig struct {
	//datastore          ds.Batching
	updateContext  ir.UpdateContext
	assembler      ir.AssemblerContext
	gcPeriod       time.Duration
	protocolPrefix protocol.ID
}

// Option type for smart records
type ServerOption func(*serverConfig) error

// defaults are the default smart record env options. This option will be automatically
// prepended to any options you pass to the constructor.
var serverDefaults = func(o *serverConfig) error {
	o.updateContext = ir.DefaultUpdateContext{}
	o.assembler = ir.AssemblerContext{Grammar: base.BaseGrammar}
	o.protocolPrefix = DefaultPrefix

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

// ClientProtocolPrefix configures the smart-record client protocol ID prefix.
func ClientProtocolPrefix(p protocol.ID) ClientOption {
	return func(c *clientConfig) error {
		c.protocolPrefix = p
		return nil
	}
}

// ServerProtocolPrefix configures the smart-record client protocol ID prefix.
func ServerProtocolPrefix(p protocol.ID) ServerOption {
	return func(c *serverConfig) error {
		c.protocolPrefix = p
		return nil
	}
}

// Assembler  configures the assembler to use in the smart record VM
func Assembler(asm ir.AssemblerContext) ServerOption {
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

// VMGcPeriod configures the garbage collection granularity in the server VM.
func VMGcPeriod(gcP time.Duration) ServerOption {
	return func(c *serverConfig) error {
		c.gcPeriod = gcP
		return nil
	}
}

// Options is a structure containing all the options that can be used when constructing the smart records env
type clientConfig struct {
	protocolPrefix protocol.ID
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
	o.protocolPrefix = DefaultPrefix
	return nil
}
