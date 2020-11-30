package state

import (
	"bytes"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/libs/orderedmap"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

type sandbox struct {
	lk deadlock.RWMutex

	state      *state
	store      *store.Store
	valChanges *orderedmap.OrderedMap
	accChanges *orderedmap.OrderedMap
}

type validatorInfo struct {
	validator *validator.Validator
}

type accountInfo struct {
	account *account.Account
}

func lessFn(l, r interface{}) bool {
	return bytes.Compare(l.(crypto.Address).RawBytes(), r.(crypto.Address).RawBytes()) < 0
}

func newSandbox(store *store.Store, state *state) *sandbox {
	sb := &sandbox{
		state:      state,
		store:      store,
		valChanges: orderedmap.NewMap(lessFn),
		accChanges: orderedmap.NewMap(lessFn),
	}
	return sb
}

func (sb *sandbox) reset() {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.accChanges = orderedmap.NewMap(lessFn)
	sb.valChanges = orderedmap.NewMap(lessFn)
}

func (sb *sandbox) commit(set *validator.ValidatorSet) error {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.accChanges.Iter(func(key, value interface{}) (more bool) {
		i := value.(*accountInfo)

		sb.store.UpdateAccount(i.account)

		return true
	})

	sb.valChanges.Iter(func(key, value interface{}) (more bool) {
		i := value.(*validatorInfo)
		sb.store.UpdateValidator(i.validator)
		return true
	})

	/// reset sandbox
	sb.accChanges = orderedmap.NewMap(lessFn)
	sb.valChanges = orderedmap.NewMap(lessFn)

	return nil
}

func (sb *sandbox) HasAccount(addr crypto.Address) bool {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	_, ok := sb.accChanges.GetOk(addr)
	if ok {
		return true
	}

	return sb.store.HasAccount(addr)
}

func (sb *sandbox) Account(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	i, ok := sb.accChanges.GetOk(addr)
	if ok {
		return i.(*accountInfo).account
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
	i, ok := sb.accChanges.GetOk(addr)
	if ok {
		i.(*accountInfo).account = acc
	} else {
		sb.accChanges.Set(addr, &accountInfo{account: acc})
	}
}

func (sb *sandbox) HasValidator(addr crypto.Address) bool {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	_, ok := sb.valChanges.GetOk(addr)
	if ok {
		return true
	}

	return sb.store.HasValidator(addr)
}

func (sb *sandbox) Validator(addr crypto.Address) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	i, ok := sb.valChanges.GetOk(addr)
	if ok {
		return i.(*validatorInfo).validator
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
	i, ok := sb.valChanges.GetOk(addr)
	if ok {
		i.(*validatorInfo).validator = val
	} else {
		sb.valChanges.Set(addr, &validatorInfo{validator: val})
	}
}

func (sb *sandbox) CurrentHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.state.LastBlockHeight() + 1
}

func (sb *sandbox) MaxMemoLenght() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.state.params.MaximumMemoLength
}

func (sb *sandbox) FeeFraction() float64 {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.state.params.FeeFraction
}

func (sb *sandbox) MinFee() int64 {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

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
