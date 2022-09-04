package executor

import (
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/types/tx/payload"
	"github.com/zarbchain/zarb-go/util/errors"
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
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve sender account")
	}
	if senderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidSequence,
			"expected: %v, got: %v", senderAcc.Sequence()+1, trx.Sequence())
	}
	val := sb.Validator(pld.Receiver)
	if val == nil {
		if pld.PublicKey == nil {
			return errors.Errorf(errors.ErrInvalidPublicKey,
				"public key is not set")
		}
		val = sb.MakeNewValidator(pld.PublicKey)
	} else {
		if pld.PublicKey != nil {
			return errors.Errorf(errors.ErrInvalidPublicKey,
				"public key set")
		}
	}
	if val.UnbondingHeight() > 0 {
		return errors.Errorf(errors.ErrInvalidHeight,
			"validator has unbonded at height %v", val.UnbondingHeight())
	}
	if e.strict {
		// In strict mode, bond transactions will be rejected if a validator is
		// in committee.
		// In non-strict mode, we accept it and keep it inside the tx pool to
		// process it when validator leaves the committee.
		if sb.Committee().Contains(pld.Receiver) {
			return errors.Errorf(errors.ErrInvalidTx,
				"validator %v is in committee", pld.Receiver)
		}

		// In strict mode, bond transactions will be rejected if a validator is
		// going to be in committee for the next height.
		// In non-strict mode, we accept it and keep it inside the tx pool to
		// process it when validator leaves the committee.
		if val.LastJoinedHeight() == sb.CurrentHeight() {
			return errors.Errorf(errors.ErrInvalidTx,
				"validator %v joins committee in the next height", pld.Receiver)
		}
	}
	if senderAcc.Balance() < pld.Stake+trx.Fee() {
		return errors.Error(errors.ErrInsufficientFunds)
	}
	if val.Stake()+pld.Stake > sb.Params().MaximumStake {
		return errors.Errorf(errors.ErrInvalidTx,
			"validator's stake can't be more than %v", sb.Params().MaximumStake)
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
