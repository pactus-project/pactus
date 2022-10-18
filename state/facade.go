package state

import (
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
)

type Facade interface {
	GenesisHash() hash.Hash
	LastBlockHeight() uint32
	LastBlockHash() hash.Hash
	LastBlockTime() time.Time
	LastCertificate() *block.Certificate
	BlockTime() time.Duration
	UpdateLastCertificate(lastCertificate *block.Certificate) error
	ProposeBlock(round int16) (*block.Block, error)
	ValidateBlock(block *block.Block) error
	CommitBlock(height uint32, block *block.Block, cert *block.Certificate) error
	CommitteeValidators() []*validator.Validator
	IsInCommittee(addr crypto.Address) bool
	Proposer(round int16) *validator.Validator
	IsProposer(addr crypto.Address, round int16) bool
	IsValidator(addr crypto.Address) bool
	TotalPower() int64
	CommitteePower() int64
	PendingTx(id tx.ID) *tx.Tx
	AddPendingTx(trx *tx.Tx) error
	AddPendingTxAndBroadcast(trx *tx.Tx) error
	StoredBlock(height uint32) *store.StoredBlock
	StoredTx(id tx.ID) *store.StoredTx
	BlockHash(height uint32) hash.Hash
	BlockHeight(hash hash.Hash) uint32
	AccountByAddress(addr crypto.Address) *account.Account
	ValidatorByAddress(addr crypto.Address) *validator.Validator
	ValidatorByNumber(number int32) *validator.Validator
	Params() param.Params
	// RewardAddress returns the rewards address that is associated to this node.
	// Reward address can be set through the config file,
	// and if it is not set, it will be the same as validator address
	RewardAddress() crypto.Address

	// ValidatorAddress return the validator address that is associated to this node
	// Validator address is different from the reward address.
	ValidatorAddress() crypto.Address
	Close() error
	Fingerprint() string
}
