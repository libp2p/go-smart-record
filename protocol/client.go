package protocol

import (
	"context"
	"fmt"

	logging "github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	pb "github.com/libp2p/go-smart-record/protocol/pb"
	"github.com/libp2p/go-smart-record/vm"
	"github.com/libp2p/go-smart-record/xr"
)

var log = logging.Logger("smart-records")

// SmartRecordClient sends smart-record requesets to other peers.
type SmartRecordClient interface {
	Get(ctx context.Context, k string, p peer.ID) (*vm.RecordValue, error)
	Update(ctx context.Context, k string, p peer.ID, rec xr.Dict, ttl uint64) error
	// NOTE: we won't support queries until we figure out selectors
	// Query(ctx context.Context, k string, p peer.ID, selector ir.Dict) (*ir.Dict, error)
}

// smartRecordClient is responsible for sending smart-record
// requests to other peers.
type smartRecordClient struct {
	ctx       context.Context
	host      host.Host
	self      peer.ID
	protocols []protocol.ID

	senderManager *messageSenderImpl
}

// NewSmartRecordClient starts a smartRecordClient instance
func NewSmartRecordClient(ctx context.Context, h host.Host, options ...ClientOption) (SmartRecordClient, error) {
	return newSmartRecordClient(ctx, h, options...)
}

func newSmartRecordClient(ctx context.Context, h host.Host, options ...ClientOption) (*smartRecordClient, error) {
	var cfg clientConfig
	if err := cfg.apply(append([]ClientOption{clientDefaults}, options...)...); err != nil {
		return nil, err
	}
	protocols := []protocol.ID{srProtocol}

	// Start a smartRecordClient
	e := &smartRecordClient{
		ctx:       ctx,
		host:      h,
		self:      h.ID(),
		protocols: protocols,

		senderManager: &messageSenderImpl{
			host:      h,
			strmap:    make(map[peer.ID]*peerMessageSender),
			protocols: protocols,
		},
	}

	return e, nil
}

func (e *smartRecordClient) Get(ctx context.Context, k string, p peer.ID) (*vm.RecordValue, error) {
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

func (e *smartRecordClient) Update(ctx context.Context, k string, p peer.ID, rec xr.Dict, ttl uint64) error {
	// Send a new request and wait for response
	recB, err := xr.Marshal(rec)
	if err != nil {
		return err
	}
	req := &pb.Message{
		Type:  pb.Message_UPDATE,
		Key:   []byte(k),
		Value: recB,
		TTL:   ttl,
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
