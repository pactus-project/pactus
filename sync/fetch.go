package sync

import (
	"encoding/hex"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func (syncer *Synchronizer) ParsMessage(data []byte, from peer.ID) {
	msg := new(message.Message)
	err := msg.UnmarshalCBOR(data)
	if err != nil {
		syncer.logger.Error("Error decoding message", "from", from.ShortString(), "message", msg, "err", err)
		return
	}
	syncer.logger.Trace("Received a message", "from", from.ShortString(), "message", msg)

	if err = msg.SanityCheck(); err != nil {
		syncer.stats.IncreaseInvalidMessageCounter(from)
		syncer.logger.Error("Peer sent us invalid msg", "from", from.ShortString(), "data", hex.EncodeToString(data), "err", err)
		return
	}

	syncer.stats.ProcessMessage(msg, from)

	switch msg.PayloadType() {
	case message.PayloadTypeSalam:
		pld := msg.Payload.(*message.SalamPayload)
		syncer.processSalamPayload(pld)

	case message.PayloadTypeHeartBeat:
		pld := msg.Payload.(*message.HeartBeatPayload)
		syncer.processHeartBeatPayload(pld)

	case message.PayloadTypeBlocksReq:
		pld := msg.Payload.(*message.BlocksReqPayload)
		syncer.processBlocksReqPayload(pld)

	case message.PayloadTypeBlocks:
		pld := msg.Payload.(*message.BlocksPayload)
		syncer.processBlocksResPayload(pld)

	case message.PayloadTypeTxsReq:
		pld := msg.Payload.(*message.TxsReqPayload)
		syncer.processTxsReqPayload(pld)

	case message.PayloadTypeTxs:
		pld := msg.Payload.(*message.TxsPayload)
		syncer.processTxsPayload(pld)

	case message.PayloadTypeProposal:
		pld := msg.Payload.(*message.ProposalPayload)
		syncer.consensus.SetProposal(&pld.Proposal)

	case message.PayloadTypeVote:
		pld := msg.Payload.(*message.VotePayload)
		syncer.processVotePayload(pld)

	case message.PayloadTypeVoteSet:
		pld := msg.Payload.(*message.VoteSetPayload)
		syncer.processVoteSetPayload(pld)

	default:
		syncer.logger.Error("Unknown message type", "type", msg.PayloadType())
	}
}

func (syncer *Synchronizer) processSalamPayload(pld *message.SalamPayload) {
	ourHeight := syncer.state.LastBlockHeight()

	switch h := pld.Height; {
	case h > ourHeight:
		{
			syncer.broadcastBlocksReq(ourHeight+1, pld.Height)
		}
	case h < ourHeight:
		{
			// Reply salam
			syncer.broadcastSalam()
			syncer.sendBlocks(h+1, ourHeight)
		}
	}
}

func (syncer *Synchronizer) processBlocksReqPayload(pld *message.BlocksReqPayload) {
	b, err := syncer.store.BlockByHeight(pld.From)
	if err == nil {
		if b.Header().LastBlockHash().EqualsTo(pld.LastBlockHash) {
			syncer.sendBlocks(pld.From, pld.To)
		}
	}
}

func (syncer *Synchronizer) processBlocksResPayload(pld *message.BlocksPayload) {
	height := pld.From
	ourHeight := syncer.state.LastBlockHeight()

	for _, block := range pld.Blocks {
		if height < ourHeight {
			continue // We already have committed this block
		}

		if height > ourHeight+1 {
			syncer.blockPool.AppendCommit(
				block.Header().LastBlockHash(),
				block.LastCommit())
		}

		syncer.blockPool.AppendBlock(height, block)

		commitBlock := syncer.blockPool.Block(ourHeight + 1)
		commit := syncer.blockPool.Commit(commitBlock.Hash())

		if commitBlock != nil && commit != nil {
			syncer.logger.Info("Committing block", "height", ourHeight+1, "block", commitBlock)
			if err := syncer.state.ApplyBlock(*commitBlock, *commit); err != nil {
				syncer.logger.Error("Committing block failed", "block", commitBlock, "err", err)
				// Ask peers to send us this block again
				syncer.broadcastBlocksReq(ourHeight+1, ourHeight+2)
			} else {
				syncer.consensus.ScheduleNewHeight()
			}

			syncer.blockPool.RemoveBlock(ourHeight + 1)
			syncer.blockPool.RemoveCommit(commitBlock.Hash())
		}
		height = height + 1
	}
}

func (syncer *Synchronizer) processTxsReqPayload(pld *message.TxsReqPayload) {
	txs := make([]tx.Tx, 0, len(pld.Hashes))
	for _, h := range pld.Hashes {
		trx := syncer.txPool.PendingTx(h)
		if trx == nil {
			trx, _, _ = syncer.store.Tx(h)
		}
		txs = append(txs, *trx)
	}

	syncer.broadcastTxs(txs)
}

func (syncer *Synchronizer) processTxsPayload(pld *message.TxsPayload) {
	syncer.txPool.AppendTxs(pld.Txs)
}

func (syncer *Synchronizer) processVotePayload(pld *message.VotePayload) {
	syncer.consensus.AddVote(pld.Vote)
}

func (syncer *Synchronizer) processVoteSetPayload(pld *message.VoteSetPayload) {
	hrs := syncer.consensus.HRS()
	if pld.Height == hrs.Height() {
		// Sending votes to peer
		ourVotes := syncer.consensus.AllVotes()
		peerVotes := pld.Hashes

		for _, v1 := range ourVotes {
			hasVote := false
			for _, v2 := range peerVotes {
				if v1.Hash() == v2 {
					hasVote = true
					break
				}
			}

			if !hasVote {
				syncer.broadcastVote(v1)
			}
		}
	}
}

func (syncer *Synchronizer) processHeartBeatPayload(pld *message.HeartBeatPayload) {
	if !pld.HasProposal {
		p := syncer.consensus.LastProposal()
		if p != nil {
			syncer.broadcastProposal(*p)
		}
	}
	hrs := syncer.consensus.HRS()
	if pld.HRS.Height() == hrs.Height() {
		if pld.HRS.GreaterThan(hrs) {
			// We are behind of the peer.
			// Let's ask for more votes
			hashes := syncer.consensus.AllVotesHashes()
			syncer.broadcastVoteSet(hrs.Height(), hashes)
		} else if pld.HRS.LessThan(hrs) {
			// We are ahead of the peer.
			// Let's inform him know about our status
			syncer.broadcastHeartBeat()
		} else {
			// We are at the same step with this peer
		}
	}
}

func (syncer *Synchronizer) sendBlocks(from, to int) {
	ourHeight := syncer.state.LastBlockHeight()

	to = util.Min(to, ourHeight)
	to = util.Min(to, to+syncer.config.BlockPerMessage)

	// Help peer to catch up
	txs := make([]tx.Tx, 0)
	blocks := make([]block.Block, to-from+1)
	for h := from; h <= to; h++ {
		b, err := syncer.store.BlockByHeight(h)
		if err != nil {
			syncer.logger.Error("An error occurred while retriveng a block", "err", err, "height", h)
			break
		}
		hashes := b.TxHashes().Hashes()
		for _, hash := range hashes {
			tx, _, _ := syncer.store.Tx(hash)
			if tx != nil {
				txs = append(txs, *tx)
			} else {
				syncer.logger.Warn("We don't have transation for the block", "hash", hash.Fingerprint())
			}
		}
		blocks[h-from] = *b
	}

	syncer.broadcastTxs(txs)
	syncer.broadcastBlocks(from, blocks, nil)
}
