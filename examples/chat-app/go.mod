module github.com/libp2p/go-smart-record/examples/chat-app

go 1.16

require (
	github.com/gdamore/tcell/v2 v2.2.1 // indirect
	github.com/libp2p/go-libp2p v0.13.1-0.20210420165741-6a5da01b0449
	github.com/libp2p/go-libp2p-core v0.8.6-0.20210415043615-525a0b130172
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-smart-record v0.0.0-00010101000000-000000000000
	github.com/multiformats/go-multiaddr v0.3.1
	github.com/rivo/tview v0.0.0-20210426144334-3ac88670ddeb
)

replace github.com/libp2p/go-smart-record => ../../
