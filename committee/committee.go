package committee

import (
	"fmt"
	"sort"

	"github.com/fxamacker/cbor/v2"
	"github.com/sasha-s/go-deadlock"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/validator"
)

var _ CommitteeReader = &Committee{}

type Committee struct {
	lk deadlock.RWMutex

	committeeSize int
	validators    []*validator.Validator
	proposerIndex int
}

func NewCommittee(validators []*validator.Validator, committeeSize int, proposer crypto.Address) (*Committee, error) {

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
	validators2 := make([]*validator.Validator, len(validators))
	copy(validators2, validators)
	return &Committee{
		committeeSize: committeeSize,
		validators:    validators2,
		proposerIndex: index,
	}, nil
}

func (committee *Committee) currentPower() int64 {
	p := int64(0)
	for _, v := range committee.validators {
		p += v.Power()
	}
	return p
}

func (committee *Committee) Update(lastRound int, joined []*validator.Validator) error {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	for _, v := range joined {
		if committee.contains(v.Address()) {
			return errors.Errorf(errors.ErrGeneric, "validator.Validator already is in the committee")
		}
	}

	if len(joined) > (committee.committeeSize / 3) {
		return errors.Errorf(errors.ErrGeneric, "In each update only 1/3 of validator committee can be changed")
	}

	sort.SliceStable(joined, func(i, j int) bool {
		return joined[i].Number() < joined[j].Number()
	})

	// First update proposer index
	committee.proposerIndex = (committee.proposerIndex + lastRound + 1) % len(committee.validators)

	committee.validators = append(committee.validators, joined...)
	if len(committee.validators) > committee.committeeSize {
		//
		shouldLeave := len(committee.validators) - committee.committeeSize
		committee.validators = committee.validators[shouldLeave:]
	}
	// Correcting proposer index
	committee.proposerIndex = committee.proposerIndex - len(joined)
	if committee.proposerIndex < 0 {
		committee.proposerIndex = 0
	}

	return nil
}

func (committee *Committee) CopyValidators() []*validator.Validator {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	vals := make([]*validator.Validator, len(committee.validators))
	for i, v := range committee.validators {
		vals[i] = v
	}
	return vals
}

func (committee *Committee) Contains(addr crypto.Address) bool {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	return committee.contains(addr)
}

func (committee *Committee) contains(addr crypto.Address) bool {
	for _, v := range committee.validators {
		if v.Address().EqualsTo(addr) {
			return true
		}
	}
	return false
}

func (committee *Committee) Validator(addr crypto.Address) *validator.Validator {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	for _, v := range committee.validators {
		if v.Address().EqualsTo(addr) {
			return v
		}
	}
	return nil
}

// IsProposer checks if the address is proposer for this run at the given round
func (committee *Committee) IsProposer(addr crypto.Address, round int) bool {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	idx := (committee.proposerIndex + round) % len(committee.validators)
	return committee.validators[idx].Address().EqualsTo(addr)
}

// Proposer returns proposer info for this run at the given round
func (committee *Committee) Proposer(round int) *validator.Validator {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	idx := (committee.proposerIndex + round) % len(committee.validators)
	return committee.validators[idx]
}

func (committee *Committee) Members() []int {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	return committee.members()
}

func (committee *Committee) members() []int {
	members := make([]int, len(committee.validators))
	for i, v := range committee.validators {
		members[i] = v.Number()
	}

	return members
}

func (committee *Committee) CommitteeHash() crypto.Hash {
	committee.lk.Lock()
	defer committee.lk.Unlock()

	bz, _ := cbor.Marshal(committee.members())
	return crypto.HashH(bz)
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
