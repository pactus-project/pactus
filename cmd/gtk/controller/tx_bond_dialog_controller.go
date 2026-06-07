//go:build gtk

//nolint:staticcheck // Using depreciated widgets
package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TxBondDialogController struct {
	view  *view.TxBondDialogView
	model *model.WalletModel
}

func NewTxBondDialogController(
	view *view.TxBondDialogView,
	model *model.WalletModel,
) *TxBondDialogController {
	return &TxBondDialogController{view: view, model: model}
}

func (c *TxBondDialogController) Show() {
	c.applyDefaults()
	c.populateCombos()

	gtkutil.DropDownOnChanged(c.view.SenderDrop, c.onSenderChanged)
	gtkutil.ComboBoxOnChanged(c.view.ReceiverCombo, c.onReceiverChanged)
	gtkutil.EntryOnChanged(c.view.FeeEntry, c.onFeeChanged)

	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	c.view.SenderDrop.SetSelected(0)
	c.onSenderChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func (c *TxBondDialogController) applyDefaults() {
	setDefaultFee(c.model, c.view.FeeEntry)
}

func (c *TxBondDialogController) populateCombos() {
	gtkutil.DropDownFromAddressList(c.view.SenderDrop, c.model.ListAccountAddresses())
	gtkutil.ComboBoxFromAddressList(c.view.ReceiverCombo, c.model.ListValidatorAddresses())
}

func (c *TxBondDialogController) onSenderChanged() {
	sender := gtkutil.DropDownGetSelectedText(c.view.SenderDrop)
	if info := c.model.AddressInfo(sender); info != nil && info.Label != "" {
		setHint(c.view.SenderHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHint(c.view.SenderHint, "")
	}

	bal, err := c.model.Balance(sender)
	if err == nil {
		setHint(c.view.AmountHint, fmt.Sprintf("Account Balance: %s", bal))
	} else {
		setHint(c.view.AmountHint, "")
	}
}

func (c *TxBondDialogController) onReceiverChanged() {
	hint := ""

	receiver := c.view.ReceiverCombo.ActiveText()
	stake, err := c.model.Stake(receiver)
	if err == nil {
		hint = fmt.Sprintf("stake: %s", stake)
	}
	if info := c.model.AddressInfo(receiver); info != nil && info.Label != "" {
		if hint != "" {
			hint += ", "
		}
		hint += "label: " + info.Label
	}
	setHint(c.view.ReceiverHint, hint)
}

func (c *TxBondDialogController) onFeeChanged() {
	_ = payload.TypeBond
	setHint(c.view.FeeHint, "")
}

func (c *TxBondDialogController) onSend() {
	sender := gtkutil.DropDownGetSelectedText(c.view.SenderDrop)
	receiver := gtkutil.ComboBoxGetSelectedText(c.view.ReceiverCombo)
	publicKey := gtkutil.EntryGetText(c.view.PublicKeyEntry)
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

	trx, err := c.model.MakeBondTx(sender, receiver, publicKey, amt, fee, memo)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:   Bond
From:   %s
To:     %s
Amount: %s
Fee:    %s
Memo:   %s
</tt>

You are going to sign and broadcast this transaction.
<b>⚠️ This action cannot be undone.</b>
Do you want to continue with this transaction?`, sender, receiver, amt, trx.Fee(), trx.Memo())

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxBondDialogController) onCancel() {
	c.view.Window.Close()
}
