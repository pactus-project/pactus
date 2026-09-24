package executor

import (
	"bytes"

	"github.com/pactus-project/gopkg/logger"
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

func (e *AnchorExecutor) Check(_ sandbox.SandboxReader, _ bool) error {
	if err := e.pld.BasicCheck(); err != nil {
		return err
	}
	if e.fee < 0 || e.fee > amount.MaxNanoPAC {
		return ErrAmountOverflow
	}

	if e.pld.Action == payload.AnchorActionDelete {
		return e.checkDelete()
	}

	return e.checkSet()
}

// Execute applies the anchor without validating it.
// The caller must run Check first, or trust a certified block.
func (e *AnchorExecutor) Execute(sbx sandbox.Sandbox) {
	if e.pld.Action == payload.AnchorActionDelete {
		refund := e.acc.LockedDeposit() - e.fee
		e.acc.ClearAnchor()
		e.acc.AddToBalance(refund)
		sbx.UpdateAccount(e.pld.From, e.acc)

		return
	}

	if err := e.acc.SetAnchor(e.nextAnchor(sbx)); err != nil {
		logger.Panic("anchor executed without a valid check", "error", err)
	}
	e.acc.SubtractFromBalance(e.pld.Deposit + e.fee)
	sbx.UpdateAccount(e.pld.From, e.acc)
}

func (e *AnchorExecutor) checkSet() error {
	deposit := e.pld.Deposit
	if deposit > amount.MaxNanoPAC-e.fee {
		return ErrAmountOverflow
	}
	if e.acc.Balance() < deposit+e.fee {
		return ErrInsufficientFunds
	}

	locked := e.acc.LockedDeposit()
	if locked < 0 || locked > amount.MaxNanoPAC {
		return ErrAmountOverflow
	}
	if locked > amount.MaxNanoPAC-deposit {
		return ErrAmountOverflow
	}
	if !e.acc.HasAnchor() && deposit < MinAnchorDeposit {
		return ErrAnchorDepositTooSmall
	}

	return nil
}

func (e *AnchorExecutor) checkDelete() error {
	if !e.acc.HasAnchor() {
		return ErrAnchorNotFound
	}

	locked := e.acc.LockedDeposit()
	if locked < 0 || locked > amount.MaxNanoPAC {
		return ErrAmountOverflow
	}
	if locked < e.fee {
		return ErrInsufficientFunds
	}

	refund := locked - e.fee
	if e.acc.Balance() > amount.MaxNanoPAC-refund {
		return ErrAmountOverflow
	}

	return nil
}

// nextAnchor builds the anchor a Set stores. It does not validate anything.
//
// CreatedAt is the block that created the slot. UpdatedAt is the block that set
// the current content (root hash, manifest URI and anchor type), so a Set that
// only adds deposit keeps the date the current digest was attested.
func (e *AnchorExecutor) nextAnchor(sbx sandbox.SandboxReader) account.AnchorData {
	height := sbx.CurrentHeight()
	unixTime := sbx.CurrentUnixTime()

	createdHeight, createdTime := height, unixTime
	updatedHeight, updatedTime := height, unixTime
	if e.acc.HasAnchor() {
		createdHeight = e.acc.CreatedAtHeight()
		createdTime = e.acc.CreatedAtTime()
		if e.keepsContent() {
			updatedHeight = e.acc.UpdatedAtHeight()
			updatedTime = e.acc.UpdatedAtTime()
		}
	}

	return account.AnchorData{
		RootHash:        e.pld.RootHash,
		ManifestURI:     e.pld.ManifestURI,
		AnchorType:      e.pld.AnchorType,
		LockedDeposit:   e.acc.LockedDeposit() + e.pld.Deposit,
		CreatedAtHeight: createdHeight,
		CreatedAtTime:   createdTime,
		UpdatedAtHeight: updatedHeight,
		UpdatedAtTime:   updatedTime,
	}
}

// keepsContent reports whether the Set stores the content the anchor already has.
func (e *AnchorExecutor) keepsContent() bool {
	return bytes.Equal(e.acc.RootHash(), e.pld.RootHash) &&
		e.acc.ManifestURI() == e.pld.ManifestURI &&
		e.acc.AnchorType() == e.pld.AnchorType
}
