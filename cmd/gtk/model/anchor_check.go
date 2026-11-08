//go:build gtk

package model

import (
	"errors"

	"github.com/pactus-project/pactus/types/amount"
)

var (
	// ErrAnchorNotFound is returned when a delete targets an account without an anchor.
	ErrAnchorNotFound = errors.New("this account has no anchor, so there is nothing to delete")
	// ErrAnchorLockedBelowFee is returned when the locked deposit cannot pay the fee of a delete.
	ErrAnchorLockedBelowFee = errors.New("the locked deposit is lower than the fee, so the delete would be rejected")
	// ErrAnchorDepositTooSmall is returned when a new anchor locks less than the minimum.
	ErrAnchorDepositTooSmall = errors.New("a new anchor needs a deposit of at least 1 PAC")
	// ErrAnchorNegativeDeposit is returned when the deposit is below zero.
	ErrAnchorNegativeDeposit = errors.New("the deposit cannot be negative")
	// ErrAnchorInsufficientFunds is returned when the balance cannot pay the deposit and the fee.
	ErrAnchorInsufficientFunds = errors.New("the balance does not cover the deposit plus the fee")
)

// CurrentAnchor is what the dialog knows about the anchor an account has now.
type CurrentAnchor struct {
	Exists bool
	Type   uint8
	Locked amount.Amount
}

// AnchorSetCheck is the outcome of checking a set before signing it.
type AnchorSetCheck struct {
	// Err is set when the transaction would be rejected.
	Err error
	// Replaces is true when the account already has an anchor.
	Replaces bool
	// TypeChanges is true when the anchor it replaces has another type.
	TypeChanges bool
	// LargeDeposit is true when the set locks more than the minimum.
	LargeDeposit bool
	// LockedAfter is the deposit the anchor holds after the transaction.
	LockedAfter amount.Amount
}

// CheckAnchorSet checks a set of the given type and deposit on an account.
func CheckAnchorSet(current CurrentAnchor, anchorType uint8,
	deposit, fee, balance amount.Amount,
) AnchorSetCheck {
	check := AnchorSetCheck{
		Replaces:     current.Exists,
		TypeChanges:  current.Exists && current.Type != anchorType,
		LargeDeposit: deposit > MinAnchorDeposit,
		LockedAfter:  current.Locked + deposit,
	}

	switch {
	case deposit < 0:
		check.Err = ErrAnchorNegativeDeposit
	case !current.Exists && deposit < MinAnchorDeposit:
		check.Err = ErrAnchorDepositTooSmall
	case balance < deposit+fee:
		check.Err = ErrAnchorInsufficientFunds
	}

	return check
}

// AnchorDeleteCheck is the outcome of checking a delete before signing it.
type AnchorDeleteCheck struct {
	// Err is set when the transaction would be rejected.
	Err error
	// Refund is the amount that returns to the balance: the locked deposit minus the fee.
	Refund amount.Amount
}

// CheckAnchorDelete checks a delete of the account's anchor with the given fee.
func CheckAnchorDelete(current CurrentAnchor, fee amount.Amount) AnchorDeleteCheck {
	switch {
	case !current.Exists:
		return AnchorDeleteCheck{Err: ErrAnchorNotFound}
	case !CanDeleteAnchor(current.Locked, fee):
		return AnchorDeleteCheck{Err: ErrAnchorLockedBelowFee}
	default:
		return AnchorDeleteCheck{Refund: current.Locked - fee}
	}
}

// DefaultAnchorDeposit returns the deposit the dialog proposes for a set:
// what is missing to reach the target of a link, or, without a link, the
// minimum for a new anchor and nothing for an update.
func DefaultAnchorDeposit(current CurrentAnchor, linkTarget *amount.Amount) amount.Amount {
	locked := amount.Amount(0)
	if current.Exists {
		locked = current.Locked
	}
	if linkTarget != nil {
		return AnchorDepositToAdd(*linkTarget, locked)
	}
	if current.Exists {
		return 0
	}

	return MinAnchorDeposit
}
