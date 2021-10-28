package main

import (
	"context"
	"math/big"
	"sync"

	peer "github.com/libp2p/go-libp2p-core/peer"
	xr "github.com/libp2p/go-routing-language/syntax"
	"github.com/libp2p/go-smart-record/protocol"
)

type clientConfig struct {
	lk       sync.Mutex
	ctx      context.Context
	seqId    int64 // Keeps track of my sequence ID
	syncId   int64 // Determines seqID from last sync
	client   protocol.SmartRecordClient
	room     string
	serverID peer.ID
	self     peer.ID
	nick     string
}

type syncUpdate struct {
	nick string
	msg  string
}

// Genereates the data model for messages for the chat application
func (e *clientConfig) generateChatMessage(msg string) xr.Dict {
	e.lk.Lock()
	// Increase seqID
	e.seqId++
	defer e.lk.Unlock()
	// Message data
	d := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.Int{Int: big.NewInt(e.seqId)}, Value: xr.String{Value: msg}},
		},
	}

	// Include message data into a wrapper with seqID for synchronization
	return xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "nick"}, Value: xr.String{Value: e.nick}},
			xr.Pair{Key: xr.String{Value: "msgs"}, Value: d},
		},
	}
}
