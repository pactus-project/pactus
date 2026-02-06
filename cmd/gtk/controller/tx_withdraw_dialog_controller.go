//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TxWithdrawDialogController struct {
	view  *view.TxWithdrawDialogView
	model *model.WalletModel
}

func NewTxWithdrawDialogController(
	view *view.TxWithdrawDialogView,
	model *model.WalletModel,
) *TxWithdrawDialogController {
	return &TxWithdrawDialogController{view: view, model: model}
}

func (c *TxWithdrawDialogController) Run() {
	c.applyDefaults()
	c.populateCombos()

	onCancel := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_sender_changed":   func() { c.onSenderChanged() },
		"on_receiver_changed": func() { c.onReceiverChanged() },
		"on_fee_changed":      func() { c.onFeeChanged() },
		"on_send":             func() { c.onSend() },
		"on_cancel":           onCancel,
	})

	c.onSenderChanged()
	gtkutil.RunDialog(c.view.Dialog)
}

func (c *TxWithdrawDialogController) applyDefaults() {
	if info, err := c.model.WalletInfo(); err == nil {
		c.view.FeeEntry.SetText(fmt.Sprintf("%g", info.DefaultFee.ToPAC()))
	}
}

func (c *TxWithdrawDialogController) populateCombos() {
	for _, ai := range c.model.ListAddresses(crypto.AddressTypeValidator) {
		c.view.ValidatorCombo.Append(ai.Address, ai.Address)
	}
	c.view.ValidatorCombo.SetActive(0)

	for _, ai := range c.model.ListAddresses(crypto.AddressTypeBLSAccount, crypto.AddressTypeEd25519Account) {
		c.view.ReceiverCombo.Append(ai.Address, ai.Address)
	}
}

func (c *TxWithdrawDialogController) onSenderChanged() {
	sender := c.view.ValidatorCombo.GetActiveID()

	stake, err := c.model.Stake(sender)

	hint := ""
	if err == nil {
		hint = fmt.Sprintf("stake: %s", stake)
	}
	if info := c.model.AddressInfo(sender); info != nil && info.Label != "" {
		if hint != "" {
			hint += ", "
		}
		hint += "label: " + info.Label
	}
	setHintLabel(c.view.ValidatorHint, hint)

	if err == nil {
		setHintLabel(c.view.StakeHint, fmt.Sprintf("Validator Stake: %s", stake))
	} else {
		setHintLabel(c.view.StakeHint, "")
	}
}

func (c *TxWithdrawDialogController) onReceiverChanged() {
	receiverEntry, _ := c.view.ReceiverCombo.GetEntry()
	receiver := gtkutil.GetEntryText(receiverEntry)
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
	sender := c.view.ValidatorCombo.GetActiveID()
	receiverEntry, _ := c.view.ReceiverCombo.GetEntry()
	receiver := gtkutil.GetEntryText(receiverEntry)
	amountStr := gtkutil.GetEntryText(c.view.StakeEntry)
	feeStr := gtkutil.GetEntryText(c.view.FeeEntry)
	memo := gtkutil.GetEntryText(c.view.MemoEntry)

	amt, err := amount.FromString(amountStr)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	fee, err := amount.FromString(feeStr)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	trx, err := c.model.MakeWithdrawTx(sender, receiver, amt, fee, memo)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:   Withdraw
From:   %s
To:     %s
Amount: %s
Fee:    %s
Memo:   %s
</tt>

You are going to sign and broadcast this transaction.
<b>⚠️ This action cannot be undone.</b>
Do you want to continue with this transaction?`, sender, receiver, amt, trx.Fee(), trx.Memo())

	if !gtkutil.ShowQuestionDialog(c.view.Dialog, msg) {
		return
	}

	password, ok := PasswordProvider(c.model)
	if !ok {
		return
	}

	if err := c.model.SignTransaction(password, trx); err != nil {
		gtkutil.ShowError(err)

		return
	}
	txID, err := c.model.BroadcastTransaction(trx)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	gtkutil.ShowInfoDialog(c.view.Dialog,
		fmt.Sprintf("✅ Transaction sent successfully!\n\n"+
			"Transaction ID: <a href=\"https://pacviewer.com/transaction/%s\">%s</a>", txID, txID))
	c.view.Dialog.Close()
}
