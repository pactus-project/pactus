package executor

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type batchRecipient struct {
	Address crypto.Address
	Account *account.Account
	Amount  amount.Amount
	IsNew   bool
}

type BatchTransferExecutor struct {
	pld        *payload.BatchTransferPayload
	fee        amount.Amount
	sender     *account.Account
	recipients []batchRecipient
}

func newBatchTransferExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*BatchTransferExecutor, error) {
	pld := trx.Payload().(*payload.BatchTransferPayload)

	sender := sbx.Account(pld.From)
	if sender == nil {
		return nil, AccountNotFoundError{Address: pld.From}
	}

	recipients := make([]batchRecipient, len(pld.Recipients))
	for i, r := range pld.Recipients {
		if r.To == pld.From {
			recipients[i].Account = sender
		} else {
			receiver := sbx.Account(r.To)
			if receiver == nil {
				// sbx.MakeNewAccount is deferred to Execute, which only runs
				// once Check has accepted the transaction, so a transaction
				// Check later rejects never leaks an empty account into the
				// sandbox. This placeholder is for validation only and is
				// never written back.
				receiver = account.NewAccount(0)
				recipients[i].IsNew = true
			}
			recipients[i].Account = receiver
		}

		recipients[i].Address = r.To
		recipients[i].Amount = r.Amount
	}

	return &BatchTransferExecutor{
		pld:        pld,
		fee:        trx.Fee(),
		sender:     sender,
		recipients: recipients,
	}, nil
}

func (e *BatchTransferExecutor) Check(_ sandbox.SandboxReader, _ bool) error {
	if e.sender.Balance() < e.pld.Value()+e.fee {
		return ErrInsufficientFunds
	}

	return nil
}

func (e *BatchTransferExecutor) Execute(sbx sandbox.Sandbox) {
	e.sender.SubtractFromBalance(e.pld.Value() + e.fee)
	sbx.UpdateAccount(e.pld.From, e.sender)

	for i := range e.recipients {
		r := &e.recipients[i]
		if r.IsNew {
			r.Account = sbx.MakeNewAccount(r.Address)
		}

		r.Account.AddToBalance(r.Amount)

		sbx.UpdateAccount(r.Address, r.Account)
	}
}
