package state

import (
	"time"

	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
)

type StateReader interface {
	StoreReader() store.StoreReader
	Committee() committee.CommitteeReader
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
	ProposeBlock(round int) (*block.Block, error)
	ValidateBlock(block block.Block) error
	CommitBlock(height int, block block.Block, commit block.Commit) error
}
