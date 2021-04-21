package protocol

import (
	"fmt"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
)

// ModeOpt publicly exposes the modes in which the smartRecordManager can operate
type ModeOpt int

const (
	ModeClient ModeOpt = iota
	ModeServer
)

// Options is a structure containing all the options that can be used when constructing the smart records env
type config struct {
	//datastore          ds.Batching
	updateContext ir.UpdateContext
	assembler     ir.Assembler
	mode          ModeOpt
}

// Option type for smart records
type Option func(*config) error

// defaults are the default smart record env options. This option will be automatically
// prepended to any options you pass to the constructor.
var defaults = func(o *config) error {
	o.updateContext = ir.DefaultUpdateContext{}
	o.assembler = base.BaseGrammar
	o.mode = ModeServer

	return nil
}

// apply applies the given options to this Option
func (c *config) apply(opts ...Option) error {
	for i, opt := range opts {
		if err := opt(c); err != nil {
			return fmt.Errorf("smart record option %d failed: %s", i, err)
		}
	}
	return nil
}

// Mode configures which mode the smartRecordManager operates in (Client, Server).
//
// Defaults to Server.
func Mode(m ModeOpt) Option {
	return func(c *config) error {
		c.mode = m
		return nil
	}
}

// Assembler  configures the assembler to use in the smart record VM
func Assembler(asm ir.Assembler) Option {
	return func(c *config) error {
		c.assembler = asm
		return nil
	}
}

// UpdateContext configures the context to use for updates in the smart record VM
func UpdateContext(uc ir.UpdateContext) Option {
	return func(c *config) error {
		c.updateContext = uc
		return nil
	}
}
