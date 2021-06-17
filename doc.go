/*
Package go-smart-record includes the go implementation of Smart Records. Smart Records (SRs) provide
a public blackboard for protocols.

Overview

We currently don't have a standardized, shared, public
medium for:
	- writing and reading
	- by multiple participants
	- talking multiple protocols
	- scattered in multiple locations.

Traditionally, DHTs have been used for this purpose. With Smart Record we generalize DHT's
key/value put/get as a separate protocol that can be leveraged by any other protocol
(including DHT protocols).
	- UPDATE/GET interface of SRs are used to interact with records stored in a peer,
	delegating the FIND operations of records to the DHT (or other available protocols),
	decoupling the storage of records from transport protocols.
	- Records become portable *data state machines*. They can be sent and updated using pubsub,
	aggregated with other versions of the records even if a peer doesn't fully understand it, etc.


SRs are a mixture between a CRDT and a smart contract. A record (for a key) is a replicated state machine holding generic data.
It supports reading, writing, merging and "smart services" (through smart tags included in the SR data model).
	- The layman description of smart records: *"they are DHT values which become publicly updatable JSON/IPLD documents by any peer*.

Model

SRs work as follows:
	- Each peer (identified by its public key) writes to a peer-specific documents.
	- Peers can overwrite their own documents.
	- Every document node has a TTL specified (and eventually paid for) by the writing peer.
	- Users of SR can get the full record and process the information stored in the different user-spaces.

Architecture

The SR system has the following architecture:
	- Syntactic representation (xr): Data model used by protocol and application developers
	to interact with smart records. In their current implementation smart-records can
	be transformed into the IPLD data model, and serialized/desearialized seamlessly for
	transmission or any other purpose.

	- Semantic representation (ir): Intemediate representation used by the SR VM.
	Syntactic nodes are assembled into semantic nodes. In the assembly process, tags are parsed
	and certain nodes may be transformed into smart nodes and trigger additional (i.e. "smart")
	operations in the VM.

	- VM (vm): The VM is responsible for storing and updating the SR stored in a peer. It exposes the
	SR interface to the "outside world" and triggers smart-tag operations when appropiate.
	The "outside world" use syntctic nodes to intercat with the VM interface that the VM assembles
	and stores in its datastore in its semantic form.

	- Libp2p SR request/response protocol (protocol) to interact with other peers SR.
	It includes a server implementation that instantiates a SR VM and exposes the SR interface
	through the network to other peers, and a client implementation that can be leveraged by
	applications and protocols to make requests to SR servers.


Use cases

Some examples of things you can do with SR:
	- Deploy new applications without upgrading the whole network.
	- Design protocols that can interact with other protocols.
	- Facilitate cryptographic protocols that require a "trusted" party
		- Fair exchange
	- Unlock application development on the DHT to the public
		- Private group chat, custom routing, decentralized limit-order marketbook, etc.
		- New app types: Interaction between trustless parties, using a public jury.
*/
package main
