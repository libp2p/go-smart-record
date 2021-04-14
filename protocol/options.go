package protocol

import (
	"fmt"

	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
)

// Options is a structure containing all the options that can be used when constructing the smart records env
type config struct {
	//datastore          ds.Batching
	mergeContext ir.MergeContext
	assembler    ir.Assembler
}

// Option type for smart records
type Option func(*config) error

// defaults are the default smart record env options. This option will be automatically
// prepended to any options you pass to the DHT constructor.
var defaults = func(o *config) error {
	o.mergeContext = ir.DefaultMergeContext{}
	o.assembler = base.BaseGrammar
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
