package sync

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/zarbchain/zarb-go/message"
	"github.com/zarbchain/zarb-go/tx"
)

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

	//ourHeight, _ := syncer.state.LastBlockInfo()
	switch msg.PayloadType() {
	case message.PayloadTypeStatusReq:
		// pld := msg.Payload.(*message.StatusReqPayload)
		// switch h := pld.LastBlockHeight; {
		// case h == ourHeight:
		// case h == ourHeight+1:
		// case h == ourHeight-1:
		// 	{
		// 		// Do nothing
		// 		// Consensus lagging?
		// 	}

		// case h > ourHeight+1:
		// 	{

		// 	}

		// case h < ourHeight-1:
		// 	{
		// 		// Help peer to catch up
		// 		from := h
		// 		end := util.Min(h+10, ourHeight)
		// 		blocks := make([]block.Block, end-from)
		// 		for h := from; h <= end; h++ {
		// 			b, err := syncer.store.BlockByHeight(h)
		// 			if err != nil {
		// 				syncer.logger.Error("An error occurred while retriveng a block", "err", err, "height", h)
		// 				return
		// 			}
		// 			blocks[h-from] = *b
		// 		}

		// 		syncer.BroadcastBlocks(from, blocks)
		// 	}
		// }

	case message.PayloadTypeBlocksReq:
		// pld := msg.Payload.(*message.BlocksPayload)
		// height := pld.From
		// for _, b := range pld.Blocks {
		// 	if height > ourHeight {
		// 		bp, has := syncer.blockPool[height]
		// 		if has {
		// 			if !bp.Hash().EqualsTo(b.Hash()) {
		// 				syncer.logger.Error("We have received twoblock from same height but different hash", "from", m.ReceivedFrom.Pretty(), "height", height)
		// 			}
		// 		} else {
		// 			syncer.blockPool[height] = &b
		// 		}
		// 	}
		// 	height++
		// }
	case message.PayloadTypeTxRes:
		pld := msg.Payload.(*message.TxResPayload)
		syncer.txPool.AppendTx(&pld.Tx)

	case message.PayloadTypeTxReq:
		pld := msg.Payload.(*message.TxReqPayload)

		if syncer.txPool.HasTx(pld.Hash) {
			trx, _ := syncer.txPool.PendingTx(pld.Hash)
			msg := message.NewTxResMessage(*trx)
			syncer.publishMessage(msg)
		}
	case message.PayloadTypeHRS:
		pld := msg.Payload.(*message.HRSPayload)
		if !pld.HasProposal {
			p := syncer.consensus.LastProposal()
			if p != nil {
				hashes := p.Block().TxHashes().Hashes()
				txs := make([]tx.Tx, 0)
				for _, hash := range hashes {
					tx := syncer.txPool.GetTx(hash)
					if tx != nil {
						txs = append(txs, *tx)
					} else {
						syncer.logger.Warn("We don't have transation for a valid proposal", "hash", hash.Fingerprint())
					}
				}
				syncer.broadcastProposal(*p, txs)
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
				hasProposal := syncer.consensus.HasProposal()
				syncer.broadcastHRS(hrs, hasProposal)
			} else {
				// We are at the same step with this peer
			}
		}

	case message.PayloadTypeProposal:
		pld := msg.Payload.(*message.ProposalPayload)

		syncer.txPool.AppendTxs(pld.Txs)
		syncer.consensus.SetProposal(&pld.Proposal)
	case message.PayloadTypeBlock:
		//pld := msg.Payload.(*message.BlockPayload)

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
		syncer.logger.Error("Unknown message type", "msg", msg)
	}

}
