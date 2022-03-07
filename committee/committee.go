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

var _ Committee = &committee{}

type committee struct {
	lk sync.RWMutex

	committeeSize int
	validatorList *list.List
	proposerPos   *list.Element
}

func NewCommittee(validators []*validator.Validator, committeeSize int, proposerAddress crypto.Address) (Committee, error) {
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

	return &committee{
		committeeSize: committeeSize,
		validatorList: validatorList,
		proposerPos:   proposerPos,
	}, nil
}

func (c *committee) TotalPower() int64 {
	c.lk.RLock()
	defer c.lk.RUnlock()

	p := int64(0)
	c.iterate(func(v *validator.Validator) (stop bool) {
		p += v.Power()
		return false
	})
	return p
}

func (c *committee) Update(lastRound int, joined []*validator.Validator) error {
	c.lk.Lock()
	defer c.lk.Unlock()

	for _, v := range joined {
		if c.contains(v.Address()) {
			return errors.Errorf(errors.ErrGeneric, "Validator is already in the committee")
		}
	}

	sort.SliceStable(joined, func(i, j int) bool {
		return joined[i].Number() < joined[j].Number()
	})

	// First update validator list
	for _, val := range joined {
		c.validatorList.InsertBefore(val, c.proposerPos)
	}

	// Now adjust the list

	oldestFirst := make([]*list.Element, c.validatorList.Len())
	i := 0
	for e := c.validatorList.Front(); e != nil; e = e.Next() {
		oldestFirst[i] = e
		i++
	}

	sort.SliceStable(oldestFirst, func(i, j int) bool {
		return oldestFirst[i].Value.(*validator.Validator).LastJoinedHeight() < oldestFirst[j].Value.(*validator.Validator).LastJoinedHeight()
	})

	for i := 0; i <= lastRound; i++ {
		c.proposerPos = c.proposerPos.Next()
		if c.proposerPos == nil {
			c.proposerPos = c.validatorList.Front()
		}
	}

	adjust := c.validatorList.Len() - c.committeeSize
	for i := 0; i < adjust; i++ {
		if oldestFirst[i] == c.proposerPos {
			c.proposerPos = c.proposerPos.Next()
			if c.proposerPos == nil {
				c.proposerPos = c.validatorList.Front()
			}
		}
		c.validatorList.Remove(oldestFirst[i])
	}

	return nil
}

func (c *committee) Validators() []*validator.Validator {
	c.lk.Lock()
	defer c.lk.Unlock()

	vals := make([]*validator.Validator, c.validatorList.Len())
	i := 0
	c.iterate(func(v *validator.Validator) (stop bool) {
		vals[i] = v
		i++
		return false
	})

	return vals
}

func (c *committee) Contains(addr crypto.Address) bool {
	c.lk.Lock()
	defer c.lk.Unlock()

	return c.contains(addr)
}

func (c *committee) contains(addr crypto.Address) bool {
	found := false
	c.iterate(func(v *validator.Validator) (stop bool) {
		if v.Address().EqualsTo(addr) {
			found = true
			return true
		}
		return false
	})
	return found
}

// func (c *committee) Validator(addr crypto.Address) *validator.Validator {
// 	c.lk.Lock()
// 	defer c.lk.Unlock()

// 	var val *validator.Validator
// 	c.iterate(func(v *validator.Validator) (stop bool) {
// 		if v.Address().EqualsTo(addr) {
// 			val = v
// 			return true
// 		}
// 		return false
// 	})
// 	return val
// }

// IsProposer checks if the address is proposer for this run at the given round
func (c *committee) IsProposer(addr crypto.Address, round int) bool {
	c.lk.Lock()
	defer c.lk.Unlock()

	p := c.proposer(round)
	return p.Address().EqualsTo(addr)
}

// Proposer returns proposer info for this run at the given round
func (c *committee) Proposer(round int) *validator.Validator {
	c.lk.Lock()
	defer c.lk.Unlock()

	return c.proposer(round)
}

func (c *committee) proposer(round int) *validator.Validator {
	pos := c.proposerPos
	for i := 0; i < round; i++ {
		pos = pos.Next()
		if pos == nil {
			pos = c.validatorList.Front()
		}
	}

	return pos.Value.(*validator.Validator)
}

func (c *committee) Committers() []int {
	c.lk.RLock()
	defer c.lk.RUnlock()

	return c.committers()
}

func (c *committee) Size() int {
	c.lk.RLock()
	defer c.lk.RUnlock()

	return c.validatorList.Len()
}

func (c *committee) committers() []int {
	committers := make([]int, c.validatorList.Len())
	i := 0
	c.iterate(func(v *validator.Validator) (stop bool) {
		committers[i] = v.Number()
		i++
		return false
	})

	return committers
}

// iterate uses for easy iteration over validators in list
func (c *committee) iterate(consumer func(*validator.Validator) (stop bool)) {
	for e := c.validatorList.Front(); e != nil; e = e.Next() {
		if consumer(e.Value.(*validator.Validator)) {
			return
		}
	}
}

// GenerateTestCommittee generates a validator committee for testing purpose
func GenerateTestCommittee() (Committee, []crypto.Signer) {
	val1, s1 := validator.GenerateTestValidator(0)
	val2, s2 := validator.GenerateTestValidator(1)
	val3, s3 := validator.GenerateTestValidator(2)
	val4, s4 := validator.GenerateTestValidator(3)

	signers := []crypto.Signer{s1, s2, s3, s4}
	vals := []*validator.Validator{val1, val2, val3, val4}
	committee, _ := NewCommittee(vals, 4, val1.Address())
	return committee, signers
}
