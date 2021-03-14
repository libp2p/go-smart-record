// Package ir defines the Intermediate Representation (informally, in-memory representation) of smart records.
package ir

/*

The IR is a "vocabulary" of nodes which can be used to construct "documents".
The vocabulary consists of two sets of nodes: syntactic and smart.

	Syntactic nodes:
		Dict
		String
		Number
		Int64	(TODO: Are fixed precision literals appropriate?)
		Blob

		(TODO: What is the right set of literal types?)

	Smart nodes:
		Cid
		Multiaddress
		Peer
		Record
		Sign
		Verify, Verified

SERIALIZATION

	JSON --(XXX)--> Syntactic IR --(interpretation)--> Syntactic+Smart IR

	JSON <--(XXX)-- Syntactic IR <--(generation)-- Syntactic+Smart IR

*/
