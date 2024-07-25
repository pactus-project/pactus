package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TransferExecutor struct {
	sb       sandbox.Sandbox
	pld      *payload.TransferPayload
	fee      amount.Amount
	sender   *account.Account
	receiver *account.Account
}

func newTransferExecutor(trx *tx.Tx, sb sandbox.Sandbox) (*TransferExecutor, error) {
	pld := trx.Payload().(*payload.TransferPayload)

	sender := sb.Account(pld.From)
	if sender == nil {
		return nil, AccountNotFoundError{Address: pld.From}
	}

	var receiver *account.Account
	if pld.To == pld.From {
		receiver = sender
	} else {
		receiver = sb.Account(pld.To)
		if receiver == nil {
			receiver = sb.MakeNewAccount(pld.To)
		}
	}

	return &TransferExecutor{
		sb:       sb,
		pld:      pld,
		fee:      trx.Fee(),
		sender:   sender,
		receiver: receiver,
	}, nil
}

func (e *TransferExecutor) Check(_ bool) error {
	if e.sender.Balance() < e.pld.Amount+e.fee {
		return ErrInsufficientFunds
	}

	return nil
}

func (e *TransferExecutor) Execute() {
	e.sender.SubtractFromBalance(e.pld.Amount + e.fee)
	e.receiver.AddToBalance(e.pld.Amount)

	e.sb.UpdateAccount(e.pld.From, e.sender)
	e.sb.UpdateAccount(e.pld.To, e.receiver)
}
