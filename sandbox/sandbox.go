package sandbox

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

type SandboxConcrete struct {
	lk deadlock.RWMutex

	store            store.StoreReader
	sortition        *sortition.Sortition
	committee        committee.CommitteeReader
	accounts         map[crypto.Address]*AccountStatus
	validators       map[crypto.Address]*ValidatorStatus
	params           param.Params
	lastHeight       int
	totalAccounts    int
	totalValidators  int
	totalStakeChange int64
}

type ValidatorStatus struct {
	Validator       validator.Validator
	Updated         bool
	JoinedCommittee bool
}

type AccountStatus struct {
	Account account.Account
	Updated bool
}

func NewSandbox(store store.StoreReader, params param.Params, lastHeight int, sortition *sortition.Sortition, committee committee.CommitteeReader) *SandboxConcrete {
	sb := &SandboxConcrete{
		store:      store,
		sortition:  sortition,
		committee:  committee,
		lastHeight: lastHeight,
		params:     params,
	}

	sb.accounts = make(map[crypto.Address]*AccountStatus)
	sb.validators = make(map[crypto.Address]*ValidatorStatus)
	sb.totalAccounts = sb.store.TotalAccounts()
	sb.totalValidators = sb.store.TotalValidators()
	sb.totalStakeChange = 0

	return sb
}

func (sb *SandboxConcrete) shouldPanicForDuplicatedAddress() {
	//
	// Why we should panic here?
	//
	// Try to make a new item which already exists in store.
	//
	logger.Panic("Duplicated address")
}

func (sb *SandboxConcrete) shouldPanicForUnknownAddress() {
	//
	// Why we should panic here?
	//
	// We only update accounts or validators which we have them inside the sandbox.
	// We must either make a new one (i.e. `MakeNewAccount`) or get it from store (i.e. `Account`) in advance.
	//
	logger.Panic("Unknown address")
}

func (sb *SandboxConcrete) Account(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.accounts[addr]
	if ok {
		copy := new(account.Account)
		*copy = s.Account
		return copy
	}

	acc, err := sb.store.Account(addr)
	if err != nil {
		return nil
	}
	sb.accounts[addr] = &AccountStatus{
		Account: *acc,
	}

	return acc
}
func (sb *SandboxConcrete) MakeNewAccount(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	if sb.store.HasAccount(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	acc := account.NewAccount(addr, sb.totalAccounts)
	sb.accounts[addr] = &AccountStatus{
		Account: *acc,
		Updated: true,
	}
	sb.totalAccounts++
	return acc
}

func (sb *SandboxConcrete) UpdateAccount(acc *account.Account) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := acc.Address()
	s, ok := sb.accounts[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}
	s.Account = *acc
	s.Updated = true
}

func (sb *SandboxConcrete) Validator(addr crypto.Address) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if ok {
		copy := new(validator.Validator)
		*copy = s.Validator
		return copy
	}

	val, err := sb.store.Validator(addr)
	if err != nil {
		return nil
	}
	sb.validators[addr] = &ValidatorStatus{
		Validator: *val,
	}
	return val
}

func (sb *SandboxConcrete) MakeNewValidator(pub crypto.PublicKey) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := pub.Address()
	if sb.store.HasValidator(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	val := validator.NewValidator(pub, sb.totalValidators, sb.lastHeight+1)
	sb.validators[addr] = &ValidatorStatus{
		Validator: *val,
		Updated:   true,
	}
	sb.totalValidators++
	return val
}

func (sb *SandboxConcrete) UpdateValidator(val *validator.Validator) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := val.Address()
	s, ok := sb.validators[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}

	sb.totalStakeChange += val.Stake() - s.Validator.Stake()
	s.Validator = *val
	s.Updated = true
}

func (sb *SandboxConcrete) EnterCommittee(blockHash crypto.Hash, addr crypto.Address) error {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if !ok {
		return errors.Errorf(errors.ErrGeneric, "unknown validator")
	}

	if sb.committee.Contains(addr) {
		return errors.Errorf(errors.ErrGeneric, "this validator already is in the committee")
	}

	joined := 0
	for _, s := range sb.validators {
		if s.JoinedCommittee {
			joined++
		}
	}
	if joined >= (sb.params.CommitteeSize / 3) {
		return errors.Errorf(errors.ErrGeneric, "in each height only 1/3 of validator can be changed")
	}
	h, _ := sb.store.BlockHeight(blockHash)
	b, err := sb.store.Block(h)
	if err != nil {
		return errors.Errorf(errors.ErrGeneric, "invalid block hash")
	}
	commiters := b.LastCertificate().Committers()
	for _, num := range commiters {
		if s.Validator.Number() == num {
			return errors.Errorf(errors.ErrGeneric, "this validator was in the committee in time of sending the sortition")
		}
	}

	s.JoinedCommittee = true
	return nil
}

func (sb *SandboxConcrete) MaxMemoLength() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.params.MaximumMemoLength
}

func (sb *SandboxConcrete) FeeFraction() float64 {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.params.FeeFraction
}

func (sb *SandboxConcrete) MinFee() int64 {
	return sb.params.MinimumFee
}

func (sb *SandboxConcrete) TransactionToLiveInterval() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.TransactionToLiveInterval
}

func (sb *SandboxConcrete) BlockHeight(hash crypto.Hash) int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	if hash.EqualsTo(crypto.UndefHash) {
		return 0
	}

	h, err := sb.store.BlockHeight(hash)
	if err != nil {
		return -1
	}

	return h
}

func (sb *SandboxConcrete) CurrentHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.lastHeight + 1
}

func (sb *SandboxConcrete) LastHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.lastHeight
}

func (sb *SandboxConcrete) LastBlockHash() crypto.Hash {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	b, _ := sb.store.Block(sb.lastHeight)
	return b.Hash()
}

func (sb *SandboxConcrete) VerifySortition(blockHash crypto.Hash, proof sortition.Proof, val *validator.Validator) bool {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.sortition.VerifyProof(blockHash, proof, val.PublicKey(), val.Stake())
}

func (sb *SandboxConcrete) IterateAccounts(consumer func(*AccountStatus)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, as := range sb.accounts {
		consumer(as)
	}
}

func (sb *SandboxConcrete) IterateValidators(consumer func(*ValidatorStatus)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, vs := range sb.validators {
		consumer(vs)
	}
}

func (sb *SandboxConcrete) CommitteeSize() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.CommitteeSize
}

func (sb *SandboxConcrete) IsInCommittee(addr crypto.Address) bool {
	return sb.committee.Contains(addr)
}
