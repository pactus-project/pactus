package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
)

func (syncer *Synchronizer) ParsMessage(data []byte, from peer.ID) {
	msg := syncer.stats.ParsMessage(data, from)
	if msg == nil {
		return
	}

	syncer.logger.Trace("Received a message", "from", from.ShortString(), "message", msg)

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
		syncer.processBlocksPayload(pld)

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

func (syncer *Synchronizer) processBlocksPayload(pld *message.BlocksPayload) {
	height := pld.From
	ourHeight := syncer.state.LastBlockHeight()

	if pld.LastCommit != nil {
		if height+len(pld.Blocks) > ourHeight {
			syncer.blockPool.AppendCommit(
				pld.Blocks[len(pld.Blocks)-1].Header().Hash(),
				pld.LastCommit)
		}
	}

	for _, block := range pld.Blocks {

		// For preventing any race condition situation
		ourHeight := syncer.state.LastBlockHeight()

		if height < ourHeight {
			continue // We already have committed this block
		}

		if height > ourHeight+1 {
			syncer.blockPool.AppendCommit(
				block.Header().LastBlockHash(),
				block.LastCommit())
		}

		syncer.blockPool.AppendBlock(height, block)
		syncer.tryCommitBlocks()

		height = height + 1
	}

	newHeight := syncer.state.LastBlockHeight()
	if height-2 > newHeight {
		// Ask peers to send us the missed block
		syncer.broadcastBlocksReq(newHeight+1, newHeight+2)
	}
}

func (syncer *Synchronizer) tryCommitBlocks() {
	for {
		ourHeight := syncer.state.LastBlockHeight()
		commitBlock := syncer.blockPool.Block(ourHeight + 1)

		if commitBlock == nil {
			break
		}
		commit := syncer.blockPool.Commit(commitBlock.Hash())
		if commit == nil {
			break
		}
		syncer.logger.Info("Committing block", "height", ourHeight+1, "block", commitBlock)
		if err := syncer.state.ApplyBlock(ourHeight+1, *commitBlock, *commit); err != nil {
			syncer.logger.Error("Committing block failed", "block", commitBlock, "err", err)
			// We will ask peers to send this block later ...
		}

		syncer.consensus.ScheduleNewHeight()
		syncer.blockPool.RemoveBlock(ourHeight + 1)
		syncer.blockPool.RemoveCommit(commitBlock.Hash())
	}
}

func (syncer *Synchronizer) processTxsReqPayload(pld *message.TxsReqPayload) {
	txs := make([]tx.Tx, 0, len(pld.Hashes))
	for _, h := range pld.Hashes {
		trx := syncer.txPool.PendingTx(h)
		if trx == nil {
			trx, _, _ = syncer.store.Tx(h)
		}
		if trx != nil {
			txs = append(txs, *trx)
		}
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
	hrs := syncer.consensus.HRS()
	if pld.HRS.Height() == hrs.Height() {
		if !pld.HasProposal {
			p := syncer.consensus.LastProposal()
			if p != nil {
				syncer.broadcastProposal(*p)
			}
		}

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
	} else if pld.HRS.Height() > hrs.Height() {
		// Ask for more blocks from this peer
		syncer.broadcastBlocksReq(hrs.Height()+1, pld.HRS.Height())
	} else {
		// We are ahead of this peer
	}
}
