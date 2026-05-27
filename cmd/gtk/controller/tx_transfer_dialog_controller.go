//go111:build gtk

package controller

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TxTransferDialogController struct {
	view  *view.TxTransferDialogView
	model *model.WalletModel
}

func NewTxTransferDialogController(
	view *view.TxTransferDialogView,
	model *model.WalletModel,
) *TxTransferDialogController {
	return &TxTransferDialogController{view: view, model: model}
}

func setHintLabel(lbl *gtk.Label, hint string) {
	if hint == "" {
		lbl.SetMarkup("")

		return
	}
	lbl.SetMarkup(gtkutil.SmallGray(hint))
}

func (c *TxTransferDialogController) Run() {
	// Defaults
	if info, err := c.model.WalletInfo(); err == nil {
		c.view.FeeEntry.SetText(fmt.Sprintf("%g", info.DefaultFee.ToPAC()))
	}

	// Fill sender accounts
	for _, ai := range c.model.ListAccountAddresses() {
		c.view.SenderCombo.Append(ai.Address, ai.Address)
	}
	c.view.SenderCombo.SetActive(0)

	c.view.ConnectSignals(map[string]any{
		"on_sender_changed":   c.onSenderChanged,
		"on_receiver_changed": c.onReceiverChanged,
		"on_fee_changed":      c.onFeeChanged,
		"on_send":             c.onSend,
		"on_cancel":           c.onCancel,
	})

	c.onSenderChanged()

	gtkutil.ShowNonModalDialog(c.view.Window)
}

func (c *TxTransferDialogController) onSenderChanged() {
	sender := c.view.SenderCombo.ActiveID()
	if info := c.model.AddressInfo(sender); info != nil && info.Label != "" {
		setHintLabel(c.view.SenderHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHintLabel(c.view.SenderHint, "")
	}

	bal, err := c.model.Balance(sender)
	if err == nil {
		setHintLabel(c.view.AmountHint, fmt.Sprintf("Account Balance: %s", bal))
	} else {
		setHintLabel(c.view.AmountHint, "")
	}
}

func (c *TxTransferDialogController) onReceiverChanged() {
	receiver := gtkutil.GetEntryText(c.view.ReceiverEntry)
	if info := c.model.AddressInfo(receiver); info != nil && info.Label != "" {
		setHintLabel(c.view.ReceiverHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHintLabel(c.view.ReceiverHint, "")
	}
}

func (c *TxTransferDialogController) onFeeChanged() {
	// Placeholder (confirmation time estimation)
	_ = payload.TypeTransfer
	setHintLabel(c.view.FeeHint, "")
}

func (c *TxTransferDialogController) onSend() {
	sender := c.view.SenderCombo.ActiveID()
	receiver := gtkutil.GetEntryText(c.view.ReceiverEntry)
	amountStr := gtkutil.GetEntryText(c.view.AmountEntry)
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

	trx, err := c.model.MakeTransferTx(sender, receiver, amt, fee, memo)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:     Transfer
From:     %s
To:       %s
Amount:   %s
Fee:      %s
Memo:     %s
</tt>

You are going to sign and broadcast this transaction.
<b>⚠️ This action cannot be undone.</b>
Do you want to continue with this transaction?`, sender, receiver, amt, trx.Fee(), trx.Memo())

	if !confirmAndSend(c.view.Window, c.model, msg, trx) {
		return
	}

	c.view.Window.Close()
}

func (c *TxTransferDialogController) onCancel() {
	c.view.Window.Close()
}
