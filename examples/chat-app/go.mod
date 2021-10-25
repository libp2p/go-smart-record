module github.com/libp2p/go-smart-record/examples/chat-app

go 1.16

require (
	github.com/libp2p/go-libp2p v0.15.1
	github.com/libp2p/go-libp2p-core v0.11.0
	github.com/libp2p/go-routing-language v0.0.0-20210531170722-12dc033e88ac
	github.com/libp2p/go-smart-record v0.0.0-00010101000000-000000000000
	github.com/multiformats/go-multiaddr v0.4.0
)

replace github.com/libp2p/go-smart-record => ../../
