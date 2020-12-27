package sync

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/util"
)

func (syncer *Synchronizer) ParsMessage(data []byte, from peer.ID) {
	msg := syncer.stats.ParsMessage(data, from)
	if msg == nil {
		return
	}

	syncer.logger.Debug("Received a message", "from", util.FingerprintPeerID(from), "message", msg)

	switch msg.PayloadType() {
	case payload.PayloadTypeSalam:
		pld := msg.Payload.(*payload.SalamPayload)
		syncer.processSalamPayload(pld)

	case payload.PayloadTypeAleyk:
		pld := msg.Payload.(*payload.AleykPayload)
		syncer.processAleykPayload(pld)

	case payload.PayloadTypeHeartBeat:
		pld := msg.Payload.(*payload.HeartBeatPayload)
		syncer.processHeartBeatPayload(pld)

	case payload.PayloadTypeBlocksReq:
		pld := msg.Payload.(*payload.BlocksReqPayload)
		syncer.processBlocksReqPayload(pld)

	case payload.PayloadTypeBlocks:
		pld := msg.Payload.(*payload.BlocksPayload)
		syncer.processBlocksPayload(pld)

	case payload.PayloadTypeTxsReq:
		pld := msg.Payload.(*payload.TxsReqPayload)
		syncer.processTxsReqPayload(pld)

	case payload.PayloadTypeTxs:
		pld := msg.Payload.(*payload.TxsPayload)
		syncer.processTxsPayload(pld)

	case payload.PayloadTypeProposalReq:
		pld := msg.Payload.(*payload.ProposalReqPayload)
		syncer.processProposalReqPayload(pld)

	case payload.PayloadTypeProposal:
		pld := msg.Payload.(*payload.ProposalPayload)
		syncer.processProposalPayload(pld)

	case payload.PayloadTypeVote:
		pld := msg.Payload.(*payload.VotePayload)
		syncer.processVotePayload(pld)

	case payload.PayloadTypeVoteSet:
		pld := msg.Payload.(*payload.VoteSetPayload)
		syncer.processVoteSetPayload(pld)

	default:
		syncer.logger.Error("Unknown message type", "type", msg.PayloadType())
	}

	syncer.sendBlocksReqIfWeAreBehind()
}

func (syncer *Synchronizer) processSalamPayload(pld *payload.SalamPayload) {
	syncer.logger.Trace("Process salam payload", "pld", pld)

	if !pld.GenesisHash.EqualsTo(syncer.state.GenesisHash()) {
		syncer.logger.Info("Received a message from different chain", "genesis_hash", pld.GenesisHash)
		// Reply salam
		syncer.broadcastAleyk(payload.SalamResponseCodeRejected, "Invalid genesis hash")
		return
	}

	// Reply salam
	syncer.broadcastAleyk(payload.SalamResponseCodeOK, "Welcome!")
}

func (syncer *Synchronizer) processAleykPayload(pld *payload.AleykPayload) {
	syncer.logger.Trace("Process Aleyk payload", "pld", pld)
}

func (syncer *Synchronizer) processBlocksReqPayload(pld *payload.BlocksReqPayload) {
	syncer.logger.Trace("Process blocks request payload", "pld", pld)

	ourHeight := syncer.state.LastBlockHeight()
	if ourHeight < pld.From {
		return
	}

	b := syncer.cache.GetBlock(pld.From)
	if b != nil {
		if !b.Header().LastBlockHash().EqualsTo(pld.LastBlockHash) {
			syncer.logger.Info("We don't have previous block hash",
				"height", pld.From-1,
				"ourHash", b.Header().LastBlockHash(),
				"peerHash", pld.LastBlockHash)

			return
		}
	}

	syncer.sendBlocks(pld.From, pld.To)
}

func (syncer *Synchronizer) processBlocksPayload(pld *payload.BlocksPayload) {
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

	probableNewHeight := height - 1
	lastClaimedHeight := syncer.stats.LastClaimedHeight()

	if lastClaimedHeight > probableNewHeight {
		// Request for more blocks
		blockhash := pld.Blocks[len(pld.Blocks)-1].Hash()
		syncer.broadcastBlocksReq(probableNewHeight+1, lastClaimedHeight, blockhash)
	}

	syncer.tryCommitBlocks()

	// When a peer send us the last commit, we are probably synced with the network
	if pld.LastCommit != nil {
		syncer.informConsensusToMoveToNewHeight()
	}
}

func (syncer *Synchronizer) processTxsReqPayload(pld *payload.TxsReqPayload) {
	syncer.logger.Trace("Process txs request Payload", "pld", pld)

	syncer.sendTransactions(pld.IDs)
}

func (syncer *Synchronizer) processTxsPayload(pld *payload.TxsPayload) {
	syncer.logger.Trace("Process txs payload", "pld", pld)

	for _, trx := range pld.Txs {
		syncer.cache.AddTransaction(trx)
	}
}

func (syncer *Synchronizer) processVotePayload(pld *payload.VotePayload) {
	syncer.logger.Trace("Process vote payload", "pld", pld)

	syncer.consensus.AddVote(pld.Vote)
}

func (syncer *Synchronizer) processVoteSetPayload(pld *payload.VoteSetPayload) {
	syncer.logger.Trace("Process vote-set payload", "pld", pld)

	hrs := syncer.consensus.HRS()
	if pld.Height == hrs.Height() {

		// Check peers vote and send the votes he doesn't have
		ourVotes := syncer.consensus.RoundVotes(pld.Round)
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
func (syncer *Synchronizer) processProposalReqPayload(pld *payload.ProposalReqPayload) {
	syncer.logger.Trace("Process proposal request payload", "pld", pld)

	hrs := syncer.consensus.HRS()
	if pld.Height == hrs.Height() {
		p := syncer.consensus.LastProposal()
		if p != nil {
			if p.Round() >= pld.Round {
				syncer.broadcastProposal(p)
			}
		}
	}
}

func (syncer *Synchronizer) processProposalPayload(pld *payload.ProposalPayload) {
	syncer.logger.Trace("Process proposal payload", "pld", pld)

	syncer.consensus.SetProposal(pld.Proposal)
}

func (syncer *Synchronizer) processHeartBeatPayload(pld *payload.HeartBeatPayload) {
	syncer.logger.Trace("Process heartbeat payload", "pld", pld)

	hrs := syncer.consensus.HRS()
	if pld.Pulse.Height() == hrs.Height() {
		if pld.Pulse.GreaterThan(hrs) {
			syncer.logger.Trace("Our consensus is behind of this peer.")
			// Let's ask for more votes
			hashes := syncer.consensus.RoundVotesHash(hrs.Round())
			syncer.broadcastVoteSet(hrs.Height(), hrs.Round(), hashes)
		} else if pld.Pulse.LessThan(hrs) {
			syncer.logger.Trace("Our consensus is ahead of this peer.")
		} else {
			syncer.logger.Trace("Our consensus is at the same step with this peer.")
		}
	} else if pld.Pulse.Height() > hrs.Height() {
		syncer.logger.Trace("Our state is behind of this peer.")
	} else {
		syncer.logger.Trace("Our state is ahead of this peer.")
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
			break
		}
	}
}
