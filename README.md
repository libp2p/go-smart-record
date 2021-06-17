# go-smart-record

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23yellow)
<!-- TODO: Uncomment when available
[![GoDoc](https://godoc.org/github.com/libp2p/go-libp2p-xor?status.svg)](https://godoc.org/github.com/libp2p/go-smart-record)
[![Build Status](https://travis-ci.org/libp2p/go-libp2p-xor.svg?branch=master)](https://travis-ci.org/libp2p/go-smart-record)
-->
[![Discourse posts](https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg)](https://discuss.libp2p.io)

> Go implementation of smart records. Smart Records (SRs) provide
a *public blackboard for protocols*. 

## Summary
We currently don't have a standardized, shared, public medium decoupled from the transport protocol for the interaction of different protocols.
With Smart Records (SRs) we generalize DHT's key/value put/get as a separate protocol that can be leveraged by any other protocol
(including DHT protocols) to store arbitrary data. 

Smart Records leverage the [go-routing-language](https://github.com/libp2p/go-routing-language) for their data model. 
SRs are a mixture between a CRDT and a smart contract. A record (for a key) is a replicated state machine holding generic data.
It supports reading, writing, merging and "smart services" (through smart tags included in the SR data model which adds additional
logic to records).

SR model works as follows:
- Each peer (identified by its public key) writes to a peer-specific documents (their own peer private space).
- Peers can overwrite their own documents.
- Every document node has a TTL specified (and eventually paid for) by the writing peer.
- Users of SR can get the full record and process the information stored in the different user-spaces.

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/libp2p/go-libp2p-xor/issues).

Check out our [contributing document](https://github.com/libp2p/community/blob/master/CONTRIBUTE.md) for more information on how we work, and about contributing in general. Please be aware that all interactions related to libp2p are subject to the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

## License

[MIT](LICENSE) © Protocol Labs Inc.

