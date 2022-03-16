package executor

import (
	"sort"

	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
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
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve validator")
	}

	if sb.CurrentHeight()-val.LastBondingHeight() < sb.BondInterval() {
		return errors.Errorf(errors.ErrInvalidTx, "Validator has bonded at height %v", val.LastBondingHeight())
	}
	// Power for parked validators (unbonded) set to zero.
	// So the proof is not valid, even they have enough stake.
	ok := sb.VerifyProof(trx.Stamp(), pld.Proof, val)
	if !ok {
		return errors.Errorf(errors.ErrInvalidTx, "Sortition proof is invalid")
	}
	if e.strict {
		// A validator might produce more than one sortition transaction before entring into the committee
		// In non-strict mode we don't check the sequence number
		if val.Sequence()+1 != trx.Sequence() {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
		}
		if sb.Committee().Size() >= sb.CommitteeSize() {
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
				return errors.Errorf(errors.ErrGeneric, "in each height only 1/3 of stake can join")
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
				return errors.Errorf(errors.ErrGeneric, "in each height only 1/3 of stake can leave")
			}

			oldestJoinedHeight := currentHeight
			for _, v := range committee.Validators() {
				if v.LastJoinedHeight() < oldestJoinedHeight {
					oldestJoinedHeight = v.LastJoinedHeight()
				}
			}
			if currentHeight-oldestJoinedHeight < sb.CommitteeSize() {
				return errors.Errorf(errors.ErrGeneric, "oldest validator still didn't propose any block")
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
