package sandbox

import (
	"sync"

	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/committee"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/logger"
	"github.com/zarbchain/zarb-go/param"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/store"
	"github.com/zarbchain/zarb-go/validator"
)

var _ Sandbox = &sandbox{}

type sandbox struct {
	lk sync.RWMutex

	store           store.Reader
	committee       committee.Reader
	accounts        map[crypto.Address]*AccountStatus
	validators      map[crypto.Address]*ValidatorStatus
	params          param.Params
	totalAccounts   int
	totalValidators int
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

func (sb *sandbox) UpdateAccount(acc *account.Account) {
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

// func (sb *sandbox) EnterCommittee(blockHash hash.Hash, addr crypto.Address) error {
// 	sb.lk.Lock()
// 	defer sb.lk.Unlock()

// 	valS, ok := sb.validators[addr]
// 	if !ok {
// 		return errors.Errorf(errors.ErrGeneric, "unknown validator")
// 	}

// 	_, bi := sb.latestBlocks.Last()
// 	if bi == nil {
// 		return errors.Errorf(errors.ErrGeneric, "Unable to retrieve last block info")
// 	}
// 	lastHeight := bi.(*BlockInfo).height

// 	if sb.committee.Size() >= sb.params.CommitteeSize {
// 		oldestJoinedHeight := lastHeight
// 		committeeStake := int64(0)
// 		for _, v := range sb.committee.Validators() {
// 			committeeStake += v.Stake()
// 			if v.LastJoinedHeight() < oldestJoinedHeight {
// 				oldestJoinedHeight = v.LastJoinedHeight()
// 			}
// 		}
// 		if lastHeight-oldestJoinedHeight < sb.params.CommitteeSize {
// 			return errors.Errorf(errors.ErrGeneric, "oldest validator still didn't propose any block")
// 		}
// 		joinedStake := int64(0)
// 		for _, s := range sb.validators {
// 			if s.JoinedCommittee {
// 				joinedStake += s.Validator.Stake()
// 			}
// 		}

// 		joinedStake += valS.Validator.Stake()
// 		if joinedStake >= (committeeStake / 3) {
// 			return errors.Errorf(errors.ErrGeneric, "in each height only 1/3 of stake can be changed")
// 		}
// 	}

// 	valS.JoinedCommittee = true
// 	return nil
// }

func (sb *sandbox) CommitteeAge() int {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	oldestJoinedHeight := sb.currentHeight()
	for _, v := range sb.committee.Validators() {
		if v.LastJoinedHeight() < oldestJoinedHeight {
			oldestJoinedHeight = v.LastJoinedHeight()
		}
	}

	return oldestJoinedHeight
}

func (sb *sandbox) CommitteePower() int64 {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.committee.TotalPower()
}

func (sb *sandbox) JoinedPower() int64 {
	joinedPower := int64(0)
	for _, s := range sb.validators {
		if s.Validator.LastJoinedHeight() == sb.currentHeight() {
			joinedPower += s.Validator.Power()
		}
	}

	return joinedPower
}

func (sb *sandbox) CommitteeHasFreeSeats() bool {
	sb.lk.Lock()
	defer sb.lk.Unlock()

	return sb.committee.Size() < sb.params.CommitteeSize
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

func (sb *sandbox) CurrentHeight() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.currentHeight()
}

func (sb *sandbox) currentHeight() int {
	h, _, err := sb.store.LastCertificate()
	if err != nil {
		return -1
	}

	return h + 1
}

func (sb *sandbox) PrevBlockHash() hash.Hash {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	_, cert, err := sb.store.LastCertificate()
	if err != nil {
		return cert.BlockHash()
	}

	return hash.UndefHash
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

func (sb *sandbox) BlockHeightByStamp(stamp hash.Stamp) int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.store.BlockHeightByStamp(stamp)
}

func (sb *sandbox) BlockSeedByStamp(stamp hash.Stamp) sortition.VerifiableSeed {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	height := sb.store.BlockHeightByStamp(stamp)
	b, err := sb.store.Block(height)
	if err != nil {
		return sortition.UndefVerifiableSeed
	}
	return b.Header().SortitionSeed()
}

func (sb *sandbox) CommitteeSize() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.CommitteeSize
}

func (sb *sandbox) UnbondInterval() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.UnbondInterval
}
func (sb *sandbox) BondInterval() int {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.params.BondInterval
}

func (sb *sandbox) IsInCommittee(addr crypto.Address) bool {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	return sb.committee.Contains(addr)
}

func (sb *sandbox) TotalPower() int64 {
	sb.lk.RLock()
	defer sb.lk.RUnlock()

	p := int64(0)
	sb.store.IterateValidators(func(val *validator.Validator) bool {
		p += val.Power()
		return false
	})
	return p
}
