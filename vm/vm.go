// Package vm defines and implements a Virtual Machine for smart records.
package vm

import (
	"fmt"

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
	ctx ir.MergeContext     // MergeContext the VM uses to resolve conflicts
	s   map[string]*ir.Dict // State of the VM storing the map of records.
}

func NewVM(ctx ir.MergeContext) *VM {
	return &VM{
		ctx: ctx,
		s:   make(map[string]*ir.Dict),
	}
}

// Update updates the record in key with the provided dict
func (v *VM) Update(k string, d ir.Dict) error {

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
func (v *VM) Query(k string, selector ir.Dict) (ir.Dict, error) {
	src := v.s[k]
	if src == nil {
		return ir.Dict{}, fmt.Errorf("empty key in state")
	}

	// Traverse the selector and check if it exists in the stored dict for the key
	// If the path exists return it, if not do nothing
	return queryDict(*src, selector)

}

func queryDict(src ir.Dict, selector ir.Dict) (ir.Dict, error) {
	// Check if selector or src equals nil
	out := ir.Dict{}
	if src.Tag == selector.Tag {
		out.Tag = selector.Tag
		for _, p := range selector.Pairs {
			// Check if selector and source are the same type
			// and is not a wildcard (i.e. value of selector == nil).
			srcP := src.Get(p.Key)
			if p.Value != nil && !ir.IsEqualType(p.Value, srcP) {
				continue
			}

			switch p.Value.(type) {
			case ir.Dict:
				srcDict := src.Get(p.Key).(ir.Dict)
				selectorDict := selector.Get(p.Key).(ir.Dict)

				// If no pairs specified it means that the full Dict wants to be returned.
				// For now the wildcard is an empty Node.
				if len(selectorDict.Pairs) == 0 {
					// Check if the Tag is equal
					if selectorDict.Tag == srcDict.Tag {
						out = out.CopySet(p.Key, srcDict)
					}
				} else {
					// If not query the specfied Pairs in the dict.
					tmpQuery, err := queryDict(src.Get(p.Key).(ir.Dict), selector.Get(p.Key).(ir.Dict))
					if err != nil {
						return ir.Dict{}, err
					}
					out = out.CopySet(p.Key, tmpQuery)

				}
			default:
				value := src.Get(p.Key)
				out = out.CopySet(p.Key, value)
			}
		}
	}
	return out, nil
}

func checkTypes(x, y ir.Node) {
	// Check if both are the same type. This is somethign to do
	// above to avoid panic when comparing different values.
}

// Gets the whole record stored in a key
func (v *VM) Get(k string) ir.Dict {
	return *v.s[k]
}
