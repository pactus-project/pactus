package validator

import (
	"github.com/sasha-s/go-deadlock"
	"gitlab.com/zarb-chain/zarb-go/crypto"
)

type ValidatorSet struct {
	lk deadlock.RWMutex

	maximumPower  int
	validators    []*Validator
	proposerIndex int
	joined        int
}

func NewValidatorSet(validators []*Validator, maximumPower int) *ValidatorSet {
	validators2 := make([]*Validator, len(validators))
	copy(validators2, validators)
	set := &ValidatorSet{
		maximumPower:  maximumPower,
		validators:    validators2,
		proposerIndex: 0,
	}
	return set
}

// TotalPower equals to the number of validator in the set
func (set *ValidatorSet) TotalPower() int {
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

func (set *ValidatorSet) Join(val *Validator) error {

	panic("Not supported yet")

	return nil
}

func (set *ValidatorSet) ForceLeave(val *Validator) error {
	set.lk.Lock()
	defer set.lk.Unlock()

	// Slashing validators should be supported
	panic("Not supported yet")

	return nil
}

func (set *ValidatorSet) MoveProposer(round int) {
	set.proposerIndex = (set.proposerIndex + round + 1) % len(set.validators)

}

func (set *ValidatorSet) Contains(addr crypto.Address) bool {
	for _, v := range set.validators {
		if v.Address().EqualsTo(addr) {
			return true
		}
	}
	return false
}

func (set *ValidatorSet) Validator(addr crypto.Address) *Validator {
	for _, v := range set.validators {
		if v.Address().EqualsTo(addr) {
			return v
		}
	}
	return nil
}

func (set *ValidatorSet) Proposer(round int) *Validator {
	idx := (set.proposerIndex + round) % len(set.validators)
	return set.validators[idx]
}
