package ir

import (
	"fmt"
)

type MergeContext interface {
	// MergeConflict merges two nodes that cannot be merged using the default merge semantics
	// of dictionaries and primitive value types.
	// MergeConflict should throw a panic, when it is unable to merge.
	MergeConflict(Node, Node) (Node, error)
}

type DefaultMergeContext struct{}

func (DefaultMergeContext) MergeConflict(Node, Node) (Node, error) {
	return nil, fmt.Errorf("cannot resolve a merge conflict")
}

func Merge(ctx MergeContext, x, y Node) (Node, error) {
	if xs, ok := x.(Smart); ok {
		return xs.MergeWith(ctx, y)
	}
	if ys, ok := y.(Smart); ok {
		return ys.MergeWith(ctx, x)
	}
	switch x1 := x.(type) {
	case Bool, String, Number, Blob: // literals merge without conflict if they are equal
		if IsEqual(x, y) {
			return x, nil
		}
	case Dict:
		switch y1 := y.(type) {
		case Dict:
			return MergeDict(ctx, x1, y1)
		}
	case Set:
		switch y1 := y.(type) {
		case Set:
			return MergeSet(ctx, x1, y1)
		}
	}
	return ctx.MergeConflict(x, y) // defer unresolvable conflicts to the context
}
