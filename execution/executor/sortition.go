package executor

import (
	"cmp"
	"slices"

	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
)

type SortitionExecutor struct {
	sbx             sandbox.Sandbox
	pld             *payload.SortitionPayload
	validator       *validator.Validator
	sortitionHeight types.Height
}

func newSortitionExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*SortitionExecutor, error) {
	pld := trx.Payload().(*payload.SortitionPayload)

	val := sbx.Validator(pld.Validator)
	if val == nil {
		return nil, ValidatorNotFoundError{
			Address: pld.Validator,
		}
	}

	return &SortitionExecutor{
		pld:             pld,
		sbx:             sbx,
		validator:       val,
		sortitionHeight: trx.LockTime(),
	}, nil
}

func (e *SortitionExecutor) Check(strict bool) error {
	if e.sbx.CurrentHeight().SafeSub(e.validator.LastBondingHeight()) < e.sbx.Params().BondInterval {
		return ErrBondingPeriod
	}

	ok := e.sbx.VerifyProof(e.sortitionHeight, e.pld.Proof, e.validator)
	if !ok {
		return ErrInvalidSortitionProof
	}

	// Check for the duplicated or expired sortition transactions
	if e.sortitionHeight <= e.validator.LastSortitionHeight() {
		return ErrExpiredSortition
	}

	if strict {
		if err := e.canJoinCommittee(); err != nil {
			return err
		}
	}

	return nil
}

func (e *SortitionExecutor) canJoinCommittee() error {
	if e.sbx.Committee().Size() < e.sbx.Params().CommitteeSize {
		// There are available seats in the committee.
		if e.sbx.Committee().Contains(e.pld.Validator) {
			return ErrValidatorInCommittee
		}

		return nil
	}

	// The committee is full, check if the validator can join the committee.
	joiningNum := 0
	joiningPower := int64(0)
	committee := e.sbx.Committee()
	e.sbx.IterateValidators(func(val *validator.Validator, _ bool, joined bool) {
		if joined {
			if !committee.Contains(val.Address()) {
				joiningPower += val.Power()
				joiningNum++
			}
		}
	})
	if !committee.Contains(e.pld.Validator) {
		joiningPower += e.validator.Power()
		joiningNum++
	}
	if joiningPower >= (committee.TotalPower() / 3) {
		return ErrCommitteeJoinLimitExceeded
	}

	vals := committee.Validators()
	slices.SortStableFunc(vals, func(a, b *validator.Validator) int {
		return cmp.Compare(a.LastSortitionHeight(), b.LastSortitionHeight())
	})
	leavingPower := int64(0)
	// The number of leaving validators is the same as the number of joining validators,
	// and the leaving validators are the ones with the oldest sortition height.
	for i := 0; i < joiningNum; i++ {
		leavingPower += vals[i].Power()
	}
	if leavingPower >= (committee.TotalPower() / 3) {
		return ErrCommitteeLeaveLimitExceeded
	}

	oldestSortitionHeight := e.sbx.CurrentHeight()
	for _, v := range committee.Validators() {
		if v.LastSortitionHeight() < oldestSortitionHeight {
			oldestSortitionHeight = v.LastSortitionHeight()
		}
	}

	// If the oldest validator in the committee still hasn't propose a block yet,
	// it stays in the committee.
	proposerHeight := e.sbx.Committee().Proposer(0).LastSortitionHeight()
	if oldestSortitionHeight >= proposerHeight {
		return ErrOldestValidatorNotProposed
	}

	return nil
}

func (e *SortitionExecutor) Execute() {
	e.validator.UpdateLastSortitionHeight(e.sortitionHeight)

	e.sbx.JoinedToCommittee(e.pld.Validator)
	e.sbx.UpdateValidator(e.validator)
}
