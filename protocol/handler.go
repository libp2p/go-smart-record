package protocol

import (
	"context"
	"errors"
	"fmt"

	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-smart-record/ir"
	pb "github.com/libp2p/go-smart-record/protocol/pb"
)

// smartRecordHandler specifies the signature of functions that handle DHT messages.
type smartRecordHandler func(context.Context, peer.ID, *pb.Message) (*pb.Message, error)

func (e *SmartRecordEnv) handlerForMsgType(t pb.Message_MessageType) smartRecordHandler {
	switch t {
	case pb.Message_GET:
		return e.handleGet
	case pb.Message_UPDATE:
		return e.handleUpdate
	case pb.Message_QUERY:
		return e.handleQuery
	}

	return nil
}

func (e *SmartRecordEnv) handleGet(ctx context.Context, p peer.ID, msg *pb.Message) (*pb.Message, error) {
	k := msg.GetKey()
	if len(k) == 0 {
		return nil, errors.New("handleGet: no key was provided")
	}

	// setup response with same type as request.
	resp := &pb.Message{
		Type: msg.GetType(),
		Key:  k,
	}
	// Get record from VM
	r := e.vm.Get(string(k))
	// Marshal record
	rb, err := ir.Marshal(r)
	if err != nil {
		return nil, err
	}

	resp.Value = rb
	return resp, nil
}

func (e *SmartRecordEnv) handleUpdate(ctx context.Context, p peer.ID, msg *pb.Message) (*pb.Message, error) {

	k := msg.GetKey()
	if len(k) == 0 {
		return nil, errors.New("handleUpdate: no key was provided")
	}

	v := msg.GetValue()
	if len(k) == 0 {
		return nil, errors.New("handleUpdate: no value was provided")
	}

	// Unmarshal the record sent
	smrec, err := ir.Unmarshal(v)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling record: %s", err)
	}
	rdict, ok := smrec.(ir.Dict)
	if !ok {
		return nil, fmt.Errorf("value sent is not a record. Won't update")
	}

	resp := &pb.Message{
		Type: msg.GetType(),
		Key:  k,
	}

	// Get record from VM
	// Update in VM
	err = e.vm.Update(string(k), rdict)
	if err != nil {
		return nil, fmt.Errorf("failed updating dict: %s", err)
	}

	// NOTE: For now if the update is successful we just send an empty response
	// with the same key and the same type. We could reutrn a response type
	// saying if the Update was OK or KO but will avoid it for now. If the
	// update fails the stream will be closed with error so the other peer
	// will be notified that it failed.
	return resp, nil
}

func (e *SmartRecordEnv) handleQuery(ctx context.Context, p peer.ID, msg *pb.Message) (_ *pb.Message, err error) {
	// TODO: For now query is the same as get. We don't understand selectors yet.
	return e.handleGet(ctx, p, msg)
}
