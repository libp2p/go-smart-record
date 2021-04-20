package protocol

import (
	"context"
	"fmt"

	logging "github.com/ipfs/go-log"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-smart-record/ir"
	pb "github.com/libp2p/go-smart-record/protocol/pb"
	"github.com/libp2p/go-smart-record/vm"
	sm "github.com/libp2p/go-smart-record/vm"
)

var log = logging.Logger("smart-records")

const (
	srProtocol protocol.ID = "/smart-record/0.0.1"
)

// SmartRecordManager interface to manage smart records
type SmartRecordManager interface {
	Get(ctx context.Context, k string, p peer.ID) (*vm.RecordValue, error)
	Update(ctx context.Context, k string, p peer.ID, rec ir.Dict) error
	// NOTE: we won't support queries until we figure out selectors
	// Query(ctx context.Context, k string, p peer.ID, selector ir.Dict) (*ir.Dict, error)
}

// smartRecordManager handles the exchange of messages
// with peers to interact with smart-records
type smartRecordManager struct {
	ctx       context.Context
	proc      goprocess.Process
	host      host.Host
	self      peer.ID
	vm        vm.Machine
	protocols []protocol.ID

	senderManager *messageSenderImpl
}

func NewSmartRecordManager(ctx context.Context, h host.Host, options ...Option) (SmartRecordManager, error) {
	return newSmartRecordManager(ctx, h, options...)
}

func newSmartRecordManager(ctx context.Context, h host.Host, options ...Option) (*smartRecordManager, error) {
	var cfg config
	if err := cfg.apply(append([]Option{defaults}, options...)...); err != nil {
		return nil, err
	}
	protocols := []protocol.ID{srProtocol}

	e := &smartRecordManager{
		ctx:       ctx,
		proc:      goprocessctx.WithContext(ctx),
		host:      h,
		self:      h.ID(),
		vm:        sm.NewVM(cfg.mergeContext, cfg.assembler),
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

	return e, nil
}

func (e *smartRecordManager) Get(ctx context.Context, k string, p peer.ID) (*vm.RecordValue, error) {
	// Send a new request and wait for response
	req := &pb.Message{
		Type: pb.Message_GET,
		Key:  []byte(k),
	}
	resp, err := e.senderManager.SendRequest(ctx, p, req)
	if err != nil {
		return nil, err
	}
	//rec, err := ir.Unmarshal(resp.GetValue())
	rv, err := vm.UnmarshalRecordValue(resp.GetValue())
	if err != nil {
		return nil, err
	}

	return &rv, nil

}

func (e *smartRecordManager) Update(ctx context.Context, k string, p peer.ID, rec ir.Dict) error {
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

/*
func (e *smartRecordManager) Query(ctx context.Context, k string, p peer.ID, selector ir.Dict) (*ir.Dict, error) {
	// NOTE: For now Query and Get are the same because we don't
	// understand selectors (yet)
	return e.Get(ctx, k, p)
}
*/
