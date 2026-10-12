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

// MinAnchorDeposit is the minimum locked deposit for a new anchor set by PIP-50,
// from protocol version 5.
const MinAnchorDeposit amount.Amount = 1_000_000_000

// anchorDepositStep is the minimum deposit for a new anchor from a protocol version on.
type anchorDepositStep struct {
	from    protocol.Version
	deposit amount.Amount
}

// minAnchorDeposits is the schedule of the minimum deposit, by protocol version.
// A later PIP that changes the deposit appends a step for its version, so blocks
// of earlier versions keep validating with the value they had.
var minAnchorDeposits = []anchorDepositStep{
	{from: anchorProtocolVersion, deposit: MinAnchorDeposit},
}

// minAnchorDeposit returns the minimum deposit for a new anchor in a block of
// the given version.
func minAnchorDeposit(blockVersion protocol.Version) amount.Amount {
	deposit := minAnchorDeposits[0].deposit
	for _, step := range minAnchorDeposits {
		if blockVersion >= step.from {
			deposit = step.deposit
		}
	}

	return deposit
}

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

func (e *AnchorExecutor) Check(sbx sandbox.SandboxReader, _ bool) error {
	if err := e.pld.BasicCheck(); err != nil {
		return err
	}
	if e.fee < 0 || e.fee > amount.MaxNanoPAC {
		return ErrAmountOverflow
	}

	if e.pld.Action == payload.AnchorActionDelete {
		return e.checkDelete()
	}

	return e.checkSet(sbx.BlockVersion())
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

func (e *AnchorExecutor) checkSet(blockVersion protocol.Version) error {
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
	if !e.acc.HasAnchor() && deposit < minAnchorDeposit(blockVersion) {
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
// the current root hash: it is the date the current digest was attested, so a Set
// that adds deposit, moves the manifest or changes the anchor type keeps it.
func (e *AnchorExecutor) nextAnchor(sbx sandbox.SandboxReader) account.AnchorData {
	height := sbx.CurrentHeight()
	unixTime := sbx.CurrentUnixTime()

	createdHeight, createdTime := height, unixTime
	updatedHeight, updatedTime := height, unixTime
	if e.acc.HasAnchor() {
		createdHeight = e.acc.CreatedAtHeight()
		createdTime = e.acc.CreatedAtTime()
		if e.keepsRootHash() {
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

// keepsRootHash reports whether the Set stores the root hash the anchor already has.
func (e *AnchorExecutor) keepsRootHash() bool {
	return bytes.Equal(e.acc.RootHash(), e.pld.RootHash)
}
