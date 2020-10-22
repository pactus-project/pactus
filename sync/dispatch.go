package sync

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
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

func (syncer *Synchronizer) broadcastStatusReq() {
	height := syncer.state.LastBlockHeight()
	msg := message.NewStatusReqMessage(height)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocksRes(from int, blocks []block.Block, txs []tx.Tx) {
	msg := message.NewBlocksMessage(from, blocks, txs)
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

func (syncer *Synchronizer) broadcastHRS(hrs hrs.HRS, hasProposal bool) {
	msg := message.NewHRSMessage(hrs, hasProposal)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastProposal(proposal vote.Proposal, txs []tx.Tx) {
	msg := message.NewProposalMessage(proposal, txs)
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

	case message.PayloadTypeStatusReq,
		message.PayloadTypeStatusRes,
		message.PayloadTypeBlock,
		message.PayloadTypeBlocksReq,
		message.PayloadTypeBlocksRes:
		return syncer.stateTopic

	case message.PayloadTypeTxReq,
		message.PayloadTypeTxRes:
		return syncer.txTopic

	case message.PayloadTypeProposal,
		message.PayloadTypeHRS,
		message.PayloadTypeVote,
		message.PayloadTypeVoteSet:
		return syncer.consensusTopic
	}

	return nil
}
