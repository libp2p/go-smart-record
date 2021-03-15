// Package ir defines the Intermediate Representation (informally, in-memory representation) of smart records.
package ir

/*

The IR is a "vocabulary" of nodes which can be used to construct "documents".
The vocabulary consists of two sets of nodes: syntactic and smart.

	Syntactic nodes:
		Dict
		String
		Int
		Float
		Blob

	Smart nodes:
		Cid
		Multiaddress
		Peer
		Record
		Sign, Signed
		Verify, Verified

We use the following nomenclature:

	Syntactic IR refers to documents comprising only syntactic nodes.
	Semantic IR refers to documents comprising syntactic and smart nodes.

SERIALIZATION

	Documents are serialization-agnostic: They can be de/serialized to any number
	of standard formats (e.g. JSON, BSON, Protocol Buffers, Flat Buffers, etc.).
	We use JSON as a running example. Serialization to JSON is also provided
	by this library out-of-the-box, due to its applicability to HTTP REST interfaces.



	JSON --(unmarshal)--> Syntactic IR --(assemble)--> Semantic IR

	JSON <--(marshal)-- Syntactic IR <--(disassemble)-- Semantic IR

*/
