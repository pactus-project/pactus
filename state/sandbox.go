package state

import (
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

type sandbox struct {
	lk deadlock.RWMutex

	store           *store.Store
	sortition       *sortition.Sortition
	validatorSet    *validator.ValidatorSet
	accounts        map[crypto.Address]*accountStatus
	validators      map[crypto.Address]*validatorStatus
	recentBlocks    *linkedmap.LinkedMap
	params          Params
	totalAccounts   int
	totalValidators int
	changeToStake   int64
}

type validatorStatus struct {
	validator validator.Validator
	updated   bool
	addToSet  bool
}

type accountStatus struct {
	account account.Account
	updated bool
}

func newSandbox(store *store.Store, params Params, lastBlockHeight int, sortition *sortition.Sortition, valset *validator.ValidatorSet) (*sandbox, error) {
	sb := &sandbox{
		store:        store,
		sortition:    sortition,
		validatorSet: valset,
		params:       params,
		recentBlocks: linkedmap.NewLinkedMap(params.TransactionToLiveInterval),
		accounts:     make(map[crypto.Address]*accountStatus),
		validators:   make(map[crypto.Address]*validatorStatus),
	}

	// First, let add genesis block (Block 0) hash
	sb.recentBlocks.PushBack(crypto.UndefHash, 0)

	// Now we try to fetch recent block hashes
	// Block zero will be kicked out of the list if we have enough blocks
	from := lastBlockHeight - params.TransactionToLiveInterval
	if from < 0 {
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

func (sb *sandbox) shouldPanicForDuplicatedAddress() {
	//
	// Why we should panic here?
	//
	// Try to make a new item which already exists in store.
	//
	logger.Panic("Duplicated address")
}

func (sb *sandbox) shouldPanicForUnknownAddress() {
	//
	// Why we should panic here?
	//
	// We only update accounts or validators which we have them inside the sandbox.
	// We must either make a new one (i.e. `MakeNewAccount`) or get it from store (i.e. `Account`) in advance.
	//
	logger.Panic("Unknown address")
}

func (sb *sandbox) Clear() {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	sb.clear()
}

func (sb *sandbox) clear() {
	sb.accounts = make(map[crypto.Address]*accountStatus)
	sb.validators = make(map[crypto.Address]*validatorStatus)
	sb.totalAccounts = sb.store.TotalAccounts()
	sb.totalValidators = sb.store.TotalValidators()
}

func (sb *sandbox) CommitChanges(round int) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	joined := make([]*validator.Validator, 0)
	for _, val := range sb.validators {
		if val.addToSet {
			joined = append(joined, &val.validator)
		}
	}
	if err := sb.validatorSet.MoveToNextHeight(0, joined); err != nil {
		//
		// We should panic here before modifying state store
		//
		logger.Panic("An error occurred", "err", err)
	}

	for _, acc := range sb.accounts {
		if acc.updated {
			sb.store.UpdateAccount(&acc.account)
		}
	}

	for _, val := range sb.validators {
		if val.updated {
			sb.store.UpdateValidator(&val.validator)
		}
	}
}

func (sb *sandbox) Account(addr crypto.Address) *account.Account {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.accounts[addr]
	if ok {
		copy := new(account.Account)
		*copy = s.account
		return copy
	}

	acc, err := sb.store.Account(addr)
	if err != nil {
		return nil
	}
	sb.accounts[addr] = &accountStatus{
		account: *acc,
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
	sb.accounts[addr] = &accountStatus{
		account: *acc,
		updated: true,
	}
	sb.totalAccounts++
	return acc
}

func (sb *sandbox) UpdateAccount(acc *account.Account) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := acc.Address()
	s, ok := sb.accounts[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}
	s.account = *acc
	s.updated = true
}

func (sb *sandbox) Validator(addr crypto.Address) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if ok {
		copy := new(validator.Validator)
		*copy = s.validator
		return copy
	}

	val, err := sb.store.Validator(addr)
	if err != nil {
		return nil
	}
	sb.validators[addr] = &validatorStatus{
		validator: *val,
	}
	return val
}

func (sb *sandbox) MakeNewValidator(pub crypto.PublicKey) *validator.Validator {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	addr := pub.Address()
	if sb.store.HasValidator(addr) {
		sb.shouldPanicForDuplicatedAddress()
	}

	val := validator.NewValidator(pub, sb.totalAccounts, sb.lastHeight()+1)
	sb.validators[addr] = &validatorStatus{
		validator: *val,
		updated:   true,
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

	sb.changeToStake += val.Stake() - s.validator.Stake()
	s.validator = *val
	s.updated = true
}

func (sb *sandbox) AddToSet(blockHash crypto.Hash, addr crypto.Address) error {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	s, ok := sb.validators[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}

	if sb.validatorSet.Contains(addr) {
		return errors.Errorf(errors.ErrGeneric, "This validator already is in the set")
	}

	joined := 0
	for _, s := range sb.validators {
		if s.addToSet {
			joined++
		}
	}
	if joined >= (sb.validatorSet.MaximumPower() / 3) {
		return errors.Errorf(errors.ErrGeneric, "In each height only 1/3 of validator can be changed")
	}
	h, _ := sb.store.BlockHeight(blockHash)
	b, err := sb.store.Block(h + 1)
	if err != nil {
		return errors.Errorf(errors.ErrGeneric, "Invalid block hash")
	}
	commiters := b.LastCommit().Committers()
	for _, c := range commiters {
		if c.Address.EqualsTo(addr) {
			return errors.Errorf(errors.ErrGeneric, "This validator was in the set in time of doing sortition")
		}
	}

	s.addToSet = true
	return nil
}

func (sb *sandbox) MaxMemoLength() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.params.MaximumMemoLength
}

func (sb *sandbox) FeeFraction() float64 {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.params.FeeFraction
}

func (sb *sandbox) MinFee() int64 {
	return sb.params.MinimumFee
}

func (sb *sandbox) TransactionToLiveInterval() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.TransactionToLiveInterval
}

func (sb *sandbox) RecentBlockHeight(hash crypto.Hash) int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	h, has := sb.recentBlocks.Get(hash)
	if !has {
		return -1
	}

	return h.(int)
}

func (sb *sandbox) lastHeight() int {
	_, v := sb.recentBlocks.Last()
	return v.(int) + 1
}

func (sb *sandbox) CurrentHeight() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.lastHeight() + 1
}

func (sb *sandbox) LastBlockHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.lastHeight()
}

func (sb *sandbox) LastBlockHash() crypto.Hash {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	k, _ := sb.recentBlocks.Last()
	return k.(crypto.Hash)
}

func (sb *sandbox) AppendNewBlock(hash crypto.Hash, height int) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	sb.recentBlocks.PushBack(hash, height)
}

func (sb *sandbox) VerifySortition(blockHash crypto.Hash, proof []byte, val *validator.Validator) bool {
	return sb.sortition.VerifySortition(blockHash, proof, val)
}
