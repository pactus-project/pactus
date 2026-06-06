//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
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

	gtkutil.DropDownOnChanged(c.view.ValidatorDrop, c.onValidatorChanged)
	gtkutil.ComboBoxOnChanged(c.view.ReceiverCombo, c.onReceiverChanged)
	gtkutil.EntryOnChanged(c.view.FeeEntry, c.onFeeChanged)

	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	c.view.ValidatorDrop.SetSelected(0)
	c.onValidatorChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func (c *TxWithdrawDialogController) applyDefaults() {
	if info, err := c.model.WalletInfo(); err == nil {
		c.view.FeeEntry.SetText(fmt.Sprintf("%g", info.DefaultFee.ToPAC()))
	}
}

func (c *TxWithdrawDialogController) populateCombos() {
	gtkutil.DropDownFromAddressList(c.view.ValidatorDrop, c.model.ListValidatorAddresses())
	gtkutil.ComboBoxFromAddressList(c.view.ReceiverCombo, c.model.ListAccountAddresses())
}

func (c *TxWithdrawDialogController) onValidatorChanged() {
	hint := ""
	sender := gtkutil.DropDownGetSelectedText(c.view.ValidatorDrop)
	stake, err := c.model.Stake(sender)
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
	sender := gtkutil.DropDownGetSelectedText(c.view.ValidatorDrop)
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

	trx, err := c.model.MakeWithdrawTx(sender, receiver, amt, fee, memo)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

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

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxWithdrawDialogController) onCancel() {
	c.view.Window.Close()
}
