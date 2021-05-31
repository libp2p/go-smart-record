package protocol

import (
	"context"
	"errors"
	"fmt"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	xr "github.com/libp2p/go-routing-language/syntax"
	pb "github.com/libp2p/go-smart-record/protocol/pb"
	"github.com/libp2p/go-smart-record/vm"
)

type SmartRecordServer interface {
	setProtocolHandler(network.StreamHandler)
}

// SmartRecordServer handles smart-record requests
type smartRecordServer struct {
	ctx       context.Context
	host      host.Host
	self      peer.ID
	vm        vm.Machine
	protocols []protocol.ID
}

// NewSmartRecordServer starts a smartRecordServer instance
func NewSmartRecordServer(ctx context.Context, h host.Host, options ...ServerOption) (SmartRecordServer, error) {
	return newSmartRecordServer(ctx, h, options...)
}

// setProtocolHandler sets new handler to the smart-record protocol
func (e *smartRecordServer) setProtocolHandler(h network.StreamHandler) {
	// For every announce protocol set this new handler.
	for _, p := range e.protocols {
		e.host.SetStreamHandler(p, h)
	}
}

func newSmartRecordServer(ctx context.Context, h host.Host, options ...ServerOption) (*smartRecordServer, error) {
	var cfg serverConfig
	if err := cfg.apply(append([]ServerOption{serverDefaults}, options...)...); err != nil {
		return nil, err
	}
	protocols := []protocol.ID{srProtocol}

	// Add host to assemblerContext
	cfg.assembler.Host = h

	vm, err := vm.NewVM(ctx, cfg.updateContext, cfg.assembler)
	if err != nil {
		return nil, err
	}
	// Start a smartRecordServer with an initialized VM.
	e := &smartRecordServer{
		ctx:       ctx,
		host:      h,
		self:      h.ID(),
		vm:        vm,
		protocols: protocols,
	}

	// Set streamhandler for smart-record protocol.
	e.setProtocolHandler(e.handleNewStream)

	return e, nil
}

// smartRecordHandler specifies the signature of functions that handle smart record messages.
type smartRecordHandler func(context.Context, peer.ID, *pb.Message) (*pb.Message, error)

func (e *smartRecordServer) handlerForMsgType(t pb.Message_MessageType) smartRecordHandler {
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

func (e *smartRecordServer) handleGet(ctx context.Context, p peer.ID, msg *pb.Message) (*pb.Message, error) {
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
	rb, err := vm.MarshalRecordValue(r)
	//rb, err := ir.Marshal(r)
	if err != nil {
		return nil, err
	}

	resp.Value = rb
	return resp, nil
}

func (e *smartRecordServer) handleUpdate(ctx context.Context, p peer.ID, msg *pb.Message) (*pb.Message, error) {

	k := msg.GetKey()
	if len(k) == 0 {
		return nil, errors.New("handleUpdate: no key was provided")
	}

	v := msg.GetValue()
	if len(k) == 0 {
		return nil, errors.New("handleUpdate: no value was provided")
	}

	// Unmarshal the record sent
	smrec, err := xr.UnmarshalJSON(v)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling record: %s", err)
	}
	rdict, ok := smrec.(xr.Dict)
	if !ok {
		return nil, fmt.Errorf("value sent is not a record. Won't update")
	}

	resp := &pb.Message{
		Type: msg.GetType(),
		Key:  k,
	}
	// Update in VM
	err = e.vm.Update(p, string(k), rdict)
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

func (e *smartRecordServer) handleQuery(ctx context.Context, p peer.ID, msg *pb.Message) (_ *pb.Message, err error) {
	// TODO: For now query is the same as get. We don't understand selectors yet.
	return e.handleGet(ctx, p, msg)
}
