package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/sortition"
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
	if e.strict {
		// A validator might produce more than one sortition transaction before entring into the committee
		// In non-strict mode we don't check the sequence number
		if val.Sequence()+1 != trx.Sequence() {
			return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", val.Sequence()+1, trx.Sequence())
		}
	}
	//////// ????????
	// Power for parked validators (unbonded) set to zero
	// if val.Power() == 0 {
	// 	return errors.Errorf(errors.ErrInvalidTx, "Validator has no Power to be in committee")
	// }
	if sb.CurrentHeight()-val.LastBondingHeight() < sb.BondInterval() {
		return errors.Errorf(errors.ErrInvalidTx, "Validator has bonded at height %v", val.LastBondingHeight())
	}
	seed := sb.BlockSeedByStamp(trx.Stamp())
	ok := sortition.VerifyProof(seed, pld.Proof, val.PublicKey(), sb.TotalPower(), val.Power())
	if !ok {
		return errors.Errorf(errors.ErrInvalidTx, "Sortition proof is invalid")
	}
	if !sb.CommitteeHasFreeSeats() {
		if sb.CommitteeAge() < sb.CommitteeSize() {
			return errors.Errorf(errors.ErrGeneric, "oldest validator still didn't propose any block")
		}

		if sb.JoinedPower() >= (sb.CommitteePower() / 3) {
			return errors.Errorf(errors.ErrGeneric, "in each height only 1/3 of stake can be changed")
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
