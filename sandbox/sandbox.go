package sandbox

import (
	"sync"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/libs/linkedmap"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

type BlockInfo struct {
	height int
	hash   hash.Hash
}

func NewBlockInfo(height int, hash hash.Hash) *BlockInfo {
	return &BlockInfo{
		height: height,
		hash:   hash,
	}
}

type Concrete struct {
	lk sync.RWMutex

	store            store.Reader
	sortition        *sortition.Sortition
	latestBlocks     *linkedmap.LinkedMap
	committee        committee.Reader
	accounts         map[crypto.Address]*AccountStatus
	validators       map[crypto.Address]*ValidatorStatus
	joinedCommittee  *crypto.Address
	params           param.Params
	totalAccounts    int
	totalValidators  int
	totalStakeChange int64
}

type ValidatorStatus struct {
	Validator validator.Validator
	Updated   bool
}

type AccountStatus struct {
	Account account.Account
	Updated bool
}

func NewSandbox(store store.Reader, params param.Params, latestBlocks *linkedmap.LinkedMap, sortition *sortition.Sortition, committee committee.Reader) *Concrete {
	sb := &Concrete{
		store:     store,
		sortition: sortition,
		committee: committee,
		params:    params,
	}

	sb.accounts = make(map[crypto.Address]*AccountStatus)
	sb.validators = make(map[crypto.Address]*ValidatorStatus)
	sb.latestBlocks = latestBlocks
	sb.totalAccounts = sb.store.TotalAccounts()
	sb.totalValidators = sb.store.TotalValidators()
	sb.totalStakeChange = 0
	sb.joinedCommittee = nil

	return sb
}

func (sb *Concrete) shouldPanicForDuplicatedAddress() {
	//
	// Why we should panic here?
	//
	// Try to make a new item which already exists in store.
	//
	logger.Panic("duplicated address")
}

func (sb *Concrete) shouldPanicForUnknownAddress() {
	//
	// Why we should panic here?
	//
	// We only update accounts or validators which we have them inside the sandbox.
	// We must either make a new one (i.e. `MakeNewAccount`) or get it from store (i.e. `Account`) in advance.
	//
	logger.Panic("unknown address")
}

func (sb *Concrete) Account(addr crypto.Address) *account.Account {
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
func (sb *Concrete) MakeNewAccount(addr crypto.Address) *account.Account {
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

func (sb *Concrete) UpdateAccount(acc *account.Account) {
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

func (sb *Concrete) Validator(addr crypto.Address) *validator.Validator {
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

func (sb *Concrete) MakeNewValidator(pub *bls.PublicKey) *validator.Validator {
	sb.lk.Lock()
	defer sb.lk.Unlock()

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

func (sb *Concrete) UpdateValidator(val *validator.Validator) {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	addr := val.Address()
	s, ok := sb.validators[addr]
	if !ok {
		sb.shouldPanicForUnknownAddress()
	}

	// shouldn't this be power??
	sb.totalStakeChange += val.Stake() - s.Validator.Stake()
	s.Validator = *val
	s.Updated = true
}

func (sb *Concrete) EnterCommittee(addr crypto.Address) error {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	valS, ok := sb.validators[addr]
	if !ok {
		return errors.Errorf(errors.ErrGeneric, "unknown validator")
	}

	if sb.joinedCommittee != nil {
		return errors.Errorf(errors.ErrGeneric, "a validator has joined into committee before")
	}

	_, lastBlockInfo := sb.latestBlocks.Last()
	if lastBlockInfo == nil {
		return errors.Errorf(errors.ErrGeneric, "Unable to retrieve last block info")
	}
	lastHeight := lastBlockInfo.(*BlockInfo).height

	if sb.committee.Size() >= sb.params.CommitteeSize {
		oldestJoinedHeight := lastHeight
		committeeStake := int64(0)
		for _, v := range sb.committee.Validators() {
			committeeStake += v.Stake()
			if v.LastJoinedHeight() < oldestJoinedHeight {
				oldestJoinedHeight = v.LastJoinedHeight()
			}
		}
		joinedStake := valS.Validator.Stake()
		if joinedStake >= (committeeStake / 3) {
			return errors.Errorf(errors.ErrGeneric, "in each height only 1/3 of stake can be changed")
		}
	}

	sb.joinedCommittee = &addr
	return nil
}

func (sb *Concrete) MaxMemoLength() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.params.MaximumMemoLength
}

func (sb *Concrete) FeeFraction() float64 {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.params.FeeFraction
}

func (sb *Concrete) MinFee() int64 {
	return sb.params.MinimumFee
}

func (sb *Concrete) TransactionToLiveInterval() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.TransactionToLiveInterval
}

func (sb *Concrete) BlockHeight(h hash.Hash) int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	if h.EqualsTo(hash.UndefHash) {
		return 0
	}

	height, err := sb.store.BlockHeight(h)
	if err != nil {
		return -1
	}

	return height
}

func (sb *Concrete) CurrentHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	_, val := sb.latestBlocks.Last()
	if val != nil {
		return val.(*BlockInfo).height + 1
	}

	return -1
}

func (sb *Concrete) PrevBlockHash() hash.Hash {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	_, val := sb.latestBlocks.Last()
	if val != nil {
		return val.(*BlockInfo).hash
	}

	return hash.UndefHash
}

func (sb *Concrete) VerifySortition(blockHash hash.Hash, proof sortition.Proof, val *validator.Validator) bool {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.sortition.VerifyProof(blockHash, proof, val.PublicKey(), val.Stake())
}

func (sb *Concrete) IterateAccounts(consumer func(*AccountStatus)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, as := range sb.accounts {
		consumer(as)
	}
}

func (sb *Concrete) IterateValidators(consumer func(*ValidatorStatus)) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	for _, vs := range sb.validators {
		consumer(vs)
	}
}

func (sb *Concrete) FindBlockInfoByStamp(stamp hash.Stamp) (int, hash.Hash) {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	val, ok := sb.latestBlocks.Get(stamp)
	if ok {
		bi := val.(*BlockInfo)
		return bi.height, bi.hash
	}

	return -1, hash.UndefHash
}

func (sb *Concrete) CommitteeSize() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.CommitteeSize
}

func (sb *Concrete) UnbondInterval() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.UnbondInterval
}

func (sb *Concrete) IsInCommittee(addr crypto.Address) bool {
	return sb.committee.Contains(addr)
}
