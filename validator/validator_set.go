package validator

import (
	"container/list"
	"fmt"
	"sort"

	"github.com/fxamacker/cbor/v2"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
)

var _ ValidatorSetReader = &ValidatorSet{}

type ValidatorSetReader interface {
	Validators() []*Validator
	Contains(addr crypto.Address) bool
	Proposer(round int) *Validator
	IsProposer(addr crypto.Address, round int) bool
	CommitteeHash() crypto.Hash
}

type ValidatorSet struct {
	lk deadlock.RWMutex

	committeeSize int
	validatorList *list.List
	proposerPos   *list.Element
}

func NewValidatorSet(validators []*Validator, committeeSize int, proposerAddress crypto.Address) (*ValidatorSet, error) {

	validatorList := list.New()
	var proposerPos *list.Element

	for _, v := range validators {
		el := validatorList.PushBack(v)
		if v.Address().EqualsTo(proposerAddress) {
			proposerPos = el
		}
	}

	if proposerPos == nil {
		return nil, fmt.Errorf("Proposer is not in the list")
	}

	return &ValidatorSet{
		committeeSize: committeeSize,
		validatorList: validatorList,
		proposerPos:   proposerPos,
	}, nil
}

func (set *ValidatorSet) currentPower() int64 {
	p := int64(0)
	set.iterate(func(v *Validator) (stop bool) {
		p += v.Power()
		return false
	})
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

	if len(joined) > (set.committeeSize / 3) {
		return errors.Errorf(errors.ErrGeneric, "In each update only 1/3 of validator set can be changed")
	}

	sort.SliceStable(joined, func(i, j int) bool {
		return joined[i].Number() < joined[j].Number()
	})

	// First update validator list
	for _, val := range joined {
		set.validatorList.InsertBefore(val, set.proposerPos)
	}

	// Now adjust the list

	oldestFirst := make([]*list.Element, set.validatorList.Len())
	i := 0
	for e := set.validatorList.Front(); e != nil; e = e.Next() {
		oldestFirst[i] = e
		i++
	}

	sort.SliceStable(oldestFirst, func(i, j int) bool {
		return oldestFirst[i].Value.(*Validator).LastJoinedHeight() < oldestFirst[j].Value.(*Validator).LastJoinedHeight()
	})

	for i := 0; i <= lastRound; i++ {
		set.proposerPos = set.proposerPos.Next()
		if set.proposerPos == nil {
			set.proposerPos = set.validatorList.Front()
		}
	}

	adjust := set.validatorList.Len() - set.committeeSize
	for i := 0; i < adjust; i++ {
		if oldestFirst[i] == set.proposerPos {
			set.proposerPos = set.proposerPos.Next()
			if set.proposerPos == nil {
				set.proposerPos = set.validatorList.Front()
			}
		}
		set.validatorList.Remove(oldestFirst[i])
	}

	return nil
}

func (set *ValidatorSet) Validators() []*Validator {
	set.lk.Lock()
	defer set.lk.Unlock()

	vals := make([]*Validator, set.validatorList.Len())
	i := 0
	set.iterate(func(v *Validator) (stop bool) {
		vals[i] = v
		i++
		return false
	})

	return vals
}

func (set *ValidatorSet) Contains(addr crypto.Address) bool {
	set.lk.Lock()
	defer set.lk.Unlock()

	return set.contains(addr)
}

func (set *ValidatorSet) contains(addr crypto.Address) bool {
	found := false
	set.iterate(func(v *Validator) (stop bool) {
		if v.Address().EqualsTo(addr) {
			found = true
			return true
		}
		return false
	})
	return found
}

func (set *ValidatorSet) Validator(addr crypto.Address) *Validator {
	set.lk.Lock()
	defer set.lk.Unlock()

	var val *Validator
	set.iterate(func(v *Validator) (stop bool) {
		if v.Address().EqualsTo(addr) {
			val = v
			return true
		}
		return false
	})
	return val
}

// IsProposer checks if the address is proposer for this run at the given round
func (set *ValidatorSet) IsProposer(addr crypto.Address, round int) bool {
	set.lk.Lock()
	defer set.lk.Unlock()

	p := set.proposer(round)
	return p.Address().EqualsTo(addr)
}

// Proposer returns proposer info for this run at the given round
func (set *ValidatorSet) Proposer(round int) *Validator {
	set.lk.Lock()
	defer set.lk.Unlock()

	return set.proposer(round)
}

func (set *ValidatorSet) proposer(round int) *Validator {
	pos := set.proposerPos
	for i := 0; i < round; i++ {
		pos = pos.Next()
		if pos == nil {
			pos = set.validatorList.Front()
		}
	}

	return pos.Value.(*Validator)
}

func (set *ValidatorSet) Committee() []int {
	set.lk.Lock()
	defer set.lk.Unlock()

	return set.committee()
}

func (set *ValidatorSet) committee() []int {
	committee := make([]int, set.validatorList.Len())
	i := 0
	set.iterate(func(v *Validator) (stop bool) {
		committee[i] = v.Number()
		i++
		return false
	})

	return committee
}

func (set *ValidatorSet) CommitteeHash() crypto.Hash {
	set.lk.Lock()
	defer set.lk.Unlock()

	bz, _ := cbor.Marshal(set.committee())
	return crypto.HashH(bz)
}

// iterate uses for easy iteration over validators in list
func (set *ValidatorSet) iterate(consumer func(*Validator) (stop bool)) {
	for e := set.validatorList.Front(); e != nil; e = e.Next() {
		if consumer(e.Value.(*Validator)) {
			return
		}
	}
}

// GenerateTestValidatorSet generates a validator set for testing purpose
func GenerateTestValidatorSet() (*ValidatorSet, []crypto.Signer) {
	val1, s1 := GenerateTestValidator(0)
	val2, s2 := GenerateTestValidator(1)
	val3, s3 := GenerateTestValidator(2)
	val4, s4 := GenerateTestValidator(3)

	signers := []crypto.Signer{s1, s2, s3, s4}
	vals := []*Validator{val1, val2, val3, val4}
	valset, _ := NewValidatorSet(vals, 4, val1.Address())
	return valset, signers
}
