package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type BondExecutor struct {
	fee    int64
	strict bool
}

func NewBondExecutor(strict bool) *BondExecutor {
	return &BondExecutor{strict: strict}
}

func (e *BondExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.BondPayload)

	bonderAcc := sb.Account(pld.Bonder)
	if bonderAcc == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve bonder account")
	}
	if bonderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence. Expected: %v, got: %v", bonderAcc.Sequence()+1, trx.Sequence())
	}
	addr := pld.PublicKey.Address()
	val := sb.Validator(addr)
	if val == nil {
		val = sb.MakeNewValidator(pld.PublicKey)
	}
	if val.UnbondingHeight() > 0 {
		return errors.Errorf(errors.ErrInvalidTx, "Validator has unbonded at height %v", val.UnbondingHeight())
	}
	if e.strict {
		// In strict mode, bond transaction will be rejected if a validator is in committee.
		// In non-strict mode, we accept it and keep it inside tx pool to process it in next blocks
		if sb.IsInCommittee(addr) {
			return errors.Errorf(errors.ErrInvalidTx, "Validator %v is in committee", addr)
		}

		// In strict mode, a validator can not evaluate sortition during bonding perion.
		// In non-strict mode, we accept it and keep it inside tx pool to process it in next blocks
		if val.LastJoinedHeight() == sb.CurrentHeight() {
			return errors.Errorf(errors.ErrInvalidTx, "Validator %v will join committee in the next height", addr)
		}
	}
	if bonderAcc.Balance() < pld.Stake+trx.Fee() {
		return errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}

	bonderAcc.IncSequence()
	bonderAcc.SubtractFromBalance(pld.Stake + trx.Fee())
	val.AddToStake(pld.Stake)
	val.UpdateLastBondingHeight(sb.CurrentHeight())

	sb.UpdateAccount(bonderAcc)
	sb.UpdateValidator(val)

	e.fee = trx.Fee()

	return nil
}

func (e *BondExecutor) Fee() int64 {
	return e.fee
}
