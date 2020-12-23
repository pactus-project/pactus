package sync

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/vote"
)

func (syncer *Synchronizer) sendBlocksReqIfWeAreBehind() {
	ourHeight := syncer.state.LastBlockHeight()
	claimedHeight := syncer.stats.LastClaimedHeight()
	if claimedHeight > ourHeight {
		syncer.logger.Debug("Ask for more blocks", "our_height", ourHeight, "claimed_height", claimedHeight)
		hash := syncer.state.LastBlockHash()
		syncer.broadcastBlocksReq(ourHeight+1, claimedHeight, hash)
	}
}

func (syncer *Synchronizer) sendBlocks(from, to int) {
	to = util.Min(to, from+syncer.config.BlockPerMessage)

	// Invalid range
	if from > to {
		return
	}

	// Help peer to catch up
	ids := make([]crypto.Hash, 0)
	blocks := make([]*block.Block, 0, to-from+1)
	for h := from; h <= to; h++ {
		b := syncer.cache.GetBlock(h)
		if b == nil {
			syncer.logger.Warn("Block can't find", "height", h)
			break
		}
		ids = append(ids, b.TxIDs().IDs()...)
		blocks = append(blocks, b)
	}

	var lastCommit *block.Commit
	ourHeight := syncer.state.LastBlockHeight()
	if to == ourHeight {
		lastCommit = syncer.state.LastCommit()
	}

	syncer.sendTransactions(ids)
	syncer.broadcastBlocks(from, blocks, lastCommit)
}

func (syncer *Synchronizer) sendTransactions(ids []crypto.Hash) {
	trxs := make([]*tx.Tx, 0, len(ids))
	for _, id := range ids {
		trx := syncer.cache.GetTransaction(id)
		if trx != nil {
			trxs = append(trxs, trx)
		} else {
			syncer.logger.Debug("Transaction can't find", "id", id.Fingerprint())
		}
	}

	if len(trxs) > 0 {
		syncer.broadcastTxs(trxs)
	}
}

func (syncer *Synchronizer) broadcastSalam() {
	msg := message.NewSalamMessage(
		syncer.config.Moniker,
		syncer.signer.PublicKey(),
		syncer.networkAPI.SelfID(),
		syncer.state.GenesisHash(),
		syncer.state.LastBlockHeight())
	syncer.publishMessage(msg)
}

func (syncer *Synchronizer) broadcastAleyk(resStatus int, resMsg string) {
	msg := message.NewAleykMessage(
		syncer.config.Moniker,
		syncer.signer.PublicKey(),
		syncer.networkAPI.SelfID(),
		syncer.state.GenesisHash(),
		syncer.state.LastBlockHeight(),
		resStatus,
		resMsg)
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
