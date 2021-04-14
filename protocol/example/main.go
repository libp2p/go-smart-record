package main

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-smart-record/ir"
	"github.com/libp2p/go-smart-record/ir/base"
	env "github.com/libp2p/go-smart-record/protocol"
)

func main() {
	ctx := context.Background()

	fmt.Println("[*] Starting hosts")

	// Variable to host h2 smartRecord manager used to
	// expose its interface.
	var sm env.SmartRecordManager
	// Option to create smart record in hosts
	sr := func(h host.Host) (env.SmartRecordManager, error) {
		var err error
		sm, err = env.NewSmartRecordManager(ctx, h)
		return sm, err
	}
	smartRecordsOpt := libp2p.SmartRecord(sr)

	// Instantiating hosts
	h1, err := libp2p.New(ctx, smartRecordsOpt)
	if err != nil {
		panic(err)
	}
	h2, err := libp2p.New(ctx, smartRecordsOpt)
	if err != nil {
		panic(err)
	}
	defer h1.Close()
	defer h2.Close()

	// Wait until hosts are ready
	time.Sleep(3 * time.Second)

	fmt.Println("[*] Connecting peers")
	// Connect h1-h2
	err = DialOtherPeer(ctx, h1, *host.InfoFromHost(h2))
	if err != nil {
		panic(err)
	}

	// Record to update
	fmt.Println("[*] Updating new record")
	ind := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "key"}, Value: ir.String{Value: "234"}},
			ir.Pair{Key: ir.String{Value: "fff"}, Value: ir.String{Value: "ff2"}},
		},
	}
	k := "234"
	r := base.Record{Key: k, User: ind}
	in := r.Disassemble()

	// Update record
	err = sm.Update(ctx, k, h1.ID(), in)
	if err != nil {
		panic(err)
	}
	fmt.Println("[*] Update successful")

	// Get Record stored
	fmt.Println("[*] Getting updated record from peer")
	out, err := sm.Get(ctx, k, h1.ID())
	if err != nil {
		panic(err)
	}

	fmt.Println("[*] It worked")
	var w bytes.Buffer
	out.WritePretty(&w)
	fmt.Println(w.String())
}

// DialOtherPeers connects to a set of peers in the experiment.
func DialOtherPeer(ctx context.Context, self host.Host, ai peer.AddrInfo) error {
	if err := self.Connect(ctx, ai); err != nil {
		return fmt.Errorf("Error while dialing peer %v: %w", ai.Addrs, err)
	}
	return nil
}
