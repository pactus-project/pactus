package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/errors"
)

type SendExecutor struct {
	fee    int64
	strict bool
}

func NewSendExecutor(strict bool) *SendExecutor {
	return &SendExecutor{strict: strict}
}

func (e *SendExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.SendPayload)

	if !e.strict && trx.IsSubsidyTx() {
		// In not-strict mode all subsidy transactions for the current height are valid.
		// There might be more than one valid subsidy transaction per height,
		// because there might be more than one proposal per height.
		if uint32(trx.Sequence()) != sb.CurrentHeight() {
			return errors.Errorf(errors.ErrInvalidSequence,
				"subsidy transaction is not for current height, expected :%d, got: %d",
				sb.CurrentHeight(), trx.Sequence())
		}

		return nil
	}

	senderAcc := sb.Account(pld.Sender)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve sender account")
	}
	var receiverAcc *account.Account
	if pld.Receiver.EqualsTo(pld.Sender) {
		receiverAcc = senderAcc
	} else {
		receiverAcc = sb.Account(pld.Receiver)
		if receiverAcc == nil {
			receiverAcc = sb.MakeNewAccount(pld.Receiver)
		}
	}
	if senderAcc.Balance() < pld.Amount+trx.Fee() {
		return errors.Error(errors.ErrInsufficientFunds)
	}
	if senderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidSequence,
			"expected: %v, got: %v", senderAcc.Sequence()+1, trx.Sequence())
	}

	senderAcc.IncSequence()
	senderAcc.SubtractFromBalance(pld.Amount + trx.Fee())
	receiverAcc.AddToBalance(pld.Amount)

	sb.UpdateAccount(pld.Sender, senderAcc)
	sb.UpdateAccount(pld.Receiver, receiverAcc)

	e.fee = trx.Fee()

	return nil
}

func (e *SendExecutor) Fee() int64 {
	return e.fee
}
