package state

import (
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/state/param"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type ChainInfo struct {
	LastBlockHeight uint32
	LastBlockHash   hash.Hash
	LastBlockTime   time.Time

	TotalPower     int64
	CommitteePower int64
	CommitteeSize  int

	TotalAccounts    int32
	TotalValidators  int32
	ActiveValidators int32

	IsPruned      bool
	PruningHeight uint32
}

// CommitteeInfo holds committee validators, protocol versions, and total power.
type CommitteeInfo struct {
	Validators       []*validator.Validator
	ProtocolVersions map[protocol.Version]float64
	CommitteePower   int64
}

type Facade interface {
	Genesis() *genesis.Genesis
	Params() *param.Params
	LastBlockHeight() uint32
	LastBlockHash() hash.Hash
	LastBlockTime() time.Time
	LastCertificate() *certificate.Certificate
	UpdateLastCertificate(v *vote.Vote) error
	ProposeBlock(valKey *bls.ValidatorKey, rewardAddr crypto.Address) (*block.Block, error)
	ValidateBlock(blk *block.Block, round int16) error
	CommitBlock(blk *block.Block, cert *certificate.Certificate) error
	CommitteeValidators() []*validator.Validator
	CommitteeInfo() *CommitteeInfo
	IsInCommittee(addr crypto.Address) bool
	Proposer(round int16) *validator.Validator
	IsProposer(addr crypto.Address, round int16) bool
	PendingTx(txID tx.ID) *tx.Tx
	AddPendingTx(trx *tx.Tx) error
	AddPendingTxAndBroadcast(trx *tx.Tx) error
	CommittedBlock(height uint32) (*store.CommittedBlock, error)
	CommittedTx(txID tx.ID) (*store.CommittedTx, error)
	BlockHash(height uint32) hash.Hash
	BlockHeight(h hash.Hash) uint32
	AccountByAddress(addr crypto.Address) (*account.Account, error)
	ValidatorByAddress(addr crypto.Address) (*validator.Validator, error)
	ValidatorByNumber(number int32) (*validator.Validator, error)
	ValidatorAddresses() []crypto.Address
	UpdateValidatorProtocolVersion(addr crypto.Address, ver protocol.Version)
	CalculateFee(amt amount.Amount, payloadType payload.Type) amount.Amount
	PublicKey(addr crypto.Address) (crypto.PublicKey, error)
	AvailabilityScore(valNum int32) float64
	AllPendingTxs() []*tx.Tx
	ChainInfo() *ChainInfo

	Close()
}
