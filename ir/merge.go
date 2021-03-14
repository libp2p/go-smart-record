package ir

type MergeContext interface {
	MergeConflict(Node, Node) Node
}

func Merge(ctx MergeContext, x, y Node) Node {
	if xs, ok := x.(Smart); ok {
		return xs.MergeWith(ctx, y)
	}
	if ys, ok := y.(Smart); ok {
		return ys.MergeWith(ctx, x)
	}
	switch x1 := x.(type) {
	case String, Int64, Number, Blob: // literals merge without conflict if they are equal
		if IsEqual(x, y) {
			return x
		}
	case Dict:
		switch y1 := y.(type) {
		case Dict:
			return MergeDict(ctx, x1, y1)
		}
	}
	return ctx.MergeConflict(x, y) // defer unresolvable conflicts to the context
}
