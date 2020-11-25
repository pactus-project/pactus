package executor

import (
	"github.com/zarbchain/zarb-go/account"
	"github.com/zarbchain/zarb-go/errors"
	"github.com/zarbchain/zarb-go/tx"
)

type SendExecutor struct {
	sandbox Sandbox
}

func NewSendExecutor(sandbox Sandbox) *SendExecutor {
	return &SendExecutor{sandbox}
}

func (e *SendExecutor) Execute(trx *tx.Tx) (*tx.Receipt, error) {
	senderAcc := e.sandbox.Account(trx.Sender())
	if senderAcc == nil {
		return nil, errors.Errorf(errors.ErrInvalidTx, "Unable to retrieve sender account")
	}
	receiverAcc := e.sandbox.Account(trx.Receiver())
	if receiverAcc == nil {
		receiverAcc = account.NewAccount(trx.Receiver())
	}
	if senderAcc.Balance() < trx.Amount()+trx.Fee() {
		return nil, errors.Errorf(errors.ErrInvalidTx, "Insufficient balance")
	}
	senderAcc.SubtractFromBalance(trx.Amount() + trx.Fee())
	receiverAcc.AddToBalance(trx.Amount())

	e.sandbox.UpdateAccount(senderAcc)
	e.sandbox.UpdateAccount(receiverAcc)

	receipt := trx.GenerateReceipt(tx.Ok)
	return receipt, nil
}
