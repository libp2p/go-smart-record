package main

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/libp2p/go-libp2p/core/peer"
	xr "github.com/libp2p/go-routing-language/syntax"
	"github.com/libp2p/go-smart-record/protocol"
	"github.com/libp2p/go-smart-record/vm"
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

// Genereates the data model of messages for the chat application
func (c *clientConfig) generateChatMessage(msg string) xr.Dict {
	c.lk.Lock()
	// Increase seqID
	c.seqId++
	defer c.lk.Unlock()
	// Message data
	d := xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.Int{Int: big.NewInt(c.seqId)}, Value: xr.String{Value: msg}},
		},
	}

	// Include message data into a wrapper with seqID for synchronization
	return xr.Dict{
		Pairs: xr.Pairs{
			xr.Pair{Key: xr.String{Value: "nick"}, Value: xr.String{Value: c.nick}},
			xr.Pair{Key: xr.String{Value: "msgs"}, Value: d},
		},
	}
}

// Submit new message to smart record server through an update
func (c *clientConfig) sendMsg(text string) error {
	msg := c.generateChatMessage(text)
	ctx, cancel := context.WithTimeout(c.ctx, reqTimeout)
	// Send message to record
	err := c.client.Update(ctx, c.room, c.serverID, msg, msgTTL)
	cancel()
	return err
}

func (c *clientConfig) processSyncMessages(out *vm.RecordValue, outCh chan string) {
	//  TODO: Check that type casts are correct throughout all the method. If not throw error
	syncMsgs := make(map[int64][]*syncUpdate) //seqID - nick - msg
	ids := make([]int, 0)
	var tmpMax int64 = -1
	update := false

	// For every peer in record
	for k, v := range *out {
		// Already have my own messages, no need to process them
		if k == c.self {
			continue
		}

		// Get peer's nick (if any)
		nickNode := v.Get(xr.String{Value: "nick"})
		nick := k.Pretty()
		if nickNode != nil {
			nick = nickNode.(xr.String).Value
		}

		msgs := v.Get(xr.String{Value: "msgs"})
		mdict, _ := msgs.(xr.Dict)

		// For all messages in peer
		for _, pv := range mdict.Pairs {
			ki := pv.Key.(xr.Int)
			i := ki.Int64()

			// If message has a seqID below the one I keep, it means I haven't seen it
			if i > c.syncId {
				// Add id as seqId to track and sync at the end
				ids = append(ids, int(i))
				// Append the message for update
				syncMsgs[i] = append(syncMsgs[i], &syncUpdate{nick, pv.Value.(xr.String).Value})
				// Update the max sequence number that I've seen so far
				if i > tmpMax {
					tmpMax = i
				}
				// Flag that at least one new message has been found
				update = true
			}
		}

	}
	// Once all the messages have been processed update sequenceIds.
	if update {
		c.lk.Lock()
		c.seqId = tmpMax
		c.syncId = tmpMax
		c.lk.Unlock()
	}

	// Sort ids seen
	sort.Ints(ids)

	// For each id print sorted messages
	for _, i := range ids {
		// Print every new message in UI
		for _, m := range syncMsgs[int64(i)] {
			outCh <- fmt.Sprintf("[%s]: %s\n", m.nick, m.msg)
		}
	}

}
