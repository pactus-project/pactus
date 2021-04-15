package network

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sync/message"
	"github.com/zarbchain/zarb-go/sync/message/payload"
)

func (n *network) PublishMessage(msg *message.Message) error {
	topic := n.topic(msg)
	if topic == nil {
		return errors.Errorf(errors.ErrNetwork, "invalid topic.")
	}
	if err := msg.SanityCheck(); err != nil {
		return err
	}
	data, err := msg.Encode()
	if err != nil {
		return err
	}

	return topic.Publish(n.ctx, data)
}

func (n *network) JoinTopics(callbackFn CallbackFn) error {
	generalTopic, err := n.joinTopic("general")
	if err != nil {
		return err
	}
	generalSub, err := generalTopic.Subscribe()
	if err != nil {
		return err
	}
	dataTopic, err := n.joinTopic("data")
	if err != nil {
		return err
	}
	dataSub, err := dataTopic.Subscribe()
	if err != nil {
		return err
	}
	consensusTopic, err := n.joinTopic("consensus")
	if err != nil {
		return err
	}
	consensusSub, err := consensusTopic.Subscribe()
	if err != nil {
		return err
	}

	n.callback = callbackFn
	n.generalTopic = generalTopic
	n.dataTopic = dataTopic
	n.consensusTopic = consensusTopic
	n.generalSub = generalSub
	n.dataSub = dataSub
	n.consensusSub = consensusSub

	n.wg.Add(1)
	go n.dataLoop()

	n.wg.Add(1)
	go n.generalLoop()

	n.wg.Add(1)
	go n.consensusLoop()

	return nil
}

func (n *network) JoinDownloadTopic() error {
	n.lk.Lock()
	defer n.lk.Unlock()

	if n.downloadSub != nil {
		return nil
	}

	downloadTopic, err := n.joinTopic("download")
	if err != nil {
		return err
	}
	downloadSub, err := downloadTopic.Subscribe()
	if err != nil {
		return err
	}
	n.downloadTopic = downloadTopic
	n.downloadSub = downloadSub

	n.wg.Add(1)
	go n.downloadLoop()

	return nil
}

func (n *network) LeaveDownloadTopic() {
	n.lk.Lock()
	defer n.lk.Unlock()

	if n.downloadSub != nil {
		n.downloadTopic.Close()
		n.downloadSub.Cancel()
	}
}

func (n *network) downloadLoop() {
	defer n.wg.Done()

	for {
		m, err := n.downloadSub.Next(n.ctx)
		if err != nil {
			n.logger.Debug("readLoop error", "err", err)
			return
		}

		n.onReceiveMessage(m)
	}
}

func (n *network) dataLoop() {
	defer n.wg.Done()

	for {
		m, err := n.dataSub.Next(n.ctx)
		if err != nil {
			n.logger.Debug("readLoop error", "err", err)
			return
		}

		n.onReceiveMessage(m)
	}
}

func (n *network) generalLoop() {
	defer n.wg.Done()

	for {
		m, err := n.generalSub.Next(n.ctx)
		if err != nil {
			n.logger.Debug("readLoop error", "err", err)
			return
		}

		n.onReceiveMessage(m)
	}
}

func (n *network) consensusLoop() {
	defer n.wg.Done()

	for {
		m, err := n.consensusSub.Next(n.ctx)
		if err != nil {
			n.logger.Debug("readLoop error", "err", err)
			return
		}

		n.onReceiveMessage(m)
	}
}

func (n *network) closeTopics() {
	n.LeaveDownloadTopic()

	n.dataTopic.Close()
	n.dataSub.Cancel()

	n.generalTopic.Close()
	n.generalSub.Cancel()

	n.consensusTopic.Close()
	n.consensusSub.Cancel()

	n.wg.Wait()
}

func (n *network) topic(msg *message.Message) *pubsub.Topic {
	switch msg.Payload.Type() {
	case payload.PayloadTypeSalam,
		payload.PayloadTypeAleyk,
		payload.PayloadTypeHeartBeat:
		return n.generalTopic

	case payload.PayloadTypeLatestBlocksRequest,
		payload.PayloadTypeLatestBlocksResponse,
		payload.PayloadTypeQueryTransactions,
		payload.PayloadTypeTransactions,
		payload.PayloadTypeBlockAnnounce:
		return n.dataTopic

	case payload.PayloadTypeQueryProposal,
		payload.PayloadTypeProposal,
		payload.PayloadTypeVote,
		payload.PayloadTypeQueryVotes:
		return n.consensusTopic

	case payload.PayloadTypeDownloadRequest,
		payload.PayloadTypeDownloadResponse:
		return n.downloadTopic

	default:
		panic("Invalid topic:")
	}
}

func (n *network) onReceiveMessage(m *pubsub.Message) {
	// only forward messages delivered by others
	if m.ReceivedFrom == n.SelfID() {
		return
	}

	n.callback(m.Data, m.ReceivedFrom)
}
