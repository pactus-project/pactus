package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/errors"
)

type BondExecutor struct {
	strict bool
}

func NewBondExecutor(strict bool) *BondExecutor {
	return &BondExecutor{strict: strict}
}

func (e *BondExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.BondPayload)

	senderAcc := sb.Account(pld.From)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve sender account")
	}

	receiverVal := sb.Validator(pld.To)
	if receiverVal == nil {
		if pld.PublicKey == nil {
			return errors.Errorf(errors.ErrInvalidPublicKey,
				"public key is not set")
		}
		if pld.Stake < sb.Params().MinimumStake {
			return errors.Errorf(errors.ErrInvalidTx,
				"validator's stake can't be less than %v", sb.Params().MinimumStake)
		}
		receiverVal = sb.MakeNewValidator(pld.PublicKey)
	} else if pld.PublicKey != nil {
		return errors.Errorf(errors.ErrInvalidPublicKey,
			"public key is set")
	}
	if receiverVal.UnbondingHeight() > 0 {
		return errors.Errorf(errors.ErrInvalidHeight,
			"validator has unbonded at height %v", receiverVal.UnbondingHeight())
	}
	if e.strict {
		// In strict mode, bond transactions will be rejected if a validator is
		// already in the committee.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if sb.Committee().Contains(pld.To) {
			return errors.Errorf(errors.ErrInvalidTx,
				"validator %v is in committee", pld.To)
		}

		// In strict mode, bond transactions will be rejected if a validator is
		// going to join the committee in the next height.
		// In non-strict mode, they are added to the transaction pool and
		// processed once eligible.
		if sb.IsJoinedCommittee(pld.To) {
			return errors.Errorf(errors.ErrInvalidTx,
				"validator %v joins committee in the next height", pld.To)
		}
	}
	if senderAcc.Balance() < pld.Stake+trx.Fee() {
		return ErrInsufficientFunds
	}
	if receiverVal.Stake()+pld.Stake > sb.Params().MaximumStake {
		return errors.Errorf(errors.ErrInvalidAmount,
			"validator's stake can't be more than %v", sb.Params().MaximumStake)
	}

	senderAcc.SubtractFromBalance(pld.Stake + trx.Fee())
	receiverVal.AddToStake(pld.Stake)
	receiverVal.UpdateLastBondingHeight(sb.CurrentHeight())

	sb.UpdatePowerDelta(int64(pld.Stake))
	sb.UpdateAccount(pld.From, senderAcc)
	sb.UpdateValidator(receiverVal)

	return nil
}
