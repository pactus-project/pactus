//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/genesis"
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

	gtkutil.ConnectButtonSignal(c.view.AmountMaxButton, c.onMax)
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
	gtkutil.DropDownFromAddressListWithAmount(c.view.SenderDrop, c.model.ListAccountAddresses(), c.model.Balance)
	gtkutil.ComboBoxFromAddressListWithAmount(c.view.ReceiverCombo, c.model.ListValidatorAddresses(), c.model.Stake)
}

func (c *TxBondDialogController) onSenderChanged() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	if info := c.model.AddressInfo(sender); info != nil && info.Label != "" {
		setHint(c.view.SenderHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHint(c.view.SenderHint, "")
	}

	bal, err := c.model.Balance(sender)
	if err == nil {
		setBalanceHint(c.view.AmountHint, bal)
	} else {
		setHint(c.view.AmountHint, "")
	}
}

func (c *TxBondDialogController) onReceiverChanged() {
	hint := ""

	receiver := gtkutil.ComboBoxGetSelectedAddress(c.view.ReceiverCombo)
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

// onMax fills in the whole balance of the sender, minus the fee, but not more
// than the stake the validator can still receive.
func (c *TxBondDialogController) onMax() {
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

	// A validator that is not on the blockchain yet has no stake.
	receiver := gtkutil.ComboBoxGetSelectedAddress(c.view.ReceiverCombo)
	stake, err := c.model.Stake(receiver)
	if err != nil {
		stake = 0
	}

	maxStake := genesis.DefaultGenesisParams().MaximumStake
	maxAmount, err := model.MaxBondAmount(balance, fee, stake, maxStake)
	setMaxAmount(c.view.Window, c.view.AmountEntry, maxAmount, err)
}

func (c *TxBondDialogController) onFeeChanged() {
	_ = payload.TypeBond
	setHint(c.view.FeeHint, "")
}

func (c *TxBondDialogController) onSend() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	receiver := gtkutil.ComboBoxGetSelectedAddress(c.view.ReceiverCombo)
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
<b>Transaction Details</b>
<tt>
Type:   Bond
From:   %s
To:     %s
Amount: %s
Fee:    %s
Memo:   %s
</tt>

You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`, sender, receiver, amt, trx.Fee(), trx.Memo())

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxBondDialogController) onCancel() {
	c.view.Window.Close()
}
