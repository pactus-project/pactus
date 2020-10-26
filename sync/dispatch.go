package sync

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/vote"
)

func (syncer *Synchronizer) broadcastLoop() {
	for {
		select {
		case <-syncer.ctx.Done():
			return

		case msg := <-syncer.broadcastCh:
			syncer.publishMessage(msg)
		}
	}
}

func (syncer *Synchronizer) broadcastSalam() {
	height := syncer.state.LastBlockHeight()
	msg := message.NewSalamMessage(height)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocksReq(from, to int) {
	msg := message.NewBlocksReqMessage(from, to)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocksRes(from int, blocks []block.Block, txs []tx.Tx) {
	msg := message.NewBlocksResMessage(from, blocks, txs)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastTxRes(trx tx.Tx) {
	msg := message.NewTxResMessage(trx)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastTxReq(hash crypto.Hash) {
	msg := message.NewTxReqMessage(hash)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastHeartBeat() {
	hasProposal := syncer.consensus.HasProposal()
	hrs := syncer.consensus.HRS()

	msg := message.NewHeartBeatMessage(hrs, hasProposal)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastProposal(proposal vote.Proposal) {
	msg := message.NewProposalMessage(proposal)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) publishMessage(msg message.Message) {
	msg.Initiator = syncer.selfAddress

	topic := syncer.topic(&msg)
	if topic != nil {
		bs, _ := msg.MarshalCBOR()
		if err := topic.Publish(syncer.ctx, bs); err != nil {
			syncer.logger.Error("Error on publishing message", "message", msg.Fingerprint(), "err", err)
		} else {
			syncer.logger.Trace("Publishing new message", "message", msg.Fingerprint())
		}
	}
}

func (syncer *Synchronizer) topic(msg *message.Message) *pubsub.Topic {
	switch msg.PayloadType() {

	case message.PayloadTypeSalam,
		message.PayloadTypeHeartBeat:
		return syncer.generalTopic

	case message.PayloadTypeBlock,
		message.PayloadTypeBlocksReq,
		message.PayloadTypeBlocksRes:
		return syncer.blockTopic

	case message.PayloadTypeTxReq,
		message.PayloadTypeTxRes:
		return syncer.txTopic

	case message.PayloadTypeProposal,
		message.PayloadTypeVote,
		message.PayloadTypeVoteSet:
		return syncer.consensusTopic
	default:
		panic("Invalid topic")
	}
}
