package executor

import (
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/types/validator"
)

type WithdrawExecutor struct {
	pld      *payload.WithdrawPayload
	fee      amount.Amount
	sender   *validator.Validator
	receiver *account.Account
}

func newWithdrawExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*WithdrawExecutor, error) {
	pld := trx.Payload().(*payload.WithdrawPayload)

	if pld.To.Type() == crypto.AddressTypeSecp256k1Account &&
		sbx.Params().BlockVersion <= protocol.ProtocolVersion3 {
		return nil, ErrSecp256k1AccountNotSupported
	}

	sender := sbx.Validator(pld.From)
	if sender == nil {
		return nil, ValidatorNotFoundError{Address: pld.From}
	}

	receiver := sbx.Account(pld.To)
	if receiver == nil {
		receiver = sbx.MakeNewAccount(pld.To)
	}

	return &WithdrawExecutor{
		pld:      pld,
		fee:      trx.Fee(),
		sender:   sender,
		receiver: receiver,
	}, nil
}

func (e *WithdrawExecutor) Check(sbx sandbox.SandboxReader, _ bool) error {
	if e.sender.Stake() < e.pld.Value()+e.fee {
		return ErrInsufficientFunds
	}

	if !e.sender.IsUnbonded() {
		return ErrValidatorBonded
	}

	if sbx.CurrentHeight() < e.sender.UnbondingHeight().SafeIncrease(sbx.Params().UnbondInterval) {
		return ErrUnbondingPeriod
	}

	// For delegated validators (PIP-49), only the stake owner can receive withdrawn principal.
	if e.sender.IsDelegated() {
		if e.pld.To != e.sender.DelegateOwner() {
			return ErrWithdrawMustGoToStakeOwner
		}
	}

	return nil
}

func (e *WithdrawExecutor) Execute(sbx sandbox.Sandbox) {
	e.sender.SubtractFromStake(e.pld.Amount + e.fee)
	e.receiver.AddToBalance(e.pld.Amount)

	sbx.UpdateValidator(e.sender)
	sbx.UpdateAccount(e.pld.To, e.receiver)
}
