package protocol

import (
	"context"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	swarmt "github.com/libp2p/go-libp2p-swarm/testing"
	bhost "github.com/libp2p/go-libp2p/p2p/host/basic"
	xr "github.com/libp2p/go-routing-language/syntax"
)

// TTL for updates in test cases
var ttl = 2 * time.Second

// Use small gcPeriod in server VM for tests
var gcPeriod = 1 * time.Second

var in1 = xr.Dict{
	Pairs: xr.Pairs{
		xr.Pair{Key: xr.String{Value: "key"}, Value: xr.String{Value: "234"}},
		xr.Pair{Key: xr.String{Value: "QmXBar"}, Value: xr.String{Value: "/ip4/multiaddr1"}},
		xr.Pair{Key: xr.String{Value: "QmXFor"}, Value: xr.String{Value: "/ip4/multiaddr2"}},
	},
}
var in2 = xr.Dict{
	Pairs: xr.Pairs{
		xr.Pair{Key: xr.String{Value: "key"}, Value: xr.String{Value: "234"}},
		xr.Pair{Key: xr.String{Value: "QmXBar2"}, Value: xr.String{Value: "/ip4/multiaddr3"}},
		xr.Pair{Key: xr.String{Value: "QmXFoo2"}, Value: xr.String{Value: "/ip4/multiaddr4"}},
	},
}
var in = xr.Dict{
	Pairs: xr.Pairs{
		xr.Pair{Key: xr.String{Value: "key"}, Value: xr.String{Value: "234"}},
		xr.Pair{Key: xr.String{Value: "QmXBar"}, Value: xr.String{Value: "/ip4/multiaddr1"}},
		xr.Pair{Key: xr.String{Value: "QmXFor"}, Value: xr.String{Value: "/ip4/multiaddr2"}},
		xr.Pair{Key: xr.String{Value: "QmXBar2"}, Value: xr.String{Value: "/ip4/multiaddr3"}},
		xr.Pair{Key: xr.String{Value: "QmXFoo2"}, Value: xr.String{Value: "/ip4/multiaddr4"}},
	},
}

func setupServer(ctx context.Context, t *testing.T) *smartRecordServer {

	s, err := newSmartRecordServer(
		ctx,
		bhost.New(swarmt.GenSwarm(t, ctx, swarmt.OptDisableReuseport)),
		[]ServerOption{VMGcPeriod(gcPeriod)}...,
	)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func setupClient(ctx context.Context, t *testing.T) *smartRecordClient {

	c, err := newSmartRecordClient(
		ctx,
		bhost.New(swarmt.GenSwarm(t, ctx, swarmt.OptDisableReuseport)),
	)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func connect(ctx context.Context, t *testing.T, h1 host.Host, h2 host.Host) {
	if err := h1.Connect(ctx, *host.InfoFromHost(h2)); err != nil {
		t.Fatal(err)
	}
}
func TestEmptyUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setupClient(ctx, t)
	s := setupServer(ctx, t)
	connect(ctx, t, c.host, s.host)

	k := "234"

	// Update record
	err := c.Update(ctx, k, s.host.ID(), in1, ttl)
	if err != nil {
		t.Fatal(err)
	}

	// Get record
	out, err := c.Get(ctx, k, s.host.ID())
	if err != nil {
		panic(err)
	}
	d := (*out)[c.host.ID()]
	if !xr.IsEqual(in1, *d) {
		t.Fatal("end-to-end update in empty key failed", in1, *out)
	}

	// Check if TTL set successfully.
	time.Sleep(3 * time.Second)
	out, err = c.Get(ctx, k, s.host.ID())
	if err != nil {
		panic(err)
	}
	if len(*out) != 0 {
		t.Fatal("TTL was not set successfully. The record didn't expire.", *out)
	}

}

func TestLocalEmptyUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s := setupServer(ctx, t)

	k := "234"

	// Put local
	err := s.UpdateLocal(k, s.host.ID(), in1, ttl)
	if err != nil {
		t.Fatal(err)
	}

	// Get record
	out := s.GetLocal(k)
	if err != nil {
		panic(err)
	}
	d := out[s.host.ID()]
	if !xr.IsEqual(in1, *d) {
		t.Fatal("local update in empty key failed", in1, out)
	}

}

func TestUpdateSameKeyDifferentPeers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c1 := setupClient(ctx, t)
	c2 := setupClient(ctx, t)
	s := setupServer(ctx, t)
	connect(ctx, t, c1.host, s.host)
	connect(ctx, t, c2.host, s.host)

	k := "234"

	// Update record
	err := c1.Update(ctx, k, s.host.ID(), in1, ttl)
	if err != nil {
		t.Fatal(err)
	}
	err = c2.Update(ctx, k, s.host.ID(), in2, ttl)
	if err != nil {
		t.Fatal(err)
	}

	// Get record
	out, err := c1.Get(ctx, k, s.host.ID())
	if err != nil {
		panic(err)
	}
	d1 := (*out)[c1.host.ID()]
	if !xr.IsEqual(in1, *d1) {
		t.Fatal("end-to-end update in empty key for client1 failed", in1, *out)
	}
	d2 := (*out)[c2.host.ID()]
	if !xr.IsEqual(in2, *d2) {
		t.Fatal("end-to-end update in empty key for client2 failed", in1, *out)
	}

}

func TestUpdateSameKeySamePeer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c := setupClient(ctx, t)
	s := setupServer(ctx, t)
	connect(ctx, t, c.host, s.host)

	k := "234"

	// Update record
	err := c.Update(ctx, k, s.host.ID(), in1, ttl)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Update(ctx, k, s.host.ID(), in2, ttl)
	if err != nil {
		t.Fatal(err)
	}

	// Get record
	out, err := c.Get(ctx, k, s.host.ID())
	if err != nil {
		panic(err)
	}
	d := (*out)[c.host.ID()]
	if !xr.IsEqual(in, *d) {
		t.Fatal("end-to-end existing key for client1 failed", in1, *out)
	}
}

func TestParallelRequests(t *testing.T) {
	//TODO
}
