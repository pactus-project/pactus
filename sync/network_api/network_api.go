package network_api

import (
	"context"
	"sync"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sync/firewall"

	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/network"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

type ParsMessageFn = func(msg *message.Message, from peer.ID)

type networkAPI struct {
	ctx            context.Context
	selfID         peer.ID
	net            *network.Network
	firewall       *firewall.Firewall
	generalTopic   *pubsub.Topic
	downloadTopic  *pubsub.Topic
	dataTopic      *pubsub.Topic
	consensusTopic *pubsub.Topic
	generalSub     *pubsub.Subscription
	downloadSub    *pubsub.Subscription
	dataSub        *pubsub.Subscription
	consensusSub   *pubsub.Subscription
	parsFn         ParsMessageFn
	wg             sync.WaitGroup
}

func NewNetworkAPI(
	ctx context.Context,
	net *network.Network,
	firewall *firewall.Firewall,
	parsFn ParsMessageFn) (NetworkAPI, error) {
	generalTopic, err := net.JoinTopic("general")
	if err != nil {
		return nil, err
	}
	generalSub, err := generalTopic.Subscribe()
	if err != nil {
		return nil, err
	}
	dataTopic, err := net.JoinTopic("data")
	if err != nil {
		return nil, err
	}
	dataSub, err := dataTopic.Subscribe()
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

	api := &networkAPI{
		ctx:            ctx,
		selfID:         net.ID(),
		net:            net,
		firewall:       firewall,
		downloadTopic:  nil,
		downloadSub:    nil,
		dataSub:        dataSub,
		dataTopic:      dataTopic,
		generalTopic:   generalTopic,
		generalSub:     generalSub,
		consensusTopic: consensusTopic,
		consensusSub:   consensusSub,
		parsFn:         parsFn,
	}

	return api, nil
}

func (api *networkAPI) Start() error {

	api.wg.Add(1)
	go api.dataLoop()

	api.wg.Add(1)
	go api.generalLoop()

	api.wg.Add(1)
	go api.consensusLoop()

	return nil
}

func (api *networkAPI) Stop() {
	api.LeaveDownloadTopic()

	api.dataTopic.Close()
	api.dataSub.Cancel()

	api.generalTopic.Close()
	api.generalSub.Cancel()

	api.consensusTopic.Close()
	api.consensusSub.Cancel()

	api.wg.Wait()

	api.downloadTopic = nil
	api.downloadSub = nil
}

func (api *networkAPI) JoinDownloadTopic() error {
	if api.downloadSub != nil {
		return nil
	}

	downloadTopic, err := api.net.JoinTopic("download")
	if err != nil {
		return err
	}
	downloadSub, err := downloadTopic.Subscribe()
	if err != nil {
		return err
	}
	api.downloadTopic = downloadTopic
	api.downloadSub = downloadSub

	api.wg.Add(1)
	go api.downloadLoop()

	return nil
}
func (api *networkAPI) LeaveDownloadTopic() {
	if api.downloadSub != nil {
		api.downloadTopic.Close()
		api.downloadSub.Cancel()
	}
}

func (api *networkAPI) parsMessage(m *pubsub.Message) {
	// only forward messages delivered by others
	if m.ReceivedFrom == api.selfID {
		return
	}

	msg := api.firewall.ParsMessage(m.Data, m.ReceivedFrom)
	if msg != nil {
		api.parsFn(msg, m.ReceivedFrom)
	}
}

func (api *networkAPI) PublishMessage(msg *message.Message) error {
	topic := api.topic(msg)
	if topic == nil {
		return errors.Errorf(errors.ErrNetwork, "Invalid topic.")
	}
	if err := msg.SanityCheck(); err != nil {
		return err
	}
	data, err := msg.Encode()
	if err != nil {
		return err
	}

	return topic.Publish(api.ctx, data)
}

func (api *networkAPI) downloadLoop() {
	defer api.wg.Done()

	for {
		m, err := api.downloadSub.Next(api.ctx)
		if err != nil {
			logger.Debug("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) dataLoop() {
	defer api.wg.Done()

	for {
		m, err := api.dataSub.Next(api.ctx)
		if err != nil {
			logger.Debug("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) generalLoop() {
	defer api.wg.Done()

	for {
		m, err := api.generalSub.Next(api.ctx)
		if err != nil {
			logger.Debug("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) consensusLoop() {
	defer api.wg.Done()

	for {
		m, err := api.consensusSub.Next(api.ctx)
		if err != nil {
			logger.Debug("readLoop error", "err", err)
			return
		}

		api.parsMessage(m)
	}
}

func (api *networkAPI) topic(msg *message.Message) *pubsub.Topic {
	switch msg.PayloadType() {
	case payload.PayloadTypeSalam,
		payload.PayloadTypeAleyk,
		payload.PayloadTypeHeartBeat:
		return api.generalTopic

	case payload.PayloadTypeLatestBlocksRequest,
		payload.PayloadTypeLatestBlocksResponse,
		payload.PayloadTypeQueryTransactions,
		payload.PayloadTypeTransactions,
		payload.PayloadTypeBlockAnnounce:
		return api.dataTopic

	case payload.PayloadTypeQueryProposal,
		payload.PayloadTypeProposal,
		payload.PayloadTypeVote,
		payload.PayloadTypeQueryVotes:
		return api.consensusTopic

	case payload.PayloadTypeDownloadRequest,
		payload.PayloadTypeDownloadResponse:
		return api.downloadTopic

	default:
		panic("Invalid topic:")
	}
}

func (api *networkAPI) SelfID() peer.ID {
	return api.selfID
}
