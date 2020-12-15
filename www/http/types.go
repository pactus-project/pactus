package http

import (
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
)

type BlockResult struct {
	Hash  crypto.Hash
	Time  time.Time
	Data  string
	Block block.Block
}
