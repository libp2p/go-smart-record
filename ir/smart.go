package ir

// Smart represents a "smart tag". A smart tag is syntactically represented by a Dict.
// Semantically, a smart tag can have custom merge logic that is different from that of a Dict.
type Smart interface {
	// Every smart tag is also a valid syntactic node.
	// This enables a smart tags to be used as keys or values in a Dict.
	Node

	// Dict returns the syntactic representation of the smart tag.
	// Every smart tag has a syntactic representation.
	// A syntactic representation includes Dict and literal nodes alone.
	// Equality of smart tags is definitionally equality of their syntactic representations.
	Disassemble() Dict
}

// Equality of smart tags is definitionally equality of their syntactic representations.
func IsEqualSmart(x, y Smart) bool {
	return IsEqualDict(x.Disassemble(), y.Disassemble())
}
