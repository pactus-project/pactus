package sync

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

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
			return
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

	syncer.broadcastBlocksRes(from, blocks, txs)
}

func (syncer *Synchronizer) appendBlock(height int, block block.Block) {
	ourHeight := syncer.state.LastBlockHeight()
	if height < ourHeight {
		// We already have committed this block
		return
	}

	if err := syncer.blockPool.AppendBlock(height, block); err != nil {
		// Invalid block? ask to send this block again
		syncer.broadcastBlocksReq(height, height+1)
		return
	}

	b1 := syncer.blockPool.Block(ourHeight + 1)
	b2 := syncer.blockPool.Block(ourHeight + 2)
	if b1 != nil && b2 != nil {
		syncer.logger.Info("Committing block", "height", ourHeight+1, "block", b1)
		if err := syncer.state.ApplyBlock(*b1, *b2.LastCommit()); err != nil {
			syncer.logger.Error("Committing block failed", "block", b1, "err", err)
			syncer.blockPool.RemoveBlock(ourHeight + 1)
			syncer.blockPool.RemoveBlock(ourHeight + 2)
		} else {
			syncer.blockPool.RemoveBlock(ourHeight + 1)

			syncer.consensus.ScheduleNewHeight()
		}
	}
}

func (syncer *Synchronizer) parsMessage(m *pubsub.Message) {

	// only forward messages delivered by others
	if m.ReceivedFrom == syncer.selfID {
		return
	}

	msg := new(message.Message)
	err := msg.UnmarshalCBOR(m.Data)
	if err != nil {
		syncer.logger.Error("Error decoding message", "from", m.ReceivedFrom.Pretty(), "message", msg, "err", err)
		return
	}
	syncer.logger.Trace("Received a message", "from", m.ReceivedFrom.Pretty(), "message", msg)

	if err = msg.SanityCheck(); err != nil {
		syncer.logger.Error("Peer sent us invalid msg", "from", m.ReceivedFrom.Pretty(), "message", msg, "err", err)
		return
	}

	syncer.stats.ParsPeerMessage(m.ReceivedFrom, msg)

	ourHeight := syncer.state.LastBlockHeight()
	switch msg.PayloadType() {
	case message.PayloadTypeSalam:
		pld := msg.Payload.(*message.SalamPayload)
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

	case message.PayloadTypeHeartBeat:
		pld := msg.Payload.(*message.HeartBeatPayload)
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
				msg := message.NewVoteSetMessage(hrs.Height(), hashes)
				syncer.publishMessage(msg)
			} else if pld.HRS.LessThan(hrs) {
				// We are ahead of the peer.
				// Let's inform him know about our status
				syncer.broadcastHeartBeat()
			} else {
				// We are at the same step with this peer
			}
		}

	case message.PayloadTypeBlocksReq:
		pld := msg.Payload.(*message.BlocksReqPayload)
		syncer.sendBlocks(pld.From, pld.To)

	case message.PayloadTypeBlocksRes:
		pld := msg.Payload.(*message.BlocksResPayload)
		syncer.txPool.AppendTxs(pld.Txs)

		height := pld.From
		for _, b := range pld.Blocks {
			syncer.appendBlock(height, b)
			height = height + 1
		}

	case message.PayloadTypeBlock:

	case message.PayloadTypeTxRes:
		pld := msg.Payload.(*message.TxResPayload)
		syncer.txPool.AppendTx(pld.Tx)

	case message.PayloadTypeTxReq:
		pld := msg.Payload.(*message.TxReqPayload)

		if syncer.txPool.HasTx(pld.Hash) {
			trx := syncer.findTx(pld.Hash)
			msg := message.NewTxResMessage(*trx)
			syncer.publishMessage(msg)
		}

	case message.PayloadTypeProposal:
		pld := msg.Payload.(*message.ProposalPayload)
		syncer.consensus.SetProposal(&pld.Proposal)

	case message.PayloadTypeVote:
		pld := msg.Payload.(*message.VotePayload)
		syncer.consensus.AddVote(pld.Vote)

	case message.PayloadTypeVoteSet:
		pld := msg.Payload.(*message.VoteSetPayload)
		hrs := syncer.consensus.HRS()
		if pld.Height == hrs.Height() {
			// Sending votes to peer
			ourVotes := syncer.consensus.AllVotes()
			peerVotes := pld.Votes

			for _, v1 := range ourVotes {
				hasVote := false
				for _, v2 := range peerVotes {
					if v1.Hash() == v2 {
						hasVote = true
						break
					}
				}

				if !hasVote {
					msg := message.NewVoteMessage(v1)
					syncer.publishMessage(msg)
				}
			}
		}

	default:
		syncer.logger.Error("Unknown message type", "type", msg.PayloadType())
	}
}

func (syncer *Synchronizer) findTx(hash crypto.Hash) *tx.Tx {
	trx := syncer.txPool.PendingTx(hash)
	if trx != nil {
		return trx
	}
	trx, _, _ = syncer.store.Tx(hash)
	if trx != nil {
		return trx
	}
	return nil
}
