package protocol

import (
	"bufio"
	"io"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-msgio"
	"github.com/libp2p/go-msgio/protoio"
	pb "github.com/libp2p/go-smart-record/protocol/pb"
)

// Idle time before the stream is closed
var streamIdleTimeout = 1 * time.Minute

// Timeout to wait for a response after a request is sent
var readMessageTimeout = 10 * time.Second

// The Protobuf writer performs multiple small writes when writing a message.
// We need to buffer those writes, to make sure that we're not sending a new
// packet for every single write.
type bufferedDelimitedWriter struct {
	*bufio.Writer
	protoio.WriteCloser
}

var writerPool = sync.Pool{
	New: func() interface{} {
		w := bufio.NewWriter(nil)
		return &bufferedDelimitedWriter{
			Writer:      w,
			WriteCloser: protoio.NewDelimitedWriter(w),
		}
	},
}

func writeMsg(w io.Writer, mes *pb.Message) error {
	bw := writerPool.Get().(*bufferedDelimitedWriter)
	bw.Reset(w)
	err := bw.WriteMsg(mes)
	if err == nil {
		err = bw.Flush()
	}
	bw.Reset(nil)
	writerPool.Put(bw)
	return err
}

// handleNewStream implements the network.StreamHandler
func (e *smartRecordManager) handleNewStream(s network.Stream) {
	if e.handleNewMessage(s) {
		// If we exited without error, close gracefully.
		_ = s.Close()
	} else {
		// otherwise, send an error.
		_ = s.Reset()
	}
}

// Returns true on orderly completion of writes (so we can Close the stream conveniently).
func (e *smartRecordManager) handleNewMessage(s network.Stream) bool {
	ctx := e.ctx
	r := msgio.NewVarintReaderSize(s, network.MessageSizeMax)

	mPeer := s.Conn().RemotePeer()

	timer := time.AfterFunc(streamIdleTimeout, func() { _ = s.Reset() })
	defer timer.Stop()

	for {
		var req pb.Message
		msgbytes, err := r.ReadMsg()
		if err != nil {
			r.ReleaseMsg(msgbytes)
			if err == io.EOF {
				return true
			}
			return false
		}
		err = req.Unmarshal(msgbytes)
		r.ReleaseMsg(msgbytes)
		if err != nil {
			return false
		}

		timer.Reset(streamIdleTimeout)

		handler := e.handlerForMsgType(req.GetType())
		if handler == nil {
			return false
		}

		resp, err := handler(ctx, mPeer, &req)
		if err != nil {
			return false
		}

		if resp == nil {
			continue
		}

		// send out response msg
		err = writeMsg(s, resp)
		if err != nil {
			return false
		}

	}
}
