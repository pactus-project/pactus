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

	val := sb.Validator(pld.Address)
	if val == nil {
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve validator")
	}

	if sb.CurrentHeight()-val.LastBondingHeight() < sb.Params().BondInterval {
		return errors.Errorf(errors.ErrInvalidHeight,
			"validator has bonded at height %v", val.LastBondingHeight())
	}
	// Power for parked validators (unbonded) set to zero.
	// So the proof is not valid, even if they have enough stake.
	ok := sb.VerifyProof(trx.Stamp(), pld.Proof, val)
	if !ok {
		return errors.Error(errors.ErrInvalidProof)
	}
	if e.strict {
		// A validator might produce more than one sortition transaction
		// before entering into the committee
		// In non-strict mode we don't check the sequence number
		if val.Sequence()+1 != trx.Sequence() {
			return errors.Errorf(errors.ErrInvalidSequence,
				"expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
		}
		if sb.Committee().Size() >= sb.Params().CommitteeSize {
			if err := e.joinCommittee(sb, val); err != nil {
				return err
			}
		} else {
			// Committee has free seats.
			// Rejecting sortition transactions from existing committee members.
			if sb.Committee().Contains(val.Address()) {
				return errors.Errorf(errors.ErrInvalidTx,
					"validator is in committee")
			}
		}
	}

	val.IncSequence()
	val.UpdateLastJoinedHeight(sb.CurrentHeight())

	sb.UpdateValidator(val)

	return nil
}

func (e *SortitionExecutor) Fee() int64 {
	return 0
}

func (e *SortitionExecutor) joinCommittee(sb sandbox.Sandbox,
	val *validator.Validator) error {
	joiningNum := 0
	joiningPower := int64(0)
	committee := sb.Committee()
	currentHeight := sb.CurrentHeight()
	sb.IterateValidators(func(vs *sandbox.ValidatorStatus) {
		if vs.Validator.LastJoinedHeight() == currentHeight {
			if !committee.Contains(vs.Validator.Address()) {
				joiningPower += vs.Validator.Power()
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
		return vals[i].LastJoinedHeight() < vals[j].LastJoinedHeight()
	})
	leavingPower := int64(0)
	for i := 0; i < joiningNum; i++ {
		leavingPower += vals[i].Power()
	}
	if leavingPower >= (committee.TotalPower() / 3) {
		return errors.Errorf(errors.ErrInvalidTx,
			"in each height only 1/3 of stake can leave")
	}

	oldestJoinedHeight := currentHeight
	for _, v := range committee.Validators() {
		if v.LastJoinedHeight() < oldestJoinedHeight {
			oldestJoinedHeight = v.LastJoinedHeight()
		}
	}

	// If the oldest validator in the committee still hasn't propose a block yet,
	// she stays in the committee.
	// We assumes all blocks has committed in round 0, in future we can consider
	// round parameter. It is backward compatible
	if currentHeight-oldestJoinedHeight < uint32(sb.Params().CommitteeSize) {
		return errors.Errorf(errors.ErrInvalidTx,
			"oldest validator still didn't propose any block")
	}
	return nil
}
