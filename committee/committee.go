package committee

import (
	"container/list"
	"fmt"
	"sort"
	"sync"

	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
)

var _ Reader = &Committee{}

type Committee struct {
	lk sync.RWMutex

	committeeSize int
	validatorList *list.List
	proposerPos   *list.Element
}

func NewCommittee(validators []*validator.Validator, committeeSize int, proposerAddress crypto.Address) (*Committee, error) {
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

	return &Committee{
		committeeSize: committeeSize,
		validatorList: validatorList,
		proposerPos:   proposerPos,
	}, nil
}

func (committee *Committee) TotalStake() int64 {
	committee.lk.RLock()
	defer committee.lk.RUnlock()

	s := int64(0)
	committee.iterate(func(v *validator.Validator) (stop bool) {
		s += v.Stake()
		return false
	})
	return s
}

func (committee *Committee) TotalPower() int64 {
	committee.lk.RLock()
	defer committee.lk.RUnlock()

	p := int64(0)
	committee.iterate(func(v *validator.Validator) (stop bool) {
		p += v.Power()
		return false
	})
	return p
}

func (committee *Committee) Update(lastRound int, joined []*validator.Validator) error {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	for _, v := range joined {
		if committee.contains(v.Address()) {
			return errors.Errorf(errors.ErrGeneric, "Validator already is in the committee")
		}
	}

	sort.SliceStable(joined, func(i, j int) bool {
		return joined[i].Number() < joined[j].Number()
	})

	// First update validator list
	for _, val := range joined {
		committee.validatorList.InsertBefore(val, committee.proposerPos)
	}

	// Now adjust the list

	oldestFirst := make([]*list.Element, committee.validatorList.Len())
	i := 0
	for e := committee.validatorList.Front(); e != nil; e = e.Next() {
		oldestFirst[i] = e
		i++
	}

	sort.SliceStable(oldestFirst, func(i, j int) bool {
		return oldestFirst[i].Value.(*validator.Validator).LastJoinedHeight() < oldestFirst[j].Value.(*validator.Validator).LastJoinedHeight()
	})

	for i := 0; i <= lastRound; i++ {
		committee.proposerPos = committee.proposerPos.Next()
		if committee.proposerPos == nil {
			committee.proposerPos = committee.validatorList.Front()
		}
	}

	adjust := committee.validatorList.Len() - committee.committeeSize
	for i := 0; i < adjust; i++ {
		if oldestFirst[i] == committee.proposerPos {
			committee.proposerPos = committee.proposerPos.Next()
			if committee.proposerPos == nil {
				committee.proposerPos = committee.validatorList.Front()
			}
		}
		committee.validatorList.Remove(oldestFirst[i])
	}

	return nil
}

func (committee *Committee) Validators() []*validator.Validator {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	vals := make([]*validator.Validator, committee.validatorList.Len())
	i := 0
	committee.iterate(func(v *validator.Validator) (stop bool) {
		vals[i] = v
		i++
		return false
	})

	return vals
}

func (committee *Committee) Contains(addr crypto.Address) bool {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	return committee.contains(addr)
}

func (committee *Committee) contains(addr crypto.Address) bool {
	found := false
	committee.iterate(func(v *validator.Validator) (stop bool) {
		if v.Address().EqualsTo(addr) {
			found = true
			return true
		}
		return false
	})
	return found
}

func (committee *Committee) Validator(addr crypto.Address) *validator.Validator {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	var val *validator.Validator
	committee.iterate(func(v *validator.Validator) (stop bool) {
		if v.Address().EqualsTo(addr) {
			val = v
			return true
		}
		return false
	})
	return val
}

// IsProposer checks if the address is proposer for this run at the given round
func (committee *Committee) IsProposer(addr crypto.Address, round int) bool {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	p := committee.proposer(round)
	return p.Address().EqualsTo(addr)
}

// Proposer returns proposer info for this run at the given round
func (committee *Committee) Proposer(round int) *validator.Validator {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	return committee.proposer(round)
}

func (committee *Committee) proposer(round int) *validator.Validator {
	pos := committee.proposerPos
	for i := 0; i < round; i++ {
		pos = pos.Next()
		if pos == nil {
			pos = committee.validatorList.Front()
		}
	}

	return pos.Value.(*validator.Validator)
}

func (committee *Committee) Committers() []int {
	committee.lk.RLock()
	defer committee.lk.RUnlock()

	return committee.committers()
}

func (committee *Committee) Size() int {
	committee.lk.RLock()
	defer committee.lk.RUnlock()

	return committee.validatorList.Len()
}

func (committee *Committee) committers() []int {
	committers := make([]int, committee.validatorList.Len())
	i := 0
	committee.iterate(func(v *validator.Validator) (stop bool) {
		committers[i] = v.Number()
		i++
		return false
	})

	return committers
}

// iterate uses for easy iteration over validators in list
func (committee *Committee) iterate(consumer func(*validator.Validator) (stop bool)) {
	for e := committee.validatorList.Front(); e != nil; e = e.Next() {
		if consumer(e.Value.(*validator.Validator)) {
			return
		}
	}
}

// GenerateTestCommittee generates a validator committee for testing purpose
func GenerateTestCommittee() (*Committee, []crypto.Signer) {
	val1, s1 := validator.GenerateTestValidator(0)
	val2, s2 := validator.GenerateTestValidator(1)
	val3, s3 := validator.GenerateTestValidator(2)
	val4, s4 := validator.GenerateTestValidator(3)

	signers := []crypto.Signer{s1, s2, s3, s4}
	vals := []*validator.Validator{val1, val2, val3, val4}
	committee, _ := NewCommittee(vals, 4, val1.Address())
	return committee, signers
}
