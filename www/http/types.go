package http

import (
	"time"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/sync/peerset"
	"github.com/zarbchain/zarb-go/tx"
)

type BlockchainResult struct {
	Height int
}
type BlockResult struct {
	Hash  hash.Hash
	Time  time.Time
	Data  string
	Block *block.Block
}

type TransactionResult struct {
	ID   hash.Hash
	Data string
	Tx   tx.Tx
}

type SendTranscationResult struct {
	Status int
	ID     hash.Hash
}

type NetworkResult struct {
	ID    peer.ID
	Peers []*peerset.Peer
}
