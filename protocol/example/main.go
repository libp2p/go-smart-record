package main

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-smart-record/protocol"
	"github.com/libp2p/go-smart-record/xr"
)

func main() {
	ctx := context.Background()

	fmt.Println("[*] Starting hosts")

	// Instantiating hosts
	h1, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}
	h2, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}
	defer h1.Close()
	defer h2.Close()

	// Wait until hosts are ready
	time.Sleep(1 * time.Second)

	fmt.Println("[*] Connecting peers")
	// Connect h1-h2
	err = DialOtherPeer(ctx, h1, *host.InfoFromHost(h2))
	if err != nil {
		panic(err)
	}

	_, _ = protocol.NewSmartRecordServer(ctx, h1)
	smClient, _ := protocol.NewSmartRecordClient(ctx, h2)

	// Record to update
	fmt.Println("[*] Updating new record")
	in1 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "key"}, Value: xr.String{Value: "234"}},
			xr.Pair{Key: xr.String{Value: "QmXBar"}, Value: xr.String{Value: "/ip4/multiaddr1"}},
			xr.Pair{Key: xr.String{Value: "QmXFor"}, Value: xr.String{Value: "/ip4/multiaddr2"}},
		},
	}
	in2 := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "key"}, Value: xr.String{Value: "234"}},
			xr.Pair{Key: xr.String{Value: "QmXBar2"}, Value: xr.String{Value: "/ip4/multiaddr3"}},
			xr.Pair{Key: xr.String{Value: "QmXFoo2"}, Value: xr.String{Value: "/ip4/multiaddr4"}},
		},
	}
	k := "234"

	// Update record with 60 seconds of TTL
	err = smClient.Update(ctx, k, h1.ID(), in1, uint64(60))
	if err != nil {
		panic(err)
	}
	fmt.Println("[*] Update 1 successful")

	// Update record with 60 seconds of TTL
	err = smClient.Update(ctx, k, h1.ID(), in2, uint64(60))
	if err != nil {
		panic(err)
	}
	fmt.Println("[*] Update 2 successful")

	// Get Record stored
	fmt.Println("[*] Getting updated record from peer")
	out, err := smClient.Get(ctx, k, h1.ID())
	if err != nil {
		panic(err)
	}

	fmt.Println("[*] It worked")
	fmt.Println("Record Key", k)
	for k, v := range *out {
		fmt.Println("Value for peer: ", k.String(), " - ", v)
	}

}

// DialOtherPeers connects to a set of peers in the experiment.
func DialOtherPeer(ctx context.Context, self host.Host, ai peer.AddrInfo) error {
	if err := self.Connect(ctx, ai); err != nil {
		return fmt.Errorf("Error while dialing peer %v: %w", ai.Addrs, err)
	}
	return nil
}
