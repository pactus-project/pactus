//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TxWithdrawDialogController struct {
	view  *view.TxWithdrawDialogView
	model *model.WalletModel

	// delegateOwner is the stake owner that signs the withdraw of a delegated
	// validator (PIP-49). It is empty for a validator that is not delegated.
	delegateOwner string

	// seq numbers the delegation reads, so a late answer for a validator the
	// user no longer has selected is ignored.
	seq uint64
}

func NewTxWithdrawDialogController(
	view *view.TxWithdrawDialogView,
	model *model.WalletModel,
) *TxWithdrawDialogController {
	return &TxWithdrawDialogController{view: view, model: model}
}

func (c *TxWithdrawDialogController) Show() {
	c.applyDefaults()
	c.populateCombos()

	gtkutil.ComboBoxOnChanged(c.view.ValidatorCombo, c.onValidatorChanged)
	gtkutil.ComboBoxOnChanged(c.view.ReceiverCombo, c.onReceiverChanged)
	gtkutil.EntryOnChanged(c.view.FeeEntry, c.onFeeChanged)

	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	gtkutil.ComboBoxSelectFirst(c.view.ValidatorCombo)
	c.onValidatorChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func (c *TxWithdrawDialogController) applyDefaults() {
	setDefaultFee(c.model, c.view.FeeEntry)
}

func (c *TxWithdrawDialogController) populateCombos() {
	gtkutil.ComboBoxFromAddressList(c.view.ValidatorCombo, c.model.ListValidatorAddresses())
	gtkutil.ComboBoxFromAddressList(c.view.ReceiverCombo, c.model.ListAccountAddresses())
}

func (c *TxWithdrawDialogController) onValidatorChanged() {
	validator := gtkutil.ComboBoxGetSelectedAddress(c.view.ValidatorCombo)
	setHintLabel(c.view.ValidatorHint, validatorHint(c.model, validator))

	if stake, err := c.model.Stake(validator); err == nil {
		setHintLabel(c.view.StakeHint, fmt.Sprintf("Validator Stake: %s", stake))
	} else {
		setHintLabel(c.view.StakeHint, "")
	}

	// Send stays off until the node says whether the validator is delegated.
	c.seq++
	seq := c.seq
	c.delegateOwner = ""
	c.view.ButtonSend.SetSensitive(false)
	loadDelegatedValidator(c.model, validator, "withdraw", func(state *delegatedValidator) {
		if seq != c.seq {
			return // the user picked another validator meanwhile
		}
		c.delegateOwner = state.owner
		showDelegationNotice(c.view.DelegationHint, state)
		c.view.ButtonSend.SetSensitive(!state.delegated || state.canSign)
	})
}

func (c *TxWithdrawDialogController) onReceiverChanged() {
	receiver := gtkutil.ComboBoxGetSelectedText(c.view.ReceiverCombo)
	if info := c.model.AddressInfo(receiver); info != nil && info.Label != "" {
		setHintLabel(c.view.ReceiverHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHintLabel(c.view.ReceiverHint, "")
	}
}

func (c *TxWithdrawDialogController) onFeeChanged() {
	_ = payload.TypeWithdraw
	setHintLabel(c.view.FeeHint, "")
}

func (c *TxWithdrawDialogController) onSend() {
	sender := gtkutil.ComboBoxGetSelectedAddress(c.view.ValidatorCombo)
	receiver := gtkutil.ComboBoxGetSelectedText(c.view.ReceiverCombo)
	amountStr := gtkutil.EntryGetText(c.view.StakeEntry)
	feeStr := gtkutil.EntryGetText(c.view.FeeEntry)
	memo := gtkutil.EntryGetText(c.view.MemoEntry)

	amt, err := amount.FromString(amountStr)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	fee, err := amount.FromString(feeStr)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	trx, err := c.model.MakeWithdrawTx(sender, receiver, amt, fee, memo, c.delegateOwner)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}
	if err := verifyWithdrawTx(trx, sender, receiver, amt, c.delegateOwner); err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	signer := ""
	if c.delegateOwner != "" {
		signer = fmt.Sprintf("\nOwner:  %s (signs this transaction)", c.delegateOwner)
	}

	msg := fmt.Sprintf(`
<b>Transaction Details</b>
<tt>
Type:   Withdraw
From:   %s
To:     %s%s
Amount: %s
Fee:    %s
Memo:   %s
</tt>

You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`, sender, receiver, signer, amt, trx.Fee(),
		gtkutil.EscapeMarkup(trx.Memo()))

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxWithdrawDialogController) onCancel() {
	c.view.Window.Close()
}

// verifyWithdrawTx checks that the node built the withdraw the user asked for.
func verifyWithdrawTx(trx *tx.Tx, validator, receiver string, amt amount.Amount, delegateOwner string) error {
	pld, ok := trx.Payload().(*payload.WithdrawPayload)
	if !ok || pld.From.String() != validator || pld.To.String() != receiver || pld.Amount != amt {
		return errTxMismatch
	}
	if delegateOwner == "" {
		if pld.IsDelegated() {
			return errTxMismatch
		}

		return nil
	}
	if !pld.IsDelegated() || pld.DelegateOwner.String() != delegateOwner {
		return errTxMismatch
	}

	return nil
}
