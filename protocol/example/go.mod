module github.com/libp2p/go-smart-record/protocol/example

go 1.16

require (
	github.com/libp2p/go-libp2p v0.13.0
	github.com/libp2p/go-libp2p-core v0.8.5
	github.com/libp2p/go-libp2p-kad-dht v0.11.1
	github.com/libp2p/go-libp2p-mplex v0.4.1
	github.com/libp2p/go-libp2p-secio v0.2.3
	github.com/libp2p/go-libp2p-yamux v0.5.1
	github.com/libp2p/go-smart-record v0.0.0-00010101000000-000000000000
	github.com/libp2p/go-tcp-transport v0.2.1
	github.com/multiformats/go-multiaddr v0.3.1
)

replace github.com/libp2p/go-smart-record => ../../
