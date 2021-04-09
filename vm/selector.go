package vm

import "github.com/libp2p/go-smart-record/ir"

type SelectorContext struct {
	s Selector
}

type Selector interface {
	Run(ctx SelectorContext, src ir.Dict) (ir.Dict, error)
}

type SyntacticDictSelector struct {
	d ir.Dict
}

func (s SyntacticDictSelector) Run(ctx SelectorContext, src ir.Dict) (ir.Dict, error) {
	// Traverse the selector and check if it exists in the stored dict for the key
	// If the path exists return it, if not do nothing
	return queryDict(src, s.d)
}

func queryDict(src ir.Dict, selector ir.Dict) (ir.Dict, error) {
	// Check if selector or src equals nil
	out := ir.Dict{}
	//if src.Tag == selector.Tag {
	//out.Tag = selector.Tag
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
	//}
	return out, nil
}
