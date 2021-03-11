package ir

// Smart represents a "smart" node.
type Smart interface {
	Node
	Dict() Dict
}

func IsEqualSmart(x, y Smart) bool {
	panic("XXX")
}
