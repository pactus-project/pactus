package executor

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/sandbox"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
)

type SendExecutor struct {
	sandbox sandbox.Sandbox
}

func NewSendExecutor(sandbox sandbox.Sandbox) *SendExecutor {
	return &SendExecutor{sandbox}
}

func (e *SendExecutor) Execute(trx *tx.Tx) error {
	pld := trx.Payload().(*payload.SendPayload)

	senderAcc := e.sandbox.Account(pld.Sender)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve sender account")
	}
	receiverAcc := e.sandbox.Account(pld.Receiver)
	if receiverAcc == nil {
		receiverAcc = account.NewAccount(pld.Receiver)
	}
	if senderAcc.Sequence()+1 != trx.Sequence() {
		return errors.Errorf(errors.ErrInvalidTx, "Invalid sequence")
	}
	if senderAcc.Balance() < pld.Amount+trx.Fee() {
		return errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}
	senderAcc.IncSequence()
	senderAcc.SubtractFromBalance(pld.Amount + trx.Fee())
	receiverAcc.AddToBalance(pld.Amount)

	e.sandbox.UpdateAccount(senderAcc)
	e.sandbox.UpdateAccount(receiverAcc)

	return nil
}
