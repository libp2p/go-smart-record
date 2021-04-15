package ir

type UpdateContext interface{}

type DefaultUpdateContext struct{}

func Update(ctx UpdateContext, old, update Node) (Node, error) {
	return old.UpdateWith(ctx, update)
}
