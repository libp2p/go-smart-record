// Package ir defines the Intermediate Representation (informally, in-memory representation) of smart records.
package ir

/*

The IR is a "vocabulary" of nodes which can be used to construct "documents".
The vocabulary consists of two types of nodes: syntactic and smart.

	Syntactic nodes, by their Go type:
		Dict
		Bool
		String
		Int
		Float
		Blob

	Smart nodes, by their Go type:
		Cid
		Multiaddress
		Peer
		Record
		Sign, Signed
		Verify, Verified

We use the following nomenclature:

	"Syntactic IR", or "syntactic documents", refers to documents comprising only syntactic nodes.
	"Semantic IR", or "semantic documents", refers to documents comprising syntactic and smart nodes.

Users generally manipulate semantic documents (or just "documents", for short),
consisting of both syntactic and smart nodes.
Syntactic nodes represent generic structured data types (and support generic merge logic).
Smart nodes represent higher concepts (with custom merge logics) that have a syntactic representation.

SERIALIZATION

Documents must be serialized when they are displayed to the user or
sent as an argument inside a network function call.

Serialization comprises two steps: disassembly and marshalling.

Disassembly converts a semantic document (with smart and syntactic nodes) into
a purely syntactic document. Simply, smart tags are substituted for their
syntactic representation.

Marshalling converts a syntactic document into a serialized form,
according to some serialization "encoding", which determines
the actual format of the serialized object.

In principle, syntactic documents can be serialized into any number
of on-the-wire formats (e.g. JSON, BSON, Protocol Buffers, Flat Buffers, etc.).
We use JSON as a running example. Serialization to JSON is also provided
by this library out-of-the-box.

Serialization:
	JSON <--(marshal)-- Syntactic IR <--(disassemble)-- Semantic IR

Deserialization:
	JSON --(unmarshal)--> Syntactic IR --(assemble)--> Semantic IR


*/
