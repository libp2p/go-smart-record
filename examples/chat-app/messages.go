package main

import (
	"context"
	"math/big"
	"sync"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-smart-record/ir"
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
func (e *clientConfig) generateChatMessage(msg string) ir.Dict {
	e.lk.Lock()
	// Increase seqID
	e.seqId++
	defer e.lk.Unlock()
	// Message data
	d := ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.Int{big.NewInt(e.seqId)}, Value: ir.String{Value: msg}},
		},
	}

	// Include message data into a wrapper with seqID for synchronization
	return ir.Dict{
		Pairs: ir.Pairs{
			ir.Pair{Key: ir.String{Value: "nick"}, Value: ir.String{Value: e.nick}},
			ir.Pair{Key: ir.String{Value: "msgs"}, Value: d},
		},
	}
}
