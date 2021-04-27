/*
	This tutorial is inspired on this two go-libp2p example:
		- https://github.com/libp2p/go-libp2p-examples/blob/master/chat/chat.go
		- https://github.com/libp2p/go-libp2p-examples/tree/master/pubsub/chat
	This tutorial builds a chat using smart-records instead of direct communication or PubSub to broadcast messages.
*/
package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	crypto "github.com/libp2p/go-libp2p-crypto"
	"github.com/libp2p/go-smart-record/protocol"

	"github.com/multiformats/go-multiaddr"
)

// Timeout for requests to smart-record server
const reqTimeout = 10 * time.Second

// Sync new messages every second
const syncTime = 1 * time.Second

func main() {
	ctx := context.Background()

	// Available flags
	sourcePort := flag.Int("sp", 0, "Source port number")
	dest := flag.String("d", "", "Destination multiaddr string")
	nick := flag.String("nick", "", "optional nickname")
	room := flag.String("room", "", "Select the chat room you want to connect to")
	server := flag.Bool("server", false, "Display help")
	help := flag.Bool("help", false, "Initialized a new smart-record server")
	debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")

	flag.Parse()

	if *help {
		fmt.Printf("This program demonstrates a simple p2p chat application using libp2p and smart-records\n\n")
		fmt.Println("Usage: ")
		fmt.Println("  - Run a chat-server './chat-app -server -sp <SOURCE_PORT>' where <SOURCE_PORT> can be any port number.")
		fmt.Println("  - Start new chat clients using './chat-app -sp <SOURCE_PORT>' -d <SERVER_MULTIADDR> -room <ROOM_ID> -nick <NICKNAME>")
		fmt.Println("<SERVER_MULTIADDR> is the multiaddress to connect to the chat server and <ROOM_ID> the id of the chat room.")
		fmt.Println("Use -nick to select an optional <NICKNAME> ")

		os.Exit(0)
	}

	// If debug is enabled, use a constant random source to generate the peer ID. Only useful for debugging,
	// off by default. Otherwise, it uses rand.Reader.
	var r io.Reader
	if *debug {
		// Use the port number as the randomness source.
		// This will always generate the same host ID on multiple executions, if the same port number is used.
		// Never do this in production code.
		r = mrand.New(mrand.NewSource(int64(*sourcePort)))
	} else {
		r = rand.Reader
	}

	// Creates a new RSA key pair for this host.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, r)
	if err != nil {
		panic(err)
	}

	// 0.0.0.0 will listen on any interface device.
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))

	// libp2p.New constructs a new libp2p Host.
	// Other options can be added here.
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(priv),
	)

	if err != nil {
		panic(err)
	}

	// If server, initialize smartRecord server and hang forever
	if *server {
		_, err = protocol.NewSmartRecordServer(ctx, host)
		if err != nil {
			panic("Couldn't initialize smartRecord server")
		}

		// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
		var port string
		for _, la := range host.Network().ListenAddresses() {
			if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
				port = p
				break
			}
		}

		if port == "" {
			panic("was not able to find actual local port")
		}

		// Send message to connect to server.
		fmt.Printf("Run './chat-app -d /ip4/127.0.0.1/tcp/%v/p2p/%s -room roomName -nick nickname' on another console to start a chat client.\n", port, host.ID().Pretty())
		fmt.Printf("\n[*] Smart record chat server running...\n\n")
		// Hang forever
		select {}
	} else {
		// Initialize smartRecord client for chat client
		smClient, _ := protocol.NewSmartRecordClient(ctx, host)
		if err != nil {
			panic("Couldn't initialize smartRecord client")
		}

		// Turn the destination into a multiaddr.
		maddr, err := multiaddr.NewMultiaddr(*dest)
		if err != nil {
			log.Fatalln(err)
		}

		// Extract the peer ID from the multiaddr.
		info, err := peer.AddrInfoFromP2pAddr(maddr)
		if err != nil {
			log.Fatalln(err)
		}

		// Add the destination's peer multiaddress in the peerstore.
		// This will be used during connection and stream creation by libp2p.
		host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

		// If no nick specified, use hostID as nick
		if *nick == "" {
			*nick = host.ID().Pretty()
		}

		// Initialize client config
		env := &clientConfig{
			ctx:      ctx,
			client:   smClient,
			room:     *room,
			serverID: info.ID,
			nick:     *nick,
			self:     host.ID(),
		}

		// Initialize chat UI
		ui := newUI(env)
		// Run UI, sync existing messages and start prompt
		err = ui.start()
		if err != nil {
			panic("Error initializing UI")
		}
	}
}
