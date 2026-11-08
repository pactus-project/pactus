//go:build gtk

package controller

import (
	"fmt"
	"time"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/types/amount"
)

// validatorHint returns the stake and the label of a validator, for the hint
// below a validator field.
func validatorHint(walletModel *model.WalletModel, validator string) string {
	hint := ""
	if stake, err := walletModel.Stake(validator); err == nil {
		hint = fmt.Sprintf("stake: %s", stake)
	}
	if info := walletModel.AddressInfo(validator); info != nil && info.Label != "" {
		if hint != "" {
			hint += ", "
		}
		hint += "label: " + info.Label
	}

	return hint
}

// delegatedValidator describes the delegation of a validator (PIP-49), as the
// unbond and withdraw dialogs need it.
type delegatedValidator struct {
	action    string
	delegated bool
	// owner is the stake owner. Only that account can sign the transaction.
	owner string
	// canSign is true when the stake owner is an account of this wallet.
	canSign bool
	// expired is true when the delegation has ended for the next block.
	expired bool
	expiry  uint32
	endsAt  time.Time
	share   amount.Amount
}

// loadDelegatedValidator reads the delegation of a validator in the background
// and calls apply on the UI thread. The node can take a while to answer when it
// is busy, for example while it syncs, so the dialog must not wait for it.
func loadDelegatedValidator(walletModel *model.WalletModel, validator, action string,
	apply func(*delegatedValidator),
) {
	go func() {
		state := checkDelegatedValidator(walletModel, validator, action)
		gtkutil.IdleAddAsync(func() { apply(&state) })
	}()
}

// checkDelegatedValidator reads the delegation of a validator from the node.
// It calls the node, so run it off the UI thread.
func checkDelegatedValidator(walletModel *model.WalletModel, validator, action string) delegatedValidator {
	state := delegatedValidator{action: action}

	info := walletModel.ValidatorInfo(validator)
	if info == nil || !info.IsDelegated {
		return state
	}

	state.delegated = true
	state.owner = info.DelegateOwner
	state.expiry = info.DelegateExpiry
	state.share = amount.Amount(info.DelegateShare)
	state.canSign = walletModel.AddressInfo(info.DelegateOwner) != nil

	status, err := walletModel.ChainStatus()
	if err == nil {
		// The next block is the first one the transaction can be in.
		state.expired = info.DelegateExpiry <= status.Height+1
		state.endsAt = model.EstimateBlockTime(status.Height, info.DelegateExpiry, time.Now())
	}

	return state
}

// showDelegationNotice explains who must sign for a delegated validator, or
// hides the notice for a validator that is not delegated.
func showDelegationNotice(label *gtk.Label, state *delegatedValidator) {
	if !state.delegated {
		label.SetVisible(false)

		return
	}

	ends := fmt.Sprintf("block %d", state.expiry)
	if !state.endsAt.IsZero() {
		ends += ", around " + state.endsAt.Format("02 Jan 2006")
	}

	var msg string
	if state.canSign {
		msg = fmt.Sprintf("<b>Delegated validator (PIP-49).</b> The stake owner %s is in this wallet "+
			"and signs the %s. The delegation ends at %s.",
			gtkutil.EscapeMarkup(state.owner), state.action, ends)
	} else {
		msg = fmt.Sprintf("<b>Delegated validator (PIP-49).</b> Only its stake owner %s can %s it, "+
			"and that account is not in this wallet. The delegation ends at %s.",
			gtkutil.EscapeMarkup(state.owner), state.action, ends)
	}
	msg += "\nOwner reward share: " + formatShare(state.share) + " of each block reward."
	label.SetMarkup(msg)
	label.SetVisible(true)
}

// formatShare shows an owner share as an amount and a percent of the block reward.
func formatShare(share amount.Amount) string {
	return fmt.Sprintf("%s (%.4g%%)", share, model.DelegateSharePercent(share))
}
