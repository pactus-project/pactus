package state

import (
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/types/vote"
)

type Facade interface {
	Genesis() *genesis.Genesis
	LastBlockHeight() uint32
	LastBlockHash() hash.Hash
	LastBlockTime() time.Time
	LastCertificate() *certificate.Certificate
	UpdateLastCertificate(v *vote.Vote) error
	ProposeBlock(valKey *bls.ValidatorKey, rewardAddr crypto.Address) (*block.Block, error)
	ValidateBlock(blk *block.Block, round int16) error
	CommitBlock(blk *block.Block, cert *certificate.Certificate) error
	CommitteeValidators() []*validator.Validator
	IsInCommittee(addr crypto.Address) bool
	Proposer(round int16) *validator.Validator
	IsProposer(addr crypto.Address, round int16) bool
	IsValidator(addr crypto.Address) bool
	TotalPower() int64
	TotalAccounts() int32
	TotalValidators() int32
	CommitteePower() int64
	PendingTx(id tx.ID) *tx.Tx
	AddPendingTx(trx *tx.Tx) error
	AddPendingTxAndBroadcast(trx *tx.Tx) error
	CommittedBlock(height uint32) *store.CommittedBlock
	CommittedTx(id tx.ID) *store.CommittedTx
	BlockHash(height uint32) hash.Hash
	BlockHeight(h hash.Hash) uint32
	AccountByAddress(addr crypto.Address) *account.Account
	ValidatorByAddress(addr crypto.Address) *validator.Validator
	ValidatorByNumber(number int32) *validator.Validator
	ValidatorAddresses() []crypto.Address
	Params() *param.Params
	Close() error
	CalculateFee(amt amount.Amount, payloadType payload.Type) amount.Amount
	PublicKey(addr crypto.Address) (crypto.PublicKey, error)
	AvailabilityScore(valNum int32) float64
}
