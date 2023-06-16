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
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/param"
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
	params          param.Params
	totalAccounts   int32
	totalValidators int32
}

type sandboxValidator struct {
	validator *validator.Validator
	updated   bool
}

type sandboxAccount struct {
	account *account.Account
	updated bool
}

func NewSandbox(store store.Reader, params param.Params, committee committee.Reader) Sandbox {
	sb := &sandbox{
		store:     store,
		committee: committee,
		params:    params,
	}

	sb.accounts = make(map[crypto.Address]*sandboxAccount)
	sb.validators = make(map[crypto.Address]*sandboxValidator)
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

func (sb *sandbox) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := pub.Address()
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

func (sb *sandbox) IterateAccounts(consumer func(crypto.Address, *account.Account, bool)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for addr, sa := range sb.accounts {
		consumer(addr, sa.account, sa.updated)
	}
}

func (sb *sandbox) IterateValidators(consumer func(*validator.Validator, bool)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, sv := range sb.validators {
		consumer(sv.validator, sv.updated)
	}
}

func (sb *sandbox) RecentBlockByStamp(stamp hash.Stamp) (uint32, *block.Block) {
	return sb.store.RecentBlockByStamp(stamp)
}

func (sb *sandbox) Committee() committee.Reader {
	return sb.committee
}

// TODO: write test for me.
// VerifyProof verifies proof of a sortition transaction.
func (sb *sandbox) VerifyProof(stamp hash.Stamp, proof sortition.Proof, val *validator.Validator) bool {
	_, b := sb.store.RecentBlockByStamp(stamp)
	if b == nil {
		return false
	}
	seed := b.Header().SortitionSeed()
	total := int64(0) // TODO: we can get it from state
	sb.store.IterateValidators(func(val *validator.Validator) bool {
		total += val.Power()
		return false
	})
	return sortition.VerifyProof(seed, proof, val.PublicKey(), total, val.Power())
}
