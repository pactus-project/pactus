package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TransferExecutor struct {
	pld         *payload.TransferPayload
	fee         amount.Amount
	sender      *account.Account
	receiver    *account.Account
	newReceiver bool
}

func newTransferExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*TransferExecutor, error) {
	pld := trx.Payload().(*payload.TransferPayload)

	sender := sbx.Account(pld.From)
	if sender == nil {
		return nil, AccountNotFoundError{Address: pld.From}
	}

	var receiver *account.Account
	newReceiver := false
	if pld.To == pld.From {
		receiver = sender
	} else {
		receiver = sbx.Account(pld.To)
		if receiver == nil {
			// sbx.MakeNewAccount is deferred to Execute, which only runs
			// once Check has accepted the transaction, so a transaction
			// Check later rejects (e.g. insufficient funds) never leaks an
			// empty account into the sandbox. This placeholder is for
			// validation only and is never written back.
			receiver = account.NewAccount(0)
			newReceiver = true
		}
	}

	return &TransferExecutor{
		pld:         pld,
		fee:         trx.Fee(),
		sender:      sender,
		receiver:    receiver,
		newReceiver: newReceiver,
	}, nil
}

func (e *TransferExecutor) Check(_ sandbox.SandboxReader, _ bool) error {
	if e.sender.Balance() < e.pld.Amount+e.fee {
		return ErrInsufficientFunds
	}

	return nil
}

func (e *TransferExecutor) Execute(sbx sandbox.Sandbox) {
	if e.newReceiver {
		e.receiver = sbx.MakeNewAccount(e.pld.To)
	}

	e.sender.SubtractFromBalance(e.pld.Amount + e.fee)
	e.receiver.AddToBalance(e.pld.Amount)

	sbx.UpdateAccount(e.pld.From, e.sender)
	sbx.UpdateAccount(e.pld.To, e.receiver)
}
