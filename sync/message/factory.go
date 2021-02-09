package message

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/consensus/hrs"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/message/payload"
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
func NewAleykMessage(code payload.ResponseCode, message string, moniker string, publicKey crypto.PublicKey, peerID peer.ID, height int, flags int) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeAleyk,
		Payload: &payload.AleykPayload{
			ResponseCode:    code,
			ResponseMessage: message,
			NodeVersion:     version.NodeVersion,
			Moniker:         moniker,
			PublicKey:       publicKey,
			PeerID:          peerID,
			Height:          height,
			Flags:           flags,
		},
	}

}

func NewLatestBlocksRequestMessage(initiator, target peer.ID, sessionID int, from int) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeLatestBlocksRequest,
		Payload: &payload.LatestBlocksRequestPayload{
			Initiator: initiator,
			Target:    target,
			SessionID: sessionID,
			From:      from,
		},
	}
}

func NewLatestBlocksResponseMessage(code payload.ResponseCode, initiator, target peer.ID, sessionID int, from int, blocks []*block.Block, transactions []*tx.Tx, commit *block.Commit) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeLatestBlocksResponse,
		Payload: &payload.LatestBlocksResponsePayload{
			ResponseCode: code,
			SessionID:    sessionID,
			Initiator:    initiator,
			Target:       target,
			From:         from,
			Blocks:       blocks,
			Transactions: transactions,
			LastCommit:   commit,
		},
	}
}

func NewHeartBeatMessage(id peer.ID, lastBlockHash crypto.Hash, hrs hrs.HRS) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeHeartBeat,
		Payload: &payload.HeartBeatPayload{
			PeerID:        id,
			Pulse:         hrs,
			LastBlockHash: lastBlockHash,
		},
	}
}

func NewOpaqueQueryProposalMessage(height, round int) *Message {
	return NewQueryProposalMessage("", height, round)
}

func NewQueryProposalMessage(querier peer.ID, height, round int) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeQueryProposal,
		Payload: &payload.QueryProposalPayload{
			Querier: querier,
			Height:  height,
			Round:   round,
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

func NewOpaqueQueryTransactionsMessage(ids []crypto.Hash) *Message {
	return NewQueryTransactionsMessage("", ids)
}

func NewQueryTransactionsMessage(querier peer.ID, ids []crypto.Hash) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeQueryTransactions,
		Payload: &payload.QueryTransactionsPayload{
			Querier: querier,
			IDs:     ids,
		},
	}
}

func NewTransactionsMessage(txs []*tx.Tx) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeTransactions,
		Payload: &payload.TransactionsPayload{
			Transactions: txs,
		},
	}
}
func NewQueryVoteMessage(querier peer.ID, height, round int) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeQueryVotes,
		Payload: &payload.QueryVotesPayload{
			Querier: querier,
			Height:  height,
			Round:   round,
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
func NewOpaqueBlockAnnounceMessage(height int, block *block.Block, commit *block.Commit) *Message {
	return NewBlockAnnounceMessage("", height, block, commit)
}

func NewBlockAnnounceMessage(peerID peer.ID, height int, block *block.Block, commit *block.Commit) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeBlockAnnounce,
		Payload: &payload.BlockAnnouncePayload{
			PeerID: peerID,
			Height: height,
			Block:  block,
			Commit: commit,
		},
	}
}

func NewDownloadRequestMessage(initiator, target peer.ID, sessionID int, from, to int) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeDownloadRequest,
		Payload: &payload.DownloadRequestPayload{
			Initiator: initiator,
			Target:    target,
			SessionID: sessionID,
			From:      from,
			To:        to,
		},
	}
}

func NewDownloadResponseMessage(code payload.ResponseCode, initiator, target peer.ID, sessionID, from int, blocks []*block.Block, txs []*tx.Tx) *Message {
	return &Message{
		Version: LastVersion,
		Type:    payload.PayloadTypeDownloadResponse,
		Payload: &payload.DownloadResponsePayload{
			ResponseCode: code,
			SessionID:    sessionID,
			Initiator:    initiator,
			Target:       target,
			From:         from,
			Blocks:       blocks,
			Transactions: txs,
		},
	}
}

func makePayload(t payload.PayloadType) payload.Payload {
	switch t {
	case payload.PayloadTypeSalam:
		return &payload.SalamPayload{}
	case payload.PayloadTypeAleyk:
		return &payload.AleykPayload{}
	case payload.PayloadTypeLatestBlocksRequest:
		return &payload.LatestBlocksRequestPayload{}
	case payload.PayloadTypeLatestBlocksResponse:
		return &payload.LatestBlocksResponsePayload{}
	case payload.PayloadTypeQueryTransactions:
		return &payload.QueryTransactionsPayload{}
	case payload.PayloadTypeTransactions:
		return &payload.TransactionsPayload{}
	case payload.PayloadTypeQueryProposal:
		return &payload.QueryProposalPayload{}
	case payload.PayloadTypeProposal:
		return &payload.ProposalPayload{}
	case payload.PayloadTypeHeartBeat:
		return &payload.HeartBeatPayload{}
	case payload.PayloadTypeQueryVotes:
		return &payload.QueryVotesPayload{}
	case payload.PayloadTypeVote:
		return &payload.VotePayload{}
	case payload.PayloadTypeBlockAnnounce:
		return &payload.BlockAnnouncePayload{}
	case payload.PayloadTypeDownloadRequest:
		return &payload.DownloadRequestPayload{}
	case payload.PayloadTypeDownloadResponse:
		return &payload.DownloadResponsePayload{}
	}

	//
	return nil
}
