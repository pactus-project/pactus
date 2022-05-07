package committee

import (
	"container/list"
	"fmt"
	"sort"
	"sync"

	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/validator"
	"github.com/zarbchain/zarb-go/util"
)

var _ Committee = &committee{}

type committee struct {
	lk sync.RWMutex

	committeeSize int
	validatorList *list.List
	proposerPos   *list.Element
}

func NewCommittee(validators []*validator.Validator, committeeSize int,
	proposerAddress crypto.Address) (Committee, error) {
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

func (c *committee) Update(lastRound int16, joined []*validator.Validator) {
	c.lk.Lock()
	defer c.lk.Unlock()

	sort.SliceStable(joined, func(i, j int) bool {
		return joined[i].Number() < joined[j].Number()
	})

	// First update validator list
	for _, val := range joined {
		committeeVal := c.find(val.Address())
		if committeeVal == nil {
			c.validatorList.InsertBefore(val, c.proposerPos)
		} else {
			*committeeVal = *val
		}
	}

	// Now adjust the list
	oldestFirst := make([]*list.Element, c.validatorList.Len())
	i := 0
	for e := c.validatorList.Front(); e != nil; e = e.Next() {
		oldestFirst[i] = e
		i++
	}

	sort.SliceStable(oldestFirst, func(i, j int) bool {
		return oldestFirst[i].Value.(*validator.Validator).LastJoinedHeight() <
			oldestFirst[j].Value.(*validator.Validator).LastJoinedHeight()
	})

	for i := 0; i <= int(lastRound); i++ {
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

	return c.find(addr) != nil
}

func (c *committee) find(addr crypto.Address) *validator.Validator {
	var found *validator.Validator
	c.iterate(func(v *validator.Validator) (stop bool) {
		if v.Address().EqualsTo(addr) {
			found = v
			return true
		}
		return false
	})
	return found
}

// IsProposer checks if the address is proposer for this height at this round.
func (c *committee) IsProposer(addr crypto.Address, round int16) bool {
	c.lk.Lock()
	defer c.lk.Unlock()

	p := c.proposer(round)
	return p.Address().EqualsTo(addr)
}

// Proposer returns proposer info for this height at this round.
func (c *committee) Proposer(round int16) *validator.Validator {
	c.lk.Lock()
	defer c.lk.Unlock()

	return c.proposer(round)
}

func (c *committee) proposer(round int16) *validator.Validator {
	pos := c.proposerPos
	for i := 0; i < int(round); i++ {
		pos = pos.Next()
		if pos == nil {
			pos = c.validatorList.Front()
		}
	}

	return pos.Value.(*validator.Validator)
}

func (c *committee) Committers() []int32 {
	c.lk.RLock()
	defer c.lk.RUnlock()

	committers := make([]int32, c.validatorList.Len())
	i := 0
	c.iterate(func(v *validator.Validator) (stop bool) {
		committers[i] = v.Number()
		i++
		return false
	})

	return committers
}

func (c *committee) Size() int {
	c.lk.RLock()
	defer c.lk.RUnlock()

	return c.validatorList.Len()
}

// iterate uses for easy iteration over validators in list.
func (c *committee) iterate(consumer func(*validator.Validator) (stop bool)) {
	for e := c.validatorList.Front(); e != nil; e = e.Next() {
		if consumer(e.Value.(*validator.Validator)) {
			return
		}
	}
}

// GenerateTestCommittee generates a committee for testing purpose.
// All committee members have same power.
func GenerateTestCommittee(num int) (Committee, []crypto.Signer) {
	signers := make([]crypto.Signer, num)
	vals := make([]*validator.Validator, num)
	h1 := util.RandInt32(100000)
	for i := int32(0); i < int32(num); i++ {
		val, s := validator.GenerateTestValidator(i)
		signers[i] = s
		vals[i] = val

		val.UpdateLastBondingHeight(h1 + i)
		val.UpdateLastJoinedHeight(h1 + 1000 + i)
		//
		val.SubtractFromStake(val.Stake())
		val.AddToStake(1 * 10e8)
	}

	committee, _ := NewCommittee(vals, num, vals[0].Address())
	return committee, signers
}
