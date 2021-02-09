package http

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
)

type BlockchainResult struct {
	Height int
}
type BlockResult struct {
	Hash  crypto.Hash
	Time  time.Time
	Data  string
	Block block.Block
}

type ReceiptResult struct {
	Hash    crypto.Hash
	Data    string
	Receipt tx.Receipt
}
type TransactionResult struct {
	ID      crypto.Hash
	Data    string
	Tx      tx.Tx
	Receipt ReceiptResult
}

type SendTranscationResult struct {
	Status int
	ID     crypto.Hash
}

type NetworkResult struct {
	ID    peer.ID
	Peers []*peerset.Peer
}
