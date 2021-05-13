package ir

type UpdateContext interface{}

type DefaultUpdateContext struct{}

func Update(ctx UpdateContext, old, update Node) error {
	return old.UpdateWith(ctx, update)
}
