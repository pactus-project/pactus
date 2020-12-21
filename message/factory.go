package message

import (
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/version"
	"github.com/zarbchain/zarb-go/vote"
)

func NewAleykMessage(genesisHash crypto.Hash, height int) *Message {
	return &Message{
		Type: payload.PayloadTypeAleyk,
		Payload: &payload.AleykPayload{
			Version:     version.NodeVersion,
			GenesisHash: genesisHash,
			Height:      height,
		},
	}

}

func NewBlocksReqMessage(from, to int, lastBlockHash crypto.Hash) *Message {
	return &Message{
		Type: payload.PayloadTypeBlocksReq,
		Payload: &payload.BlocksReqPayload{
			From:          from,
			To:            to,
			LastBlockHash: lastBlockHash,
		},
	}
}

func NewBlocksMessage(from int, blocks []*block.Block, lastCommit *block.Commit) *Message {
	return &Message{
		Type: payload.PayloadTypeBlocks,
		Payload: &payload.BlocksPayload{
			From:       from,
			Blocks:     blocks,
			LastCommit: lastCommit,
		},
	}
}

func NewHeartBeatMessage(lastBlockHash crypto.Hash, hrs hrs.HRS) *Message {
	return &Message{
		Type: payload.PayloadTypeHeartBeat,
		Payload: &payload.HeartBeatPayload{
			Pulse: hrs,
		},
	}
}

func NewProposalReqMessage(height, round int) *Message {
	return &Message{
		Type: payload.PayloadTypeProposalReq,
		Payload: &payload.ProposalReqPayload{
			Height: height,
			Round:  round,
		},
	}
}

func NewProposalMessage(proposal *vote.Proposal) *Message {
	return &Message{
		Type: payload.PayloadTypeProposal,
		Payload: &payload.ProposalPayload{
			Proposal: proposal,
		},
	}
}

func NewSalamMessage(genesisHash crypto.Hash, height int) *Message {
	return &Message{
		Type: payload.PayloadTypeSalam,
		Payload: &payload.SalamPayload{
			Version:     version.NodeVersion,
			GenesisHash: genesisHash,
			Height:      height,
		},
	}
}

func NewTxsReqMessage(ids []crypto.Hash) *Message {
	return &Message{
		Type: payload.PayloadTypeTxsReq,
		Payload: &payload.TxsReqPayload{
			IDs: ids,
		},
	}
}

func NewTxsMessage(txs []*tx.Tx) *Message {
	return &Message{
		Type: payload.PayloadTypeTxs,
		Payload: &payload.TxsPayload{
			Txs: txs,
		},
	}
}
func NewVoteSetMessage(height int, Hashes []crypto.Hash) *Message {
	return &Message{
		Type: payload.PayloadTypeVoteSet,
		Payload: &payload.VoteSetPayload{
			Height: height,
			Hashes: Hashes,
		},
	}
}

func NewVoteMessage(vote *vote.Vote) *Message {
	return &Message{
		Type: payload.PayloadTypeVote,
		Payload: &payload.VotePayload{
			Vote: vote,
		},
	}
}
