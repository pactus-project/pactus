package executor

import (
	"sort"

	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
	"github.com/pactus-project/pactus/util/errors"
)

type SortitionExecutor struct {
	strict bool
}

func NewSortitionExecutor(strict bool) *SortitionExecutor {
	return &SortitionExecutor{strict: strict}
}

func (e *SortitionExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.SortitionPayload)

	val := sb.Validator(pld.Validator)
	if val == nil {
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve validator")
	}

	if sb.CurrentHeight()-val.LastBondingHeight() < sb.Params().BondInterval {
		return errors.Errorf(errors.ErrInvalidHeight,
			"validator has bonded at height %v", val.LastBondingHeight())
	}

	sortitionHeight := trx.LockTime()
	ok := sb.VerifyProof(sortitionHeight, pld.Proof, val)
	if !ok {
		return errors.Error(errors.ErrInvalidProof)
	}

	// Check for the duplicated or expired sortition transactions
	if sortitionHeight <= val.LastSortitionHeight() {
		return errors.Errorf(errors.ErrInvalidTx,
			"duplicated sortition transaction")
	}

	if e.strict {
		if err := e.joinCommittee(sb, val); err != nil {
			return err
		}
	}

	val.UpdateLastSortitionHeight(sortitionHeight)

	sb.JoinedToCommittee(pld.Validator)
	sb.UpdateValidator(val)

	return nil
}

func (e *SortitionExecutor) joinCommittee(sb sandbox.Sandbox,
	val *validator.Validator,
) error {
	if sb.Committee().Size() < sb.Params().CommitteeSize {
		// There are available seats in the committee.
		if sb.Committee().Contains(val.Address()) {
			return errors.Errorf(errors.ErrInvalidTx,
				"validator is in committee")
		}

		return nil
	}

	// The committee is full, check if the validator can enter the committee.
	joiningNum := 0
	joiningPower := int64(0)
	committee := sb.Committee()
	sb.IterateValidators(func(val *validator.Validator, updated bool, joined bool) {
		if joined {
			if !committee.Contains(val.Address()) {
				joiningPower += val.Power()
				joiningNum++
			}
		}
	})
	if !committee.Contains(val.Address()) {
		joiningPower += val.Power()
		joiningNum++
	}
	if joiningPower >= (committee.TotalPower() / 3) {
		return errors.Errorf(errors.ErrInvalidTx,
			"in each height only 1/3 of stake can join")
	}

	vals := committee.Validators()
	sort.SliceStable(vals, func(i, j int) bool {
		return vals[i].LastSortitionHeight() < vals[j].LastSortitionHeight()
	})
	leavingPower := int64(0)
	for i := 0; i < joiningNum; i++ {
		leavingPower += vals[i].Power()
	}
	if leavingPower >= (committee.TotalPower() / 3) {
		return errors.Errorf(errors.ErrInvalidTx,
			"in each height only 1/3 of stake can leave")
	}

	oldestSortitionHeight := sb.CurrentHeight()
	for _, v := range committee.Validators() {
		if v.LastSortitionHeight() < oldestSortitionHeight {
			oldestSortitionHeight = v.LastSortitionHeight()
		}
	}

	// If the oldest validator in the committee still hasn't propose a block yet,
	// she stays in the committee.
	proposerHeight := sb.Committee().Proposer(0).LastSortitionHeight()
	if oldestSortitionHeight >= proposerHeight {
		return errors.Errorf(errors.ErrInvalidTx,
			"oldest validator still didn't propose any block")
	}

	return nil
}
