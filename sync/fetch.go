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

	case message.PayloadTypeAleyk:
		pld := msg.Payload.(*message.AleykPayload)
		syncer.processAleykPayload(pld)

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

	case message.PayloadTypeProposalReq:
		pld := msg.Payload.(*message.ProposalReqPayload)
		syncer.processProposalReqPayload(pld)

	case message.PayloadTypeProposal:
		pld := msg.Payload.(*message.ProposalPayload)
		syncer.processProposalPayload(pld)

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
	syncer.logger.Trace("Process salam payload", "pld", pld)

	if !pld.GenesisHash.EqualsTo(syncer.state.GenesisHash()) {
		syncer.logger.Info("Received a message from different chain", "Genesis hash", pld.GenesisHash)
		return
	}

	// Reply salam
	syncer.broadcastAleyk()

	syncer.sendBlocksReqIfWeAreBehind(pld.Height)
}

func (syncer *Synchronizer) processAleykPayload(pld *message.AleykPayload) {
	syncer.logger.Trace("Process Aleyk payload", "pld", pld)

	syncer.sendBlocksReqIfWeAreBehind(pld.Height)
}

func (syncer *Synchronizer) processBlocksReqPayload(pld *message.BlocksReqPayload) {
	syncer.logger.Trace("Process blocks request payload", "pld", pld)

	ourHeight := syncer.state.LastBlockHeight()
	if ourHeight < pld.From {
		return
	}

	b := syncer.cache.GetBlock(pld.From)
	if b != nil {
		if !b.Header().LastBlockHash().EqualsTo(pld.LastBlockHash) {
			syncer.logger.Info("a peer has a block which we have no trace of it",
				"height", pld.From-1,
				"ourHash", b.Header().LastBlockHash(),
				"peerHash", pld.LastBlockHash)

			return
		}
	}

	syncer.sendBlocks(pld.From, pld.To)
}

func (syncer *Synchronizer) processBlocksPayload(pld *message.BlocksPayload) {
	syncer.logger.Trace("Process blocks payload", "pld", pld)

	ourHeight := syncer.state.LastBlockHeight()
	if ourHeight >= pld.To() {
		return
	}

	if pld.LastCommit != nil {
		syncer.cache.AddCommit(
			pld.Blocks[len(pld.Blocks)-1].Header().Hash(),
			pld.LastCommit)
	}

	height := pld.From
	for _, block := range pld.Blocks {
		syncer.cache.AddCommit(
			block.Header().LastBlockHash(),
			block.LastCommit())

		syncer.cache.AddBlock(height, block)

		height = height + 1
	}

	probableNewHeight := height - 2
	networkMaxHeight := syncer.stats.MaxHeight()

	if networkMaxHeight > probableNewHeight {
		// Request for more blocks
		blockhash := pld.Blocks[len(pld.Blocks)-2].Hash()
		syncer.broadcastBlocksReq(probableNewHeight, networkMaxHeight, blockhash)
	}

	syncer.tryCommitBlocks()

	// Check if our high is not same as we expected
	ourNewHeight := syncer.state.LastBlockHeight()

	if ourNewHeight != probableNewHeight {
		blockhash := syncer.state.LastBlockHash()
		syncer.broadcastBlocksReq(ourNewHeight, networkMaxHeight, blockhash)
	}

	// When a peer send us the last commit, it probably is be in latest hight
	syncer.maybeSynced(pld.LastCommit != nil)
}

func (syncer *Synchronizer) processTxsReqPayload(pld *message.TxsReqPayload) {
	syncer.logger.Trace("Process txs request Payload", "pld", pld)

	txs := make([]*tx.Tx, 0, len(pld.IDs))
	for _, h := range pld.IDs {
		trx := syncer.cache.GetTransaction(h)
		if trx != nil {
			txs = append(txs, trx)
		}
	}

	syncer.broadcastTxs(txs)
}

func (syncer *Synchronizer) processTxsPayload(pld *message.TxsPayload) {
	syncer.logger.Trace("Process txs payload", "pld", pld)

	for _, trx := range pld.Txs {
		syncer.cache.AddTransaction(trx)
	}
}

func (syncer *Synchronizer) processVotePayload(pld *message.VotePayload) {
	syncer.logger.Trace("Process vote payload", "pld", pld)

	syncer.consensus.AddVote(pld.Vote)
}

func (syncer *Synchronizer) processVoteSetPayload(pld *message.VoteSetPayload) {
	syncer.logger.Trace("Process vote-set payload", "pld", pld)

	hrs := syncer.consensus.HRS()
	if pld.Height == hrs.Height() {

		// Check peers vote and send the votes he doesn't have
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
func (syncer *Synchronizer) processProposalReqPayload(pld *message.ProposalReqPayload) {
	syncer.logger.Trace("Process proposal request payload", "pld", pld)

	hrs := syncer.consensus.HRS()
	if pld.Height == hrs.Height() {
		p := syncer.consensus.LastProposal()
		if p != nil {
			syncer.broadcastProposal(p)
		}
	}
}

func (syncer *Synchronizer) processProposalPayload(pld *message.ProposalPayload) {
	syncer.logger.Trace("Process proposal payload", "pld", pld)

	syncer.consensus.SetProposal(pld.Proposal)
}

func (syncer *Synchronizer) processHeartBeatPayload(pld *message.HeartBeatPayload) {
	syncer.logger.Trace("Process heartbeat payload", "pld", pld)

	hrs := syncer.consensus.HRS()
	if pld.Pulse.Height() == hrs.Height() {
		if pld.Pulse.GreaterThan(hrs) {
			syncer.logger.Trace("Our consensus is behind of this peer.")
			// Let's ask for more votes
			hashes := syncer.consensus.AllVotesHashes()
			syncer.broadcastVoteSet(hrs.Height(), hashes)
		} else if pld.Pulse.LessThan(hrs) {
			syncer.logger.Trace("Our consensus is ahead of this peer.")
		} else {
			syncer.logger.Trace("Our consensus is at the same step with this peer.")
		}
	} else if pld.Pulse.Height() > hrs.Height() {
		// Ask for more blocks from this peer
		syncer.sendBlocksReqIfWeAreBehind(pld.Pulse.Height())
	} else {
		syncer.logger.Trace("We are ahead of this peer.")
	}
}

func (syncer *Synchronizer) tryCommitBlocks() {
	for {
		ourHeight := syncer.state.LastBlockHeight()
		b := syncer.cache.GetBlock(ourHeight + 1)
		if b == nil {
			break
		}
		c := syncer.cache.GetCommit(b.Hash())
		if c == nil {
			break
		}
		syncer.logger.Trace("Committing block", "height", ourHeight+1, "block", b)
		if err := syncer.state.ApplyBlock(ourHeight+1, *b, *c); err != nil {
			syncer.logger.Error("Committing block failed", "block", b, "err", err, "height", ourHeight+1)
			// We will ask peers to send this block later ...

			// TODO: add tests for me later
			// TODO: Remove this invalid block and commit from the cache
		}
	}
}

func (syncer *Synchronizer) sendBlocksReqIfWeAreBehind(peerHeight int) {

	ourHeight := syncer.state.LastBlockHeight()
	if peerHeight > ourHeight {
		hash := syncer.state.LastBlockHash()

		syncer.broadcastBlocksReq(ourHeight+1, peerHeight, hash)
	}
}
