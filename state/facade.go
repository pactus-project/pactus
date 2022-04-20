package state

import (
	"time"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

type Facade interface {
	GenesisHash() hash.Hash
	LastBlockHeight() int32
	LastBlockHash() hash.Hash
	LastBlockTime() time.Time
	LastCertificate() *block.Certificate
	BlockTime() time.Duration
	UpdateLastCertificate(lastCertificate *block.Certificate) error
	ProposeBlock(round int16) (*block.Block, error)
	ValidateBlock(block *block.Block) error
	CommitBlock(height int32, block *block.Block, cert *block.Certificate) error
	CommitteeValidators() []*validator.Validator
	IsInCommittee(addr crypto.Address) bool
	Proposer(round int16) *validator.Validator
	IsProposer(addr crypto.Address, round int16) bool
	TotalPower() int64
	CommitteePower() int64
	Transaction(id tx.ID) *tx.Tx
	PendingTx(id tx.ID) *tx.Tx
	AddPendingTx(trx *tx.Tx) error
	AddPendingTxAndBroadcast(trx *tx.Tx) error
	Block(hash hash.Hash) *block.Block // TODO: return store block (including block header data)
	BlockHash(height int32) hash.Hash
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
