package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	"github.com/gdamore/tcell/v2"
	xr "github.com/libp2p/go-routing-language/syntax"
	"github.com/libp2p/go-smart-record/vm"
	"github.com/rivo/tview"
)

type chatUI struct {
	env *clientConfig
	app *tview.Application

	msgW    io.Writer
	inputCh chan string
	doneCh  chan struct{}
}

func newUI(env *clientConfig) *chatUI {
	app := tview.NewApplication()

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf("Room: %s", env.room))

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	// an input field for typing messages into
	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(env.nick + " > ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	// the done func is called when the user hits enter, or tabs out of the field
	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}
		line := input.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}

		// bail if requested
		if line == "/quit" {
			app.Stop()
			return
		}

		// send the line onto the input chan and reset the field text
		inputCh <- line
		input.SetText("")
	})

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgBox, 0, 1, false)
		// flex is a vertical box with the chatPanel on top and the input field at the bottom.

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)
	return &chatUI{
		env:     env,
		app:     app,
		msgW:    msgBox,
		inputCh: inputCh,
		doneCh:  make(chan struct{}, 1),
	}
}

func (ui *chatUI) start() error {
	// Start event handler
	go ui.handleEvents()
	defer ui.end()
	// Run UI
	return ui.app.Run()
}

// end signals the event loop to exit gracefully
func (ui *chatUI) end() {
	ui.doneCh <- struct{}{}
}

// sends new message to smart-record server
func (ui *chatUI) sendNewMessage(text string) {
	msg := ui.env.generateChatMessage(text)
	ctx, cancel := context.WithTimeout(ui.env.ctx, reqTimeout)
	// Send message to record. Set TTL to messages
	err := ui.env.client.Update(ctx, ui.env.room, ui.env.serverID, msg, msgTTL)
	cancel()
	if err != nil {
		printErr("publish error: %s", err)
	}
	// when the user types in a line, and the record updates successfully print to the message window
	ui.displaySelfMessage(text)
}

// The chatUI orchestrates the handling of events from input and syncing messages.
func (ui *chatUI) handleEvents() {
	// defer msgSyncTicker.Stop()
	msgSyncTicker := time.NewTicker(syncTime)
	for {
		select {
		case input := <-ui.inputCh:
			// If new input message send update to record.
			ui.sendNewMessage(input)

		case <-msgSyncTicker.C:
			msgSyncTicker.Stop()
			// We periodically fetch the smart-record to check if there are new messages.
			ui.syncMessages()
			msgSyncTicker = time.NewTicker(syncTime)

		case <-ui.env.ctx.Done():
			return
		case <-ui.doneCh:
			return
		}
	}
}

// Get record, and sync new messages
func (ui *chatUI) syncMessages() {
	ctx, cancel := context.WithTimeout(ui.env.ctx, reqTimeout)
	// Get record from server.
	out, err := ui.env.client.Get(ctx, ui.env.room, ui.env.serverID)
	cancel()
	if err != nil {
		//NOTE: What should we do if synchronizing fails? For now do nothing
		printErr("sync error: %s", err)
	} else {
		// Process record to check if there are new messages.
		ui.processSyncMessages(out)
	}
}

func (ui *chatUI) processSyncMessages(out *vm.RecordValue) {
	//  TODO: Check that type casts are correct throughout all the method. If not throw error
	syncMsgs := make(map[int64][]*syncUpdate, 0) //seqID - nick - msg
	ids := make([]int, 0)
	var tmpMax int64 = -1
	update := false

	// For every peer in record
	for k, v := range *out {
		// Already have my own messages, no need to process them
		// NOTE: We could perform an initial sync to gather our previous messages
		// and print them also in case we disconnected.
		if k == ui.env.self {
			continue
		}

		// Get peer's nick (if any)
		nickNode := v.Get(xr.String{Value: "nick"})
		nick := k.Pretty()
		if nickNode != nil {
			nick = nickNode.(xr.String).Value
		}

		msgs := v.Get(xr.String{Value: "msgs"})
		mdict, ok := msgs.(xr.Dict)
		if !ok {
			printErr("sync error: dict of messages not stored inrecord")
		}

		// For all messages in peer
		for _, pv := range mdict.Pairs {
			ki := pv.Key.(xr.Int)
			i := ki.Int64()

			// If message has a seqID below the one I keep, it means I haven't seen it
			if i > ui.env.syncId {
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
		ui.env.lk.Lock()
		ui.env.seqId = tmpMax
		ui.env.syncId = tmpMax
		ui.env.lk.Unlock()
	}

	// Sort ids seen
	sort.Ints(ids)

	// For each id print sorted messages
	for _, i := range ids {
		// Print every new message in UI
		for _, m := range syncMsgs[int64(i)] {
			ui.displayChatMessage(m.nick, m.msg)
		}
	}

}

func (ui *chatUI) displaySelfMessage(msg string) {
	fmt.Fprintf(ui.msgW, "[yellow]<%s>: [-]%s\n", ui.env.nick, msg)
}

func (ui *chatUI) displayChatMessage(nick string, msg string) {
	fmt.Fprintf(ui.msgW, "[blue]<%s>: [-]%s\n", nick, msg)
}

func printErr(m string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, m, args...)
}
