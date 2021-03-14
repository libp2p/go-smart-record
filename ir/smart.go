package ir

// Smart represents a "smart tag". A smart tag is syntactically a Dict.
// Semantically, a smart tag can have custom merge logic
// that is different from that of a Dict.
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
	return IsEqualDict(x.Dict(), y.Dict())
}
