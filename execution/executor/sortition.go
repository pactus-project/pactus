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
	pld             *payload.SortitionPayload
	validator       *validator.Validator
	sortitionHeight types.Height
}

func newSortitionExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*SortitionExecutor, error) {
	pld := trx.Payload().(*payload.SortitionPayload)

	val := sbx.Validator(pld.Address)
	if val == nil {
		return nil, ValidatorNotFoundError{
			Address: pld.Address,
		}
	}

	return &SortitionExecutor{
		pld:             pld,
		validator:       val,
		sortitionHeight: trx.LockTime(),
	}, nil
}

func (e *SortitionExecutor) Check(sbx sandbox.SandboxReader, strict bool) error {
	if sbx.CurrentHeight().SafeSub(e.validator.LastBondingHeight()) < sbx.Params().BondInterval {
		return ErrBondingPeriod
	}

	ok := sbx.VerifyProof(e.sortitionHeight, e.pld.Proof, e.validator)
	if !ok {
		return ErrInvalidSortitionProof
	}

	// Check for the duplicated or expired sortition transactions
	if e.sortitionHeight <= e.validator.LastSortitionHeight() {
		return ErrExpiredSortition
	}

	if strict {
		if err := e.canJoinCommittee(sbx); err != nil {
			return err
		}
	}

	return nil
}

func (e *SortitionExecutor) canJoinCommittee(sbx sandbox.SandboxReader) error {
	if sbx.Committee().Size() < sbx.Params().CommitteeSize {
		// There are available seats in the committee.
		if sbx.Committee().Contains(e.pld.Address) {
			return ErrValidatorInCommittee
		}

		return nil
	}

	// The committee is full, check if the validator can join the committee.
	joiningNum := 0
	joiningPower := int64(0)
	committee := sbx.Committee()
	committeePower := committee.Power()
	committeeVals := committee.Validators()

	sbx.IterateValidators(func(val *validator.Validator, _ bool, joined bool) {
		if joined {
			if !committee.Contains(val.Address()) {
				joiningPower += val.Power()
				joiningNum++
			}
		}
	})
	if !committee.Contains(e.pld.Address) {
		joiningPower += e.validator.Power()
		joiningNum++
	}

	if joiningPower >= (committeePower / 3) {
		return ErrCommitteeJoinLimitExceeded
	}

	slices.SortStableFunc(committeeVals, func(a, b *validator.Validator) int {
		return cmp.Compare(a.LastSortitionHeight(), b.LastSortitionHeight())
	})
	leavingPower := int64(0)
	// The number of leaving validators is the same as the number of joining validators,
	// and the leaving validators are the ones with the oldest sortition height.
	for i := 0; i < joiningNum; i++ {
		leavingPower += committeeVals[i].Power()
	}
	if leavingPower >= (committeePower / 3) {
		return ErrCommitteeLeaveLimitExceeded
	}

	oldestSortitionHeight := sbx.CurrentHeight()
	for _, v := range committeeVals {
		if v.LastSortitionHeight() < oldestSortitionHeight {
			oldestSortitionHeight = v.LastSortitionHeight()
		}
	}

	// If the oldest validator in the committee still hasn't propose a block yet,
	// it stays in the committee.
	proposerHeight := sbx.Committee().Proposer(0).LastSortitionHeight()
	if oldestSortitionHeight >= proposerHeight {
		return ErrOldestValidatorNotProposed
	}

	return nil
}

func (e *SortitionExecutor) Execute(sbx sandbox.Sandbox) {
	e.validator.UpdateLastSortitionHeight(e.sortitionHeight)

	sbx.JoinToCommittee(e.pld.Address)
	sbx.UpdateValidator(e.validator)
}
