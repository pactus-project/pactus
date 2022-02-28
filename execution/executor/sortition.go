package executor

import (
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
	// Power for parked validators is set to zero
	if val.Power() == 0 {
		return errors.Errorf(errors.ErrInvalidTx, "Validator has no Power to be in committee")
	}
	if sb.CurrentHeight()-val.LastBondingHeight() < sb.BondInterval() {
		return errors.Errorf(errors.ErrInvalidTx, "In bonding period")
	}
	height, hash := sb.FindBlockInfoByStamp(trx.Stamp())
	ok := sb.VerifySortition(hash, pld.Proof, val)
	if !ok {
		return errors.Errorf(errors.ErrInvalidTx, "Sortition proof is invalid")
	}

	if val.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
	}

	if val.LastJoinedHeight() > 0 && height < val.LastJoinedHeight() {
		return errors.Errorf(errors.ErrInvalidTx, "Expired sortition. Last joined at: %v", val.LastJoinedHeight())
	}

	if sb.CommitteeHasFreeSeats() {
		committeeStake := sb.CommitteeStake()
		joinedStake := val.Stake()
		if joinedStake >= (committeeStake / 3) {
			return errors.Errorf(errors.ErrGeneric, "In each height less than 1/3 of stake can move")
		}
	}

	if e.strict {
		// There maybe more than one sortion transaction inside transaction pool
		if sb.HasAnyValidatorJoinedCommittee() {
			return errors.Errorf(errors.ErrGeneric, "a validator has joined into committee before")
		}
	}

	val.IncSequence()
	val.UpdateLastJoinedHeight(sb.CurrentHeight())

	sb.UpdateValidator(val)
	sb.JoinCommittee(val.Address())

	return nil
}

func (e *SortitionExecutor) Fee() int64 {
	return 0
}
