package executor

import (
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
	"github.com/zarbchain/zarb-go/util"
)

type SendExecutor struct {
	sandbox sandbox.Sandbox
	fee     int64
}

func NewSendExecutor(sandbox sandbox.Sandbox) *SendExecutor {
	return &SendExecutor{sandbox: sandbox}
}

func (e *SendExecutor) Execute(trx *tx.Tx) error {
	pld := trx.Payload().(*payload.SendPayload)

	senderAcc := e.sandbox.Account(pld.Sender)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve sender account")
	}
	receiverAcc := e.sandbox.Account(pld.Receiver)
	if receiverAcc == nil {
		receiverAcc = e.sandbox.MakeNewAccount(pld.Receiver)
	}
	if senderAcc.Balance() < pld.Amount+trx.Fee() {
		return errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}
	if senderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence, Expected: %v, got: %v", senderAcc.Sequence()+1, trx.Sequence())
	}
	if trx.IsMintbaseTx() {
		if trx.Fee() != 0 {
			return errors.Errorf(errors.ErrInvalidTx, "Fee is wrong. expected: 0, got: %v", trx.Fee())
		}
	} else {
		fee := int64(float64(trx.Payload().Value()) * e.sandbox.FeeFraction())
		fee = util.Max64(fee, e.sandbox.MinFee())
		if trx.Fee() != fee {
			return errors.Errorf(errors.ErrInvalidTx, "Fee is wrong. expected: %v, got: %v", fee, trx.Fee())
		}
	}

	senderAcc.IncSequence()
	senderAcc.SubtractFromBalance(pld.Amount + trx.Fee())
	receiverAcc.AddToBalance(pld.Amount)

	e.sandbox.UpdateAccount(senderAcc)
	e.sandbox.UpdateAccount(receiverAcc)

	e.fee = trx.Fee()

	return nil
}

func (e *SendExecutor) Fee() int64 {
	return e.fee
}
