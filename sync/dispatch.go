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
	txs := make([]tx.Tx, 0)
	blocks := make([]block.Block, to-from+1)
	for h := from; h <= to; h++ {
		b, err := syncer.store.Block(h)
		if err != nil {
			syncer.logger.Error("An error occurred while retriveng a block", "err", err, "height", h)
			break
		}
		hashes := b.TxHashes().Hashes()
		for _, hash := range hashes {
			ctrx, _ := syncer.store.Transaction(hash)
			if ctrx != nil {
				txs = append(txs, *ctrx.Tx)
			} else {
				syncer.logger.Warn("We don't have transaction for the block", "hash", hash.Fingerprint())
			}
		}
		blocks[h-from] = *b
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

func (syncer *Synchronizer) broadcastBlocksReq(to int) {
	from := syncer.state.LastBlockHeight() + 1
	hash := syncer.state.LastBlockHash()
	msg := message.NewBlocksReqMessage(from, to, hash)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastBlocks(from int, blocks []block.Block, lastCommit *block.Commit) {
	if len(blocks) == 0 {
		return
	}
	msg := message.NewBlocksMessage(from, blocks, lastCommit)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastTxs(txs []tx.Tx) {
	if len(txs) == 0 {
		return
	}
	msg := message.NewTxsMessage(txs)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastHeartBeat() {
	hrs := syncer.consensus.HRS()

	msg := message.NewHeartBeatMessage(syncer.state.LastBlockHash(), hrs)
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

func (syncer *Synchronizer) broadcastVoteSet(height int, hashes []crypto.Hash) {
	msg := message.NewVoteSetMessage(height, hashes)
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) publishMessage(msg *message.Message) {

	if err := syncer.networkAPI.PublishMessage(msg); err != nil {
		syncer.logger.Error("Error on publishing message", "message", msg.Fingerprint(), "err", err)
	} else {
		syncer.logger.Debug("Publishing new message", "message", msg.Fingerprint())
	}
}
