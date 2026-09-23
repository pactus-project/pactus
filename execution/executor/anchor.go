package executor

import (
	"github.com/pactus-project/pactus/sandbox"
	"github.com/pactus-project/pactus/types/account"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/protocol"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

const anchorProtocolVersion = protocol.Version(5)

// MinAnchorDeposit is the minimum locked deposit for a new anchor.
const MinAnchorDeposit amount.Amount = 1_000_000_000

type AnchorExecutor struct {
	pld *payload.AnchorPayload
	fee amount.Amount
	acc *account.Account
}

type anchorChange struct {
	spend  amount.Amount
	refund amount.Amount
	remove bool
	anchor account.AnchorData
}

func newAnchorExecutor(trx *tx.Tx, sbx sandbox.Sandbox) (*AnchorExecutor, error) {
	pld := trx.Payload().(*payload.AnchorPayload)

	acc := sbx.Account(pld.From)
	if acc == nil {
		return nil, AccountNotFoundError{Address: pld.From}
	}

	return &AnchorExecutor{
		pld: pld,
		fee: trx.Fee(),
		acc: acc,
	}, nil
}

func (e *AnchorExecutor) Check(sbx sandbox.SandboxReader, _ bool) error {
	_, err := e.prepare(sbx)

	return err
}

func (e *AnchorExecutor) Execute(sbx sandbox.Sandbox) {
	change, err := e.prepare(sbx)
	if err != nil {
		return
	}

	if change.remove {
		e.acc.ClearAnchor()
		e.acc.AddToBalance(change.refund)
		sbx.UpdateAccount(e.pld.From, e.acc)

		return
	}

	if err := e.acc.SetAnchor(change.anchor); err != nil {
		return
	}

	e.acc.SubtractFromBalance(change.spend)
	sbx.UpdateAccount(e.pld.From, e.acc)
}

func (e *AnchorExecutor) prepare(sbx sandbox.SandboxReader) (*anchorChange, error) {
	if err := e.pld.BasicCheck(); err != nil {
		return nil, err
	}
	if e.fee < 0 || e.fee > amount.MaxNanoPAC {
		return nil, ErrAmountOverflow
	}

	if e.pld.Action == payload.AnchorActionDelete {
		return e.prepareDelete()
	}

	return e.prepareSet(sbx)
}

func (e *AnchorExecutor) prepareSet(sbx sandbox.SandboxReader) (*anchorChange, error) {
	deposit := e.pld.Deposit
	if deposit < 0 || deposit > amount.MaxNanoPAC {
		return nil, ErrAmountOverflow
	}
	if deposit > amount.MaxNanoPAC-e.fee {
		return nil, ErrAmountOverflow
	}

	cost := deposit + e.fee
	if e.acc.Balance() < cost {
		return nil, ErrInsufficientFunds
	}

	locked := e.acc.LockedDeposit()
	if locked < 0 || locked > amount.MaxNanoPAC {
		return nil, ErrAmountOverflow
	}
	if deposit > 0 && locked > amount.MaxNanoPAC-deposit {
		return nil, ErrAmountOverflow
	}

	newLocked := locked + deposit
	createdHeight := sbx.CurrentHeight()
	createdTime := sbx.CurrentUnixTime()
	if e.acc.HasAnchor() {
		createdHeight = e.acc.CreatedAtHeight()
		createdTime = e.acc.CreatedAtTime()
	} else if newLocked < MinAnchorDeposit {
		return nil, ErrAnchorDepositTooSmall
	}

	return &anchorChange{
		spend: cost,
		anchor: account.AnchorData{
			RootHash:        e.pld.RootHash,
			ManifestURI:     e.pld.ManifestURI,
			AnchorType:      e.pld.AnchorType,
			LockedDeposit:   newLocked,
			CreatedAtHeight: createdHeight,
			CreatedAtTime:   createdTime,
			UpdatedAtHeight: sbx.CurrentHeight(),
			UpdatedAtTime:   sbx.CurrentUnixTime(),
		},
	}, nil
}

func (e *AnchorExecutor) prepareDelete() (*anchorChange, error) {
	if !e.acc.HasAnchor() {
		return nil, ErrAnchorNotFound
	}

	locked := e.acc.LockedDeposit()
	if locked < 0 || locked > amount.MaxNanoPAC {
		return nil, ErrAmountOverflow
	}
	if locked < e.fee {
		return nil, ErrInsufficientFunds
	}

	refund := locked - e.fee
	if refund > 0 && e.acc.Balance() > amount.MaxNanoPAC-refund {
		return nil, ErrAmountOverflow
	}

	return &anchorChange{
		refund: refund,
		remove: true,
	}, nil
}
