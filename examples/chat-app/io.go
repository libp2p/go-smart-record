package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

// Read new messages from STDIN, and submit update with new message
// to SR server
func (c *clientConfig) readInput(outCh chan string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		// ReadString will block until the delimiter is entered
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			return
		}

		// remove the delimeter from the string
		input = strings.TrimSuffix(input, "\n")
		// sendMsg
		err = c.sendMsg(input)
		if err != nil {
			outCh <- "Error: Message failed to send\n"
		} else {
			outCh <- fmt.Sprintf("[%s]: %s\n", c.nick, input)
		}

	}
}

// Write in STDOUT new messages received in channel
func (c *clientConfig) writeOutput(outCh chan string) {
	fmt.Println("[*] Ready! You can start typing your messages :) ")
	w := bufio.NewWriter(os.Stdout)
	for input := range outCh {
		w.WriteString(input)
		w.Flush()
	}
}

// Get record, and sync new messages
func (c *clientConfig) syncMessages(outCh chan string) {
	ctx, cancel := context.WithTimeout(c.ctx, reqTimeout)
	// Get record from server.
	out, err := c.client.Get(ctx, c.room, c.serverID)
	cancel()
	if err != nil {
		//NOTE: What should we do if synchronizing fails? For now do nothing
		// printErr("sync error: %s", err)
	} else {
		// Process record to check if there are new messages.
		c.processSyncMessages(out, outCh)
	}
}
