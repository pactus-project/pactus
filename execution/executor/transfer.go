package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util/errors"
)

type TransferExecutor struct {
	strict bool
}

func NewTransferExecutor(strict bool) *TransferExecutor {
	return &TransferExecutor{strict: strict}
}

func (*TransferExecutor) Execute(trx *tx.Tx, sb sandbox.Sandbox) error {
	pld := trx.Payload().(*payload.TransferPayload)

	senderAcc := sb.Account(pld.From)
	if senderAcc == nil {
		return errors.Errorf(errors.ErrInvalidAddress,
			"unable to retrieve sender account")
	}
	var receiverAcc *account.Account
	if pld.To == pld.From {
		receiverAcc = senderAcc
	} else {
		receiverAcc = sb.Account(pld.To)
		if receiverAcc == nil {
			receiverAcc = sb.MakeNewAccount(pld.To)
		}
	}

	if senderAcc.Balance() < pld.Amount+trx.Fee() {
		return ErrInsufficientFunds
	}

	senderAcc.SubtractFromBalance(pld.Amount + trx.Fee())
	receiverAcc.AddToBalance(pld.Amount)

	sb.UpdateAccount(pld.From, senderAcc)
	sb.UpdateAccount(pld.To, receiverAcc)

	return nil
}
