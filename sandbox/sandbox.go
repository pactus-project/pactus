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
	lk sync.RWMutex // TODO: can we get rid of this lock?

	store           store.Reader
	committee       committee.Reader
	accounts        map[crypto.Address]*accountStatus
	validators      map[crypto.Address]*validatorStatus
	params          param.Params
	totalAccounts   int32
	totalValidators int32
}

type validatorStatus struct {
	Validator validator.Validator
	Updated   bool
	Joined    bool // Is joined committee
}

type accountStatus struct {
	Account account.Account
	Updated bool
}

func NewSandbox(store store.Reader, params param.Params, committee committee.Reader) Sandbox {
	sb := &sandbox{
		store:     store,
		committee: committee,
		params:    params,
	}

	sb.accounts = make(map[crypto.Address]*accountStatus)
	sb.validators = make(map[crypto.Address]*validatorStatus)
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
	// In order to proceed, the account must have been created and present beforehand.
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
		clone := new(account.Account)
		*clone = s.Account
		return clone
	}

	acc, err := sb.store.Account(addr)
	if err != nil {
		return nil
	}
	sb.accounts[addr] = &accountStatus{
		Account: *acc,
	}

	return acc
}
func (sb *sandbox) MakeNewAccount(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	if sb.store.HasAccount(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	acc := account.NewAccount(sb.totalAccounts)
	sb.accounts[addr] = &accountStatus{
		Account: *acc,
		Updated: true,
	}
	sb.totalAccounts++
	return acc
}

func (sb *sandbox) UpdateAccount(addr crypto.Address, acc *account.Account) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.accounts[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}
	s.Account = *acc
	s.Updated = true
}

func (sb *sandbox) Validator(addr crypto.Address) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if ok {
		clone := new(validator.Validator)
		*clone = s.Validator
		return clone
	}

	val, err := sb.store.Validator(addr)
	if err != nil {
		return nil
	}
	sb.validators[addr] = &validatorStatus{
		Validator: *val,
	}
	return val
}

func (sb *sandbox) JoinedToCommittee(addr crypto.Address) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}

	s.Joined = true
}

func (sb *sandbox) IsJoinedCommittee(addr crypto.Address) bool {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	s, ok := sb.validators[addr]
	if ok {
		return s.Joined
	}
	return false
}

func (sb *sandbox) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := pub.Address()
	if sb.store.HasValidator(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	val := validator.NewValidator(pub, sb.totalValidators)
	sb.validators[addr] = &validatorStatus{
		Validator: *val,
		Updated:   true,
	}
	sb.totalValidators++
	return val
}

func (sb *sandbox) UpdateValidator(val *validator.Validator) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

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

func (sb *sandbox) IterateAccounts(
	consumer func(crypto.Address, *account.Account, bool)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for addr, as := range sb.accounts {
		consumer(addr, &as.Account, as.Updated)
	}
}

func (sb *sandbox) IterateValidators(
	consumer func(*validator.Validator, bool, bool)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, vs := range sb.validators {
		consumer(&vs.Validator, vs.Updated, vs.Joined)
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
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	height, ok := sb.store.FindBlockHeightByStamp(stamp)
	if !ok {
		return false
	}
	storedBlock, err := sb.store.Block(height)
	if err != nil {
		return false
	}
	seed := storedBlock.ToBlock().Header().SortitionSeed()
	total := int64(0) // TODO: we can get it from state
	sb.store.IterateValidators(func(val *validator.Validator) bool {
		total += val.Power()
		return false
	})
	return sortition.VerifyProof(seed, proof, val.PublicKey(), total, val.Power())
}
