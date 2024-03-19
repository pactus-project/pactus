package sandbox

import (
	"sync"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/logger"
)

var _ Sandbox = &sandbox{}

type sandbox struct {
	lk sync.RWMutex

	store           store.Reader
	committee       committee.Reader
	accounts        map[crypto.Address]*sandboxAccount
	validators      map[crypto.Address]*sandboxValidator
	committedTrxs   map[tx.ID]*tx.Tx
	params          *param.Params
	height          uint32
	totalAccounts   int32
	totalValidators int32
	totalPower      int64
	powerDelta      int64
	accumulatedFee  amount.Amount
}

type sandboxValidator struct {
	validator *validator.Validator
	updated   bool
	joined    bool // Is joined committee
}

type sandboxAccount struct {
	account *account.Account
	updated bool
}

func NewSandbox(height uint32, str store.Reader, params *param.Params,
	cmt committee.Reader, totalPower int64,
) Sandbox {
	sb := &sandbox{
		height:     height,
		store:      str,
		committee:  cmt,
		totalPower: totalPower,
		params:     params,
	}

	sb.accounts = make(map[crypto.Address]*sandboxAccount)
	sb.validators = make(map[crypto.Address]*sandboxValidator)
	sb.committedTrxs = make(map[tx.ID]*tx.Tx)
	sb.totalAccounts = sb.store.TotalAccounts()
	sb.totalValidators = sb.store.TotalValidators()

	return sb
}

func (sb *sandbox) shouldPanicForDuplicatedAddress() {
	//
	// Why is it necessary to panic here?
	//
	// An attempt is made to create a new item that already exists in the store.
	//
	logger.Panic("duplicated address")
}

func (sb *sandbox) shouldPanicForUnknownAddress() {
	//
	// Why is it necessary to panic here?
	//
	// We only update accounts or validators that are already present within the sandbox.
	// This can be achieved either by creating a new account using MakeNewAccount or
	// retrieving it from the store using Account.
	//
	logger.Panic("unknown address")
}

func (sb *sandbox) Account(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.accounts[addr]
	if ok {
		return s.account.Clone()
	}

	acc, err := sb.store.Account(addr)
	if err != nil {
		return nil
	}
	sb.accounts[addr] = &sandboxAccount{
		account: acc,
	}

	return acc.Clone()
}

func (sb *sandbox) MakeNewAccount(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	if sb.store.HasAccount(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	acc := account.NewAccount(sb.totalAccounts)
	sb.accounts[addr] = &sandboxAccount{
		account: acc,
		updated: true,
	}
	sb.totalAccounts++

	return acc.Clone()
}

// This function takes ownership of the account pointer.
// It is important that the caller should not modify the account data and
// keep it immutable.
func (sb *sandbox) UpdateAccount(addr crypto.Address, acc *account.Account) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.accounts[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}
	s.account = acc
	s.updated = true
}

func (sb *sandbox) AnyRecentTransaction(txID tx.ID) bool {
	if sb.committedTrxs[txID] != nil {
		return true
	}

	return sb.store.AnyRecentTransaction(txID)
}

func (sb *sandbox) Validator(addr crypto.Address) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if ok {
		return s.validator.Clone()
	}

	val, err := sb.store.Validator(addr)
	if err != nil {
		return nil
	}
	sb.validators[addr] = &sandboxValidator{
		validator: val,
	}

	return val.Clone()
}

func (sb *sandbox) JoinedToCommittee(addr crypto.Address) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}

	s.joined = true
}

func (sb *sandbox) IsJoinedCommittee(addr crypto.Address) bool {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if ok {
		return s.joined
	}

	return false
}

func (sb *sandbox) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := pub.ValidatorAddress()
	if sb.store.HasValidator(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	val := validator.NewValidator(pub, sb.totalValidators)
	sb.validators[addr] = &sandboxValidator{
		validator: val,
		updated:   true,
	}
	sb.totalValidators++

	return val.Clone()
}

// This function takes ownership of the validator pointer.
// It is important that the caller should not modify the validator data and
// keep it immutable.
func (sb *sandbox) UpdateValidator(val *validator.Validator) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := val.Address()
	s, ok := sb.validators[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}

	s.validator = val
	s.updated = true
}

func (sb *sandbox) Params() *param.Params {
	return sb.params
}

func (sb *sandbox) CurrentHeight() uint32 {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.height + 1
}

func (sb *sandbox) IterateAccounts(
	consumer func(crypto.Address, *account.Account, bool),
) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for addr, sa := range sb.accounts {
		consumer(addr, sa.account, sa.updated)
	}
}

func (sb *sandbox) IterateValidators(
	consumer func(*validator.Validator, bool, bool),
) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, sv := range sb.validators {
		consumer(sv.validator, sv.updated, sv.joined)
	}
}

func (sb *sandbox) Committee() committee.Reader {
	return sb.committee
}

// UpdatePowerDelta updates the change in the total power of the blockchain.
// The delta is the amount of change in the total power and can be either positive or negative.
func (sb *sandbox) UpdatePowerDelta(delta int64) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.powerDelta += delta
}

func (sb *sandbox) PowerDelta() int64 {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.powerDelta
}

// VerifyProof verifies proof of a sortition transaction.
func (sb *sandbox) VerifyProof(blockHeight uint32, proof sortition.Proof, val *validator.Validator) bool {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	seed := sb.store.SortitionSeed(blockHeight)
	if seed == nil {
		return false
	}

	return sortition.VerifyProof(*seed, proof, val.PublicKey(), sb.totalPower, val.Power())
}

func (sb *sandbox) CommitTransaction(trx *tx.Tx) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.committedTrxs[trx.ID()] = trx
	sb.accumulatedFee += trx.Fee()
}

func (sb *sandbox) AccumulatedFee() amount.Amount {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.accumulatedFee
}
