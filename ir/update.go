package ir

type UpdateContext interface{}

type DefaultUpdateContext struct{}

// Update updates the node in the first argument with
// the node in the second argument.
// NOTE: I don't think this top-level Update function
// is needed anymore. Consider removing it.
func Update(ctx UpdateContext, old, update Node) error {
	return old.UpdateWith(ctx, update)
}
