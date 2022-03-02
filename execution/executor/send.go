package executor

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
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

	// In not-restrict mode we accepts all subsidy transactions for the current height
	// There might be more than one valid subsidy transaction per height
	// because There might be more than one proposal per height
	if !e.strict && trx.IsMintbaseTx() {
		if trx.Sequence() != sb.CurrentHeight() {
			return errors.Errorf(errors.ErrInvalidTx,
				"Subsidy transaction is not for current height. Expected :%d, got: %d",
				sb.CurrentHeight(), trx.Sequence())
		}

		return nil
	}

	senderAcc := sb.Account(pld.Sender)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve sender account")
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
		return errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}
	if senderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence, Expected: %v, got: %v", senderAcc.Sequence()+1, trx.Sequence())
	}

	senderAcc.IncSequence()
	senderAcc.SubtractFromBalance(pld.Amount + trx.Fee())
	receiverAcc.AddToBalance(pld.Amount)

	sb.UpdateAccount(senderAcc)
	sb.UpdateAccount(receiverAcc)

	e.fee = trx.Fee()

	return nil
}

func (e *SendExecutor) Fee() int64 {
	return e.fee
}
