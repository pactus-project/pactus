package message

import (
	peer "github.com/libp2p/go-libp2p-peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/message/payload"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/version"
	"github.com/zarbchain/zarb-go/vote"
)

func NewSalamMessage(moniker string, publicKey crypto.PublicKey, peerID peer.ID, genesisHash crypto.Hash, height int, flags int) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeSalam,
		Payload: &payload.SalamPayload{
			NodeVersion: version.NodeVersion,
			Moniker:     moniker,
			PublicKey:   publicKey,
			PeerID:      peerID,
			GenesisHash: genesisHash,
			Height:      height,
			Flags:       flags,
		},
	}
}
func NewAleykMessage(moniker string, publicKey crypto.PublicKey, peerID peer.ID, genesisHash crypto.Hash, height int, flags int, resStatus int, resMessage string) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeAleyk,
		Payload: &payload.AleykPayload{
			NodeVersion: version.NodeVersion,
			Moniker:     moniker,
			PublicKey:   publicKey,
			PeerID:      peerID,
			GenesisHash: genesisHash,
			Height:      height,
			Flags:       flags,
			Response: payload.SalamResponse{
				Status:  resStatus,
				Message: resMessage,
			},
		},
	}

}

func NewBlocksReqMessage(from, to int, lastBlockHash crypto.Hash) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeBlocksReq,
		Payload: &payload.BlocksReqPayload{
			From:          from,
			To:            to,
			LastBlockHash: lastBlockHash,
		},
	}
}

func NewBlocksMessage(from int, blocks []*block.Block, lastCommit *block.Commit) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeBlocks,
		Payload: &payload.BlocksPayload{
			From:       from,
			Blocks:     blocks,
			LastCommit: lastCommit,
		},
	}
}

func NewHeartBeatMessage(lastBlockHash crypto.Hash, hrs hrs.HRS) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeHeartBeat,
		Payload: &payload.HeartBeatPayload{
			Pulse: hrs,
		},
	}
}

func NewProposalReqMessage(height, round int) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeProposalReq,
		Payload: &payload.ProposalReqPayload{
			Height: height,
			Round:  round,
		},
	}
}

func NewProposalMessage(proposal *vote.Proposal) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeProposal,
		Payload: &payload.ProposalPayload{
			Proposal: proposal,
		},
	}
}

func NewTxsReqMessage(ids []crypto.Hash) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeTxsReq,
		Payload: &payload.TxsReqPayload{
			IDs: ids,
		},
	}
}

func NewTxsMessage(txs []*tx.Tx) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeTxs,
		Payload: &payload.TxsPayload{
			Txs: txs,
		},
	}
}
func NewVoteSetMessage(height, round int, Hashes []crypto.Hash) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeVoteSet,
		Payload: &payload.VoteSetPayload{
			Height: height,
			Round:  round,
			Hashes: Hashes,
		},
	}
}

func NewVoteMessage(vote *vote.Vote) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeVote,
		Payload: &payload.VotePayload{
			Vote: vote,
		},
	}
}

func makePayload(t payload.PayloadType) payload.Payload {
	switch t {
	case payload.PayloadTypeSalam:
		return &payload.SalamPayload{}
	case payload.PayloadTypeAleyk:
		return &payload.AleykPayload{}
	case payload.PayloadTypeBlocksReq:
		return &payload.BlocksReqPayload{}
	case payload.PayloadTypeBlocks:
		return &payload.BlocksPayload{}
	case payload.PayloadTypeTxsReq:
		return &payload.TxsReqPayload{}
	case payload.PayloadTypeTxs:
		return &payload.TxsPayload{}
	case payload.PayloadTypeProposalReq:
		return &payload.ProposalReqPayload{}
	case payload.PayloadTypeProposal:
		return &payload.ProposalPayload{}
	case payload.PayloadTypeHeartBeat:
		return &payload.HeartBeatPayload{}
	case payload.PayloadTypeVote:
		return &payload.VotePayload{}
	case payload.PayloadTypeVoteSet:
		return &payload.VoteSetPayload{}
	}

	//
	return nil
}
