package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/peerstore"
	"github.com/libp2p/go-smart-record/protocol"

	"github.com/multiformats/go-multiaddr"
)

func main() {
	ctx := context.Background()
	room := flag.String("room", "", "Select the chat room you want to connect to")
	dest := flag.String("d", "", "Destination multiaddr string")
	flag.Parse()
	if *dest == "" || *room == "" {
		fmt.Println("Specify a destination with -d and a -room to debug")
		return
	}

	fmt.Println("[*] Starting hosts")

	// Instantiating hosts
	h2, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}
	defer h2.Close()

	smClient, _ := protocol.NewSmartRecordClient(ctx, h2)

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
	h2.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
	// Get Record stored
	fmt.Println("[*] Getting updated record from peer")
	out, err := smClient.Get(ctx, *room, info.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("[*] Record:")
	for k, v := range *out {
		fmt.Println("Value for peer: ", k.String(), " - ", v)
	}

}
