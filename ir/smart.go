package ir

// Smart is a "smart" tag.
// A smart tag is syntactically equivalent to a Dict.
// Semantically, a smart tag can have custom equality and merge logics
// that are different from those of a Dict.
type Smart interface {
	// Every smart tag is also a valid syntactic node.
	Node

	// Dict returns the syntactic representation of the smart tag.
	// A syntactic representation is built from Dict and literal nodes alone.
	Dict() Dict

	// MergeWith returns the result of merging this smart tag with the given node.
	// If x is smart as well, it should be the case that x.MergeWith(y) = y.MergeWith(x).
	MergeWith(ctx MergeContext, x Node) Node
}

func IsEqualSmart(x, y Smart) bool {
	panic("XXX")
}
