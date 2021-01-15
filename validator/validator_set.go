package validator

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

var _ ValidatorSetReader = &ValidatorSet{}

type ValidatorSetReader interface {
	MaximumPower() int
	CopyValidators() []*Validator
	Contains(addr crypto.Address) bool
	Proposer(round int) *Validator
	IsProposer(addr crypto.Address, round int) bool
	CommittersHash() crypto.Hash
}

type ValidatorSet struct {
	lk deadlock.RWMutex

	maximumPower  int
	validators    []*Validator
	proposerIndex int
}

func NewValidatorSet(validators []*Validator, maximumPower int, proposer crypto.Address) (*ValidatorSet, error) {

	index := -1
	for i, v := range validators {
		if v.Address().EqualsTo(proposer) {
			index = i
			break
		}
	}

	if index == -1 {
		return nil, fmt.Errorf("Proposer is not in the list")
	}
	validators2 := make([]*Validator, len(validators))
	copy(validators2, validators)
	return &ValidatorSet{
		maximumPower:  maximumPower,
		validators:    validators2,
		proposerIndex: index,
	}, nil
}

func (set *ValidatorSet) MaximumPower() int {
	return set.maximumPower
}

func (set *ValidatorSet) currentPower() int {
	p := 0
	for _, v := range set.validators {
		p += v.Power()
	}
	return p
}

func (set *ValidatorSet) UpdateTheSet(lastRound int, joined []*Validator) error {
	set.lk.Lock()
	defer set.lk.Unlock()

	for _, v := range joined {
		if set.contains(v.Address()) {
			return errors.Errorf(errors.ErrGeneric, "Validator already is in the set")
		}
	}

	if len(joined) > (set.MaximumPower() / 3) {
		return errors.Errorf(errors.ErrGeneric, "In each update only 1/3 of validator can be changed")
	}

	// First update proposer index
	set.proposerIndex = (set.proposerIndex + lastRound + 1) % len(set.validators)

	set.validators = append(set.validators, joined...)
	if set.currentPower() > set.maximumPower {
		shouldLeave := set.currentPower() - set.maximumPower
		set.validators = set.validators[shouldLeave:]
	}
	// Move proposer index after modifying the set
	set.proposerIndex = set.proposerIndex - len(joined)
	if set.proposerIndex < 0 {
		set.proposerIndex = 0
	}

	return nil
}

func (set *ValidatorSet) CopyValidators() []*Validator {
	set.lk.Lock()
	defer set.lk.Unlock()

	vals := make([]*Validator, len(set.validators))
	for i, v := range set.validators {
		vals[i] = v
	}
	return vals
}

func (set *ValidatorSet) Contains(addr crypto.Address) bool {
	set.lk.Lock()
	defer set.lk.Unlock()

	return set.contains(addr)
}

func (set *ValidatorSet) contains(addr crypto.Address) bool {
	for _, v := range set.validators {
		if v.Address().EqualsTo(addr) {
			return true
		}
	}
	return false
}

func (set *ValidatorSet) Validator(addr crypto.Address) *Validator {
	set.lk.Lock()
	defer set.lk.Unlock()

	for _, v := range set.validators {
		if v.Address().EqualsTo(addr) {
			return v
		}
	}
	return nil
}

// IsProposer checks if the address is proposer for this run at the given round
func (set *ValidatorSet) IsProposer(addr crypto.Address, round int) bool {
	set.lk.Lock()
	defer set.lk.Unlock()

	idx := (set.proposerIndex + round) % len(set.validators)
	return set.validators[idx].Address().EqualsTo(addr)
}

// Proposer returns proposer info for this run at the given round
func (set *ValidatorSet) Proposer(round int) *Validator {
	set.lk.Lock()
	defer set.lk.Unlock()

	idx := (set.proposerIndex + round) % len(set.validators)
	return set.validators[idx]
}

func (set *ValidatorSet) CommittersHash() crypto.Hash {
	set.lk.Lock()
	defer set.lk.Unlock()

	nums := make([]int, len(set.validators))
	for i, v := range set.validators {
		nums[i] = v.Number()
	}

	bz, _ := cbor.Marshal(nums)
	return crypto.HashH(bz)
}

// GenerateTestValidatorSet generates a validator set for testing purpose
func GenerateTestValidatorSet() (*ValidatorSet, []crypto.PrivateKey) {
	val1, pv1 := GenerateTestValidator(0)
	val2, pv2 := GenerateTestValidator(1)
	val3, pv3 := GenerateTestValidator(2)
	val4, pv4 := GenerateTestValidator(3)

	keys := []crypto.PrivateKey{pv1, pv2, pv3, pv4}
	vals := []*Validator{val1, val2, val3, val4}
	valset, _ := NewValidatorSet(vals, 4, val1.Address())
	return valset, keys
}
