package state

import (
	"time"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/block"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/validator"
)

type StateFacade interface {
	GenesisHash() crypto.Hash
	LastBlockHeight() int
	LastBlockHash() crypto.Hash
	LastBlockTime() time.Time
	LastCertificate() *block.Certificate
	BlockTime() time.Duration
	UpdateLastCertificate(lastCertificate *block.Certificate) error
	ProposeBlock(round int) (*block.Block, error)
	ValidateBlock(block *block.Block) error
	CommitBlock(height int, block *block.Block, cert *block.Certificate) error
	CommitteeValidators() []*validator.Validator
	IsInCommittee(addr crypto.Address) bool
	Proposer(round int) *validator.Validator
	IsProposer(addr crypto.Address, round int) bool
	Transaction(id tx.ID) *tx.Tx
	PendingTx(id tx.ID) *tx.Tx
	AddPendingTx(trx *tx.Tx) error
	AddPendingTxAndBroadcast(trx *tx.Tx) error
	Block(height int) *block.Block
	BlockHeight(hash crypto.Hash) int
	Account(addr crypto.Address) *account.Account
	Validator(addr crypto.Address) *validator.Validator
	ValidatorByNumber(number int) *validator.Validator
	Close() error
	Fingerprint() string
}
