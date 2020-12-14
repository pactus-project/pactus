package sync

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/network"
)

type NetworkAPI interface {
	Start() error
	Stop()
	PublishMessage(msg *message.Message) error
}

type networkAPI struct {
	ctx            context.Context
	selfAddress    crypto.Address
	selfID         peer.ID
	generalTopic   *pubsub.Topic
	txTopic        *pubsub.Topic
	blockTopic     *pubsub.Topic
	consensusTopic *pubsub.Topic
	generalSub     *pubsub.Subscription
	txSub          *pubsub.Subscription
	blockSub       *pubsub.Subscription
	consensusSub   *pubsub.Subscription
	parsMessageFn  func(data []byte, from peer.ID)
}

func newNetworkAPI(
	ctx context.Context,
	selfAddress crypto.Address,
	net *network.Network,
	parsMessageFn func(data []byte, from peer.ID)) (*networkAPI, error) {
	generalTopic, err := net.JoinTopic("general")
	if err != nil {
		return nil, err
	}
	generalSub, err := generalTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	txTopic, err := net.JoinTopic("tx")
	if err != nil {
		return nil, err
	}
	txSub, err := txTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	blockTopic, err := net.JoinTopic("block")
	if err != nil {
		return nil, err
	}
	blockSub, err := blockTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	consensusTopic, err := net.JoinTopic("consensus")
	if err != nil {
		return nil, err
	}
	consensusSub, err := consensusTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	return &networkAPI{
		ctx:            ctx,
		selfID:         net.ID(),
		selfAddress:    selfAddress,
		txTopic:        txTopic,
		txSub:          txSub,
		blockSub:       blockSub,
		blockTopic:     blockTopic,
		generalTopic:   generalTopic,
		generalSub:     generalSub,
		consensusTopic: consensusTopic,
		consensusSub:   consensusSub,
		parsMessageFn:  parsMessageFn,
	}, nil
}

func (api *networkAPI) Start() error {
	go api.txLoop()
	go api.blockLoop()
	go api.generalLoop()
	go api.consensusLoop()

	return nil
}

func (api *networkAPI) Stop() {
	api.txTopic.Close()
	api.txSub.Cancel()
	api.blockTopic.Close()
	api.blockSub.Cancel()
	api.generalTopic.Close()
	api.generalSub.Cancel()
	api.consensusTopic.Close()
	api.consensusSub.Cancel()
}

func (api *networkAPI) parsMessage(m *pubsub.Message) {
	// only forward messages delivered by others
	if m.ReceivedFrom == api.selfID {
		return
	}

	api.parsMessageFn(m.Data, m.ReceivedFrom)
}

func (api *networkAPI) PublishMessage(msg *message.Message) error {
	msg.Initiator = api.selfAddress
	topic := api.topic(msg)
	bs, _ := msg.MarshalCBOR()
	return topic.Publish(api.ctx, bs)
}

func (api *networkAPI) txLoop() {
	for {
		m, err := api.txSub.Next(api.ctx)
		if err != nil {
			logger.Error("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) blockLoop() {
	for {
		m, err := api.blockSub.Next(api.ctx)
		if err != nil {
			logger.Error("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) generalLoop() {
	for {
		m, err := api.generalSub.Next(api.ctx)
		if err != nil {
			logger.Error("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) consensusLoop() {
	for {
		m, err := api.consensusSub.Next(api.ctx)
		if err != nil {
			logger.Error("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) topic(msg *message.Message) *pubsub.Topic {
	switch msg.PayloadType() {

	case message.PayloadTypeSalam,
		message.PayloadTypeHeartBeat:
		return api.generalTopic

	case message.PayloadTypeBlocksReq,
		message.PayloadTypeBlocks:
		return api.blockTopic

	case message.PayloadTypeTxsReq,
		message.PayloadTypeTxs:
		return api.txTopic

	case message.PayloadTypeProposalReq,
		message.PayloadTypeProposal,
		message.PayloadTypeVote,
		message.PayloadTypeVoteSet:
		return api.consensusTopic
	default:
		panic("Invalid topic")
	}
}
