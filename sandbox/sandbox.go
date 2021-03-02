package sandbox

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
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
	validatorSet     validator.ValidatorSetReader
	accounts         map[crypto.Address]*AccountStatus
	validators       map[crypto.Address]*ValidatorStatus
	recentBlocks     *linkedmap.LinkedMap
	params           param.Params
	totalAccounts    int
	totalValidators  int
	totalStakeChange int64
}

type ValidatorStatus struct {
	Validator validator.Validator
	Updated   bool
	AddToSet  bool
}

type AccountStatus struct {
	Account account.Account
	Updated bool
}

func NewSandbox(store store.StoreReader, params param.Params, lastBlockHeight int, sortition *sortition.Sortition, valset validator.ValidatorSetReader) (*SandboxConcrete, error) {
	sb := &SandboxConcrete{
		store:        store,
		sortition:    sortition,
		validatorSet: valset,
		params:       params,
		recentBlocks: linkedmap.NewLinkedMap(params.TransactionToLiveInterval),
		accounts:     make(map[crypto.Address]*AccountStatus),
		validators:   make(map[crypto.Address]*ValidatorStatus),
	}

	// First, let add genesis block (Block 0) hash
	sb.recentBlocks.PushBack(crypto.UndefHash, 0)

	// Now we try to fetch recent block hashes
	// Block zero will be kicked out of the list if we have enough blocks
	from := lastBlockHeight - params.TransactionToLiveInterval
	if from <= 0 {
		from = 1
	}
	to := lastBlockHeight
	for i := from; i <= to; i++ {
		b, err := store.Block(i)
		if err != nil {
			return nil, err
		}
		sb.recentBlocks.PushBack(b.Hash(), i)
	}

	// To update total accounts and validator counters
	sb.clear()

	return sb, nil
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

func (sb *SandboxConcrete) Clear() {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.clear()
}

func (sb *SandboxConcrete) clear() {
	sb.accounts = make(map[crypto.Address]*AccountStatus)
	sb.validators = make(map[crypto.Address]*ValidatorStatus)
	sb.totalAccounts = sb.store.TotalAccounts()
	sb.totalValidators = sb.store.TotalValidators()
	sb.totalStakeChange = 0
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

	val := validator.NewValidator(pub, sb.totalValidators, sb.lastHeight()+1)
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

func (sb *SandboxConcrete) AddToSet(blockHash crypto.Hash, addr crypto.Address) error {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if !ok {
		return errors.Errorf(errors.ErrGeneric, "Unknown validator")
	}

	if sb.validatorSet.Contains(addr) {
		return errors.Errorf(errors.ErrGeneric, "This validator already is in the set")
	}

	joined := 0
	for _, s := range sb.validators {
		if s.AddToSet {
			joined++
		}
	}
	if joined >= (sb.params.CommitteeSize / 3) {
		return errors.Errorf(errors.ErrGeneric, "In each height only 1/3 of validator can be changed")
	}
	h, _ := sb.store.BlockHeight(blockHash)
	b, err := sb.store.Block(h)
	if err != nil {
		return errors.Errorf(errors.ErrGeneric, "Invalid block hash")
	}
	commiters := b.LastCommit().Committers()
	for _, c := range commiters {
		if s.Validator.Number() == c.Number {
			return errors.Errorf(errors.ErrGeneric, "This validator was in the set in time of sending the sortition")
		}
	}

	s.AddToSet = true
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

func (sb *SandboxConcrete) RecentBlockHeight(hash crypto.Hash) int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	h, has := sb.recentBlocks.Get(hash)
	if !has {
		return -1
	}

	return h.(int)
}

func (sb *SandboxConcrete) lastHeight() int {
	_, v := sb.recentBlocks.Last()
	if v == nil {
		return -1
	}
	return v.(int)
}

func (sb *SandboxConcrete) CurrentHeight() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.lastHeight() + 1
}

func (sb *SandboxConcrete) LastBlockHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.lastHeight()
}

func (sb *SandboxConcrete) LastBlockHash() crypto.Hash {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	k, _ := sb.recentBlocks.Last()
	return k.(crypto.Hash)
}

func (sb *SandboxConcrete) AppendNewBlock(hash crypto.Hash, height int) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.recentBlocks.PushBack(hash, height)
}

func (sb *SandboxConcrete) VerifySortition(blockHash crypto.Hash, proof sortition.Proof, val *validator.Validator) bool {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	h, _ := sb.store.BlockHeight(blockHash)
	b, err := sb.store.Block(h)
	if err != nil {
		return false
	}

	return sb.sortition.VerifyProof(b.Header().SortitionSeed(), proof, val.PublicKey(), val.Stake())
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

func (sb *SandboxConcrete) TotalStakeChange() int64 {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.totalStakeChange
}
