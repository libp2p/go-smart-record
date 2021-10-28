module github.com/libp2p/go-smart-record/examples/chat-app

go 1.16

require (
	github.com/gdamore/tcell/v2 v2.2.1
	github.com/libp2p/go-libp2p v0.15.1
	github.com/libp2p/go-libp2p-core v0.9.0
	github.com/libp2p/go-routing-language v0.0.0-20210531170722-12dc033e88ac
	github.com/libp2p/go-smart-record v0.0.0-00010101000000-000000000000
	github.com/multiformats/go-multiaddr v0.4.0
	github.com/rivo/tview v0.0.0-20210426144334-3ac88670ddeb
)

replace github.com/libp2p/go-smart-record => ../../
