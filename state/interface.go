package state

import (
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

type StateReader interface {
	StoreReader() store.StoreReader
	ValidatorSet() validator.ValidatorSetReader
	LastBlockHeight() int
	GenesisHash() crypto.Hash
	LastBlockHash() crypto.Hash
	LastBlockTime() time.Time
	LastCommit() *block.Commit
	BlockTime() time.Duration
	UpdateLastCommit(lastCommit *block.Commit) error
	Fingerprint() string
}

type State interface {
	StateReader

	Close() error
	ProposeBlock() block.Block
	ValidateBlock(block block.Block) error
	ApplyBlock(height int, block block.Block, commit block.Commit) error
}
