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

	senderAcc := sb.Account(pld.Sender)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidAddress, "unable to retrieve sender account")
	}
	if senderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidSequence, "expected: %v, got: %v", senderAcc.Sequence()+1, trx.Sequence())
	}
	addr := pld.PublicKey.Address()
	val := sb.Validator(addr)
	if val == nil {
		val = sb.MakeNewValidator(pld.PublicKey)
	}
	if val.UnbondingHeight() > 0 {
		return errors.Errorf(errors.ErrInvalidHeight, "validator has unbonded at height %v", val.UnbondingHeight())
	}
	if e.strict {
		// In strict mode, bond transaction will be rejected if a validator is in committee.
		// In non-strict mode, we accept it and keep it inside tx pool to process it later
		if sb.Committee().Contains(addr) {
			return errors.Errorf(errors.ErrInvalidTx, "validator %v is in committee", addr)
		}

		// In strict mode, a validator can not evaluate sortition during bonding perion.
		// In non-strict mode, we accept it and keep it inside tx pool to process it later
		if val.LastJoinedHeight() == sb.CurrentHeight() {
			return errors.Errorf(errors.ErrInvalidHeight, "validator %v will join committee in the next height", addr)
		}
	}
	if senderAcc.Balance() < pld.Stake+trx.Fee() {
		return errors.Error(errors.ErrInsufficientFunds)
	}

	senderAcc.IncSequence()
	senderAcc.SubtractFromBalance(pld.Stake + trx.Fee())
	val.AddToStake(pld.Stake)
	val.UpdateLastBondingHeight(sb.CurrentHeight())

	sb.UpdateAccount(senderAcc)
	sb.UpdateValidator(val)

	e.fee = trx.Fee()

	return nil
}

func (e *BondExecutor) Fee() int64 {
	return e.fee
}
