package sandbox

import (
	"sync"

	"github.com/pactus-project/pactus/committee"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/sortition"
	"github.com/pactus-project/pactus/store"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/param"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/logger"
)

var _ Sandbox = &sandbox{}

type sandbox struct {
	lk sync.RWMutex

	store           store.Reader
	committee       committee.Reader
	accounts        map[crypto.Address]*AccountStatus
	validators      map[crypto.Address]*ValidatorStatus
	params          param.Params
	totalAccounts   int32
	totalValidators int32
}

type ValidatorStatus struct {
	Validator validator.Validator
	Updated   bool
}

type AccountStatus struct {
	Account account.Account
	Updated bool
}

func NewSandbox(store store.Reader, params param.Params, committee committee.Reader) Sandbox {
	sb := &sandbox{
		store:     store,
		committee: committee,
		params:    params,
	}

	sb.accounts = make(map[crypto.Address]*AccountStatus)
	sb.validators = make(map[crypto.Address]*ValidatorStatus)
	sb.totalAccounts = sb.store.TotalAccounts()
	sb.totalValidators = sb.store.TotalValidators()

	return sb
}

func (sb *sandbox) shouldPanicForDuplicatedAddress() {
	//
	// Why we should panic here?
	//
	// Try to make a new item which already exists in store.
	//
	logger.Panic("duplicated address")
}

func (sb *sandbox) shouldPanicForUnknownAddress() {
	//
	// Why we should panic here?
	//
	// We only update accounts or validators which we have them inside the sandbox.
	// We must either make a new one (i.e. `MakeNewAccount`) or get it from store (i.e. `Account`) in advance.
	//
	logger.Panic("unknown address")
}

func (sb *sandbox) Account(addr crypto.Address) *account.Account {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

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
func (sb *sandbox) MakeNewAccount(addr crypto.Address) *account.Account {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

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

func (sb *sandbox) UpdateAccount(acc *account.Account) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	addr := acc.Address()
	s, ok := sb.accounts[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}
	s.Account = *acc
	s.Updated = true
}

func (sb *sandbox) Validator(addr crypto.Address) *validator.Validator {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

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

func (sb *sandbox) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	addr := pub.Address()
	if sb.store.HasValidator(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	val := validator.NewValidator(pub, sb.totalValidators)
	sb.validators[addr] = &ValidatorStatus{
		Validator: *val,
		Updated:   true,
	}
	sb.totalValidators++
	return val
}

func (sb *sandbox) UpdateValidator(val *validator.Validator) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	addr := val.Address()
	s, ok := sb.validators[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}

	s.Validator = *val
	s.Updated = true
}

func (sb *sandbox) Params() param.Params {
	return sb.params
}

func (sb *sandbox) CurrentHeight() uint32 {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.currentHeight()
}

func (sb *sandbox) currentHeight() uint32 {
	h, _ := sb.store.LastCertificate()

	return h + 1
}

func (sb *sandbox) IterateAccounts(consumer func(*AccountStatus)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, as := range sb.accounts {
		consumer(as)
	}
}

func (sb *sandbox) IterateValidators(consumer func(*ValidatorStatus)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, vs := range sb.validators {
		consumer(vs)
	}
}

func (sb *sandbox) FindBlockHashByStamp(stamp hash.Stamp) (hash.Hash, bool) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.store.FindBlockHashByStamp(stamp)
}

func (sb *sandbox) FindBlockHeightByStamp(stamp hash.Stamp) (uint32, bool) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.store.FindBlockHeightByStamp(stamp)
}

func (sb *sandbox) Committee() committee.Reader {
	return sb.committee
}

// TODO: write test for me.
// VerifyProof verifies proof of a sortition transaction.
func (sb *sandbox) VerifyProof(stamp hash.Stamp, proof sortition.Proof, val *validator.Validator) bool {
	hash, ok := sb.store.FindBlockHashByStamp(stamp)
	if !ok {
		return false
	}
	bi, err := sb.store.Block(hash)
	if err != nil {
		return false
	}
	b, _ := bi.ToFullBlock()
	seed := b.Header().SortitionSeed()

	total := int64(0) // TODO: we can get it from state
	sb.store.IterateValidators(func(val *validator.Validator) bool {
		total += val.Power()
		return false
	})
	return sortition.VerifyProof(seed, proof, val.PublicKey(), total, val.Power())
}
