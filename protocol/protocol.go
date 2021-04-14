package protocol

import (
	"context"
	"fmt"

	logging "github.com/ipfs/go-log"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/libp2p/go-smart-record/ir"
	pb "github.com/libp2p/go-smart-record/protocol/pb"
	"github.com/libp2p/go-smart-record/vm"
)

var log = logging.Logger("smart-records")

const (
	srProtocol protocol.ID = "/smart-record/0.0.1"
)

type SmartRecordEnv struct {
	ctx       context.Context
	proc      goprocess.Process
	host      host.Host
	self      peer.ID
	vm        vm.Machine
	protocols []protocol.ID

	senderManager *messageSenderImpl
}

// New creates a new DHT with the specified host and options.
// Please note that being connected to a DHT peer does not necessarily imply that it's also in the DHT Routing Table.
// If the Routing Table has more than "minRTRefreshThreshold" peers, we consider a peer as a Routing Table candidate ONLY when
// we successfully get a query response from it OR if it send us a query.
func New(ctx context.Context, h host.Host, options ...Option) (*SmartRecordEnv, error) {
	var cfg config
	if err := cfg.apply(append([]Option{defaults}, options...)...); err != nil {
		return nil, err
	}
	protocols := []protocol.ID{srProtocol}

	e := &SmartRecordEnv{
		ctx:       ctx,
		proc:      goprocessctx.WithContext(ctx),
		host:      h,
		self:      h.ID(),
		vm:        vm.NewVM(cfg.mergeContext, cfg.assembler),
		protocols: protocols,

		senderManager: &messageSenderImpl{
			host:      h,
			strmap:    make(map[peer.ID]*peerMessageSender),
			protocols: protocols,
		},
	}

	// Set streamhandler
	for _, p := range e.protocols {
		e.host.SetStreamHandler(p, e.handleNewStream)
	}

	// Create processes that will be listening to new records.
	// TODO: Is this really needed?
	return e, nil
}

func (e *SmartRecordEnv) Get(ctx context.Context, k string, p peer.ID) (*ir.Dict, error) {
	// Send a new request and wait for response
	req := &pb.Message{
		Type: pb.Message_GET,
		Key:  []byte(k),
	}
	resp, err := e.senderManager.SendRequest(ctx, p, req)
	if err != nil {
		return nil, err
	}
	rec, err := ir.Unmarshal(resp.GetValue())
	if err != nil {
		return nil, err
	}
	rdict, ok := rec.(ir.Dict)
	if !ok {
		return nil, fmt.Errorf("received value has wrong type")
	}
	// NOTE: Here we are returning a disassembled record. We may want
	// to assemble it here and return a *base.Record.
	return &rdict, nil

}

func (e *SmartRecordEnv) Update(ctx context.Context, k string, p peer.ID, rec ir.Dict) error {
	// Send a new request and wait for response
	recB, err := ir.Marshal(rec)
	if err != nil {
		return err
	}
	req := &pb.Message{
		Type:  pb.Message_UPDATE,
		Key:   []byte(k),
		Value: recB,
	}
	resp, err := e.senderManager.SendRequest(ctx, p, req)
	if err != nil {
		return err
	}
	// NOTE: We are only sending a message now if the update is successful. If we
	// don't receive a response it means that the update has failed at the other end
	// for some reason
	if resp == nil {
		return fmt.Errorf("update request failed, no response received")
	}

	return nil
}

func (e *SmartRecordEnv) Query(ctx context.Context, k string, p peer.ID) (*ir.Dict, error) {
	// NOTE: For now Query and Get are the same because we don't
	// understand selectors (yet)
	return e.Get(ctx, k, p)
}
