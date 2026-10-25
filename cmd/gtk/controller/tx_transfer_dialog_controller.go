//go:build gtk

package controller

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
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

func (c *TxTransferDialogController) Show() {
	c.applyDefaults()
	c.populateCombos()

	gtkutil.DropDownOnChanged(c.view.SenderDrop, c.onSenderChanged)
	gtkutil.EntryOnChanged(c.view.ReceiverEntry, c.onReceiverChanged)
	gtkutil.EntryOnChanged(c.view.FeeEntry, c.onFeeChanged)

	gtkutil.ConnectButtonSignal(c.view.AmountMaxButton, c.onMax)
	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	c.view.SenderDrop.SetSelected(0)
	c.onSenderChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func (c *TxTransferDialogController) applyDefaults() {
	setDefaultFee(c.model, c.view.FeeEntry)
}

func (c *TxTransferDialogController) populateCombos() {
	gtkutil.DropDownFromAddressListWithAmount(c.view.SenderDrop, c.model.ListAccountAddresses(), c.model.Balance)
}

func (c *TxTransferDialogController) onSenderChanged() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	if info := c.model.AddressInfo(sender); info != nil && info.Label != "" {
		setHintLabel(c.view.SenderHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHintLabel(c.view.SenderHint, "")
	}

	bal, err := c.model.Balance(sender)
	if err == nil {
		setBalanceHint(c.view.AmountHint, bal)
	} else {
		setHintLabel(c.view.AmountHint, "")
	}
}

func (c *TxTransferDialogController) onReceiverChanged() {
	receiver := gtkutil.EntryGetText(c.view.ReceiverEntry)
	if info := c.model.AddressInfo(receiver); info != nil && info.Label != "" {
		setHintLabel(c.view.ReceiverHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHintLabel(c.view.ReceiverHint, "")
	}
}

// onMax fills in the whole balance of the sender, minus the fee.
func (c *TxTransferDialogController) onMax() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	balance, err := c.model.Balance(sender)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	fee, err := amount.FromString(gtkutil.EntryGetText(c.view.FeeEntry))
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	maxAmount, err := model.MaxTransferAmount(balance, fee)
	setMaxAmount(c.view.Window, c.view.AmountEntry, maxAmount, err)
}

func (c *TxTransferDialogController) onFeeChanged() {
	setHintLabel(c.view.FeeHint, "")
}

func (c *TxTransferDialogController) onSend() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	receiver := gtkutil.EntryGetText(c.view.ReceiverEntry)
	amountStr := gtkutil.EntryGetText(c.view.AmountEntry)
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

	trx, err := c.model.MakeTransferTx(sender, receiver, amt, fee, memo)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	msg := fmt.Sprintf(`
<b>Transaction Details</b>
<tt>
Type:     Transfer
From:     %s
To:       %s
Amount:   %s
Fee:      %s
Memo:     %s
</tt>

You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`, sender, receiver, amt, trx.Fee(), trx.Memo())

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxTransferDialogController) onCancel() {
	c.view.Window.Close()
}
