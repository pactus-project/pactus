package state

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

type sandbox struct {
	lk deadlock.RWMutex

	state      *state
	store      *store.Store
	accounts   map[crypto.Address]accountStatus
	validators map[crypto.Address]validatorStatus
}

type validatorStatus struct {
	validator *validator.Validator
	updated   bool
	addToSet  bool
}

type accountStatus struct {
	account *account.Account
	updated bool
}

func newSandbox(store *store.Store, state *state) *sandbox {
	sb := &sandbox{
		state:      state,
		store:      store,
		accounts:   make(map[crypto.Address]accountStatus),
		validators: make(map[crypto.Address]validatorStatus),
	}
	return sb
}

func (sb *sandbox) Reset() {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.reset()
}

func (sb *sandbox) reset() {
	sb.accounts = make(map[crypto.Address]accountStatus)
	sb.validators = make(map[crypto.Address]validatorStatus)
}

func (sb *sandbox) commit(set *validator.ValidatorSet) error {
	for _, acc := range sb.accounts {
		if acc.updated {
			sb.store.UpdateAccount(acc.account)
		}
	}

	for _, val := range sb.validators {
		if val.updated {
			sb.store.UpdateValidator(val.validator)
		}

		if val.addToSet {

		}
	}

	sb.reset()

	return nil
}

func (sb *sandbox) Account(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.accounts[addr]
	if ok {
		return s.account
	}

	acc, err := sb.store.Account(addr)
	if err != nil {
		return nil
	}
	return acc
}

func (sb *sandbox) UpdateAccount(acc *account.Account) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := acc.Address()
	s, ok := sb.accounts[addr]
	if ok {
		s.account = acc
		s.updated = true
	} else {
		sb.accounts[addr] = accountStatus{
			account: acc,
			updated: true,
		}
	}
}

func (sb *sandbox) Validator(addr crypto.Address) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if ok {
		return s.validator
	}

	val, err := sb.store.Validator(addr)
	if err != nil {
		return nil
	}
	return val
}

func (sb *sandbox) UpdateValidator(val *validator.Validator) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := val.Address()
	s, ok := sb.validators[addr]
	if ok {
		s.validator = val
		s.updated = true
	} else {
		sb.validators[addr] = validatorStatus{
			validator: val,
			updated:   true,
		}
	}
}

func (sb *sandbox) AddToSet(val *validator.Validator) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := val.Address()
	s, ok := sb.validators[addr]
	if ok {
		s.validator = val
		s.addToSet = true
	} else {
		sb.validators[addr] = validatorStatus{
			validator: val,
			addToSet:  true,
		}
	}
}

func (sb *sandbox) CurrentHeight() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.state.lastBlockHeight + 1
}

func (sb *sandbox) MaxMemoLenght() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.state.params.MaximumMemoLength
}

func (sb *sandbox) FeeFraction() float64 {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.state.params.FeeFraction
}

func (sb *sandbox) MinFee() int64 {
	return sb.state.params.MinimumFee
}

func (sb *sandbox) TransactionToLiveInterval() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.state.params.TransactionToLiveInterval
}

func (sb *sandbox) RecentBlockHeight(hash crypto.Hash) int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	h, has := sb.state.recentBlockHashes.Get(hash)
	if !has {
		return -1
	}

	return h.(int)
}

func (sb *sandbox) LastBlockHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.state.lastBlockHeight
}

func (sb *sandbox) LastBlockHash() crypto.Hash {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.state.lastBlockHash
}
