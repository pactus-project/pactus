package validator

import (
	"fmt"

	"github.com/zarbchain/zarb-go/errors"

	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
)

type ValidatorSet struct {
	lk deadlock.RWMutex

	maximumPower  int
	validators    []*Validator
	joined        []*Validator
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
		joined:        make([]*Validator, 0),
		proposerIndex: index,
	}, nil
}

// TotalPower equals to the number of validator in the set
func (set *ValidatorSet) TotalPower() int {
	set.lk.Lock()
	defer set.lk.Unlock()

	return len(set.validators)
}

func (set *ValidatorSet) UpdateMaximumPower(maximumPower int) {
	set.lk.Lock()
	defer set.lk.Unlock()

	panic("Not supported yet")
	set.maximumPower = maximumPower
}

func (set *ValidatorSet) MaximumPower() int {
	return set.maximumPower
}

func (set *ValidatorSet) Power() int {
	return len(set.validators)
}

func (set *ValidatorSet) Join(val *Validator) error {
	set.lk.Lock()
	defer set.lk.Unlock()

	if set.contains(val.Address()) {
		return errors.Errorf(errors.ErrGeneric, "Validator already is in the set")
	}
	if len(set.joined) >= (set.Power() / 3) {
		return errors.Errorf(errors.ErrGeneric, "In each height only 1/3 of validator can be changed")
	}
	set.joined = append(set.joined, val)

	return nil
}

func (set *ValidatorSet) MoveToNewHeight(lastRound int) {
	set.lk.Lock()
	defer set.lk.Unlock()

	// Firts update proposer index
	set.proposerIndex = (set.proposerIndex + lastRound + 1) % len(set.validators)

	set.validators = append(set.validators, set.joined...)
	if set.Power() > set.MaximumPower() {
		shouldLeave := set.Power() - set.MaximumPower()
		set.validators = set.validators[shouldLeave:]
	}
	// Move proposer index after modifying the set
	set.proposerIndex = set.proposerIndex - len(set.joined)
	if set.proposerIndex < 0 {
		set.proposerIndex = 0
	}
	set.joined = set.joined[:0]
}

func (set *ValidatorSet) Validators() []crypto.Address {
	set.lk.Lock()
	defer set.lk.Unlock()

	vals := make([]crypto.Address, len(set.validators))
	for i, v := range set.validators {
		vals[i] = v.Address()
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

func (set *ValidatorSet) Proposer(round int) *Validator {
	set.lk.Lock()
	defer set.lk.Unlock()

	idx := (set.proposerIndex + round) % len(set.validators)
	return set.validators[idx]
}

func (set *ValidatorSet) CommittersHash() crypto.Hash {
	set.lk.Lock()
	defer set.lk.Unlock()

	data := make([][]byte, len(set.validators))

	for i, v := range set.validators {
		data[i] = make([]byte, 20)
		copy(data[i], v.Address().RawBytes())
	}
	merkle := simpleMerkle.NewTreeFromSlices(data)

	return merkle.Root()
}

// ---------
// For tests
func GenerateTestValidatorSet() (*ValidatorSet, []crypto.PrivateKey) {
	val1, pv1 := GenerateTestValidator()
	val2, pv2 := GenerateTestValidator()
	val3, pv3 := GenerateTestValidator()
	val4, pv4 := GenerateTestValidator()

	keys := []crypto.PrivateKey{pv1, pv2, pv3, pv4}
	vals := []*Validator{val1, val2, val3, val4}
	valset, _ := NewValidatorSet(vals, 4, val1.Address())
	return valset, keys
}
