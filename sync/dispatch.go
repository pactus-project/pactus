package sync

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
)

func (syncer *Synchronizer) sendBlocks(from, to int) {
	to = util.Min(to, from+syncer.config.BlockPerMessage)

	// Invalid range
	if from > to {
		return
	}

	// Help peer to catch up
	txs := make([]*tx.Tx, 0)
	blocks := make([]*block.Block, 0, to-from+1)
	for h := from; h <= to; h++ {
		b := syncer.cache.GetBlock(h)
		if b == nil {
			syncer.logger.Warn("Block can't find", "height", h)
			break
		}
		ids := b.TxIDs().IDs()
		for _, id := range ids {
			t := syncer.cache.GetTransaction(id)
			if t != nil {
				txs = append(txs, t)
			} else {
				syncer.logger.Warn("Transaction can't find", "id", id.Fingerprint())
			}
		}
		blocks = append(blocks, b)
	}

	var lastCommit *block.Commit
	ourHeight := syncer.state.LastBlockHeight()
	if to == ourHeight {
		lastCommit = syncer.state.LastCommit()
	}

	syncer.broadcastTxs(txs)
	syncer.broadcastBlocks(from, blocks, lastCommit)
}

func (syncer *Synchronizer) broadcastSalam() {
	msg := message.NewSalamMessage(
		syncer.state.GenesisHash(),
		syncer.state.LastBlockHeight())
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastAleyk() {
	msg := message.NewAleykMessage(
		syncer.state.GenesisHash(),
		syncer.state.LastBlockHeight())
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocksReq(from, to int, hash crypto.Hash) {
	msg := message.NewBlocksReqMessage(from, to, hash)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocks(from int, blocks []*block.Block, lastCommit *block.Commit) {
	msg := message.NewBlocksMessage(from, blocks, lastCommit)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastTxs(txs []*tx.Tx) {
	msg := message.NewTxsMessage(txs)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastHeartBeat() {
	hrs := syncer.consensus.HRS()

	// Probable we are syncing
	if !hrs.IsValid() {
		return
	}

	msg := message.NewHeartBeatMessage(syncer.state.LastBlockHash(), hrs)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastProposal(p *vote.Proposal) {
	msg := message.NewProposalMessage(p)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastVote(v *vote.Vote) {
	msg := message.NewVoteMessage(v)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastVoteSet(height int, hashes []crypto.Hash) {
	msg := message.NewVoteSetMessage(height, hashes)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) publishMessage(msg *message.Message) {
	if err := msg.SanityCheck(); err != nil {
		syncer.logger.Error("We have invalid message", "err", err, "message", msg)
		return
	}
	if err := syncer.networkAPI.PublishMessage(msg); err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg.Fingerprint(), "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg.Fingerprint())
	}
}
