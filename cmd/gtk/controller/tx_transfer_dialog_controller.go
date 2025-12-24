//go:build gtk

package controller

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/types"
)

type TxTransferModel interface {
	WalletInfo() (types.WalletInfo, error)
	ListAddresses(opts ...wallet.ListAddressOption) []types.AddressInfo
	AddressInfo(addr string) *types.AddressInfo
	Balance(addr string) (amount.Amount, error)

	MakeTransferTx(sender, receiver string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	SignTransaction(password string, trx *tx.Tx) error
	BroadcastTransaction(trx *tx.Tx) (string, error)
}

type TxTransferDialogController struct {
	view   *view.TxTransferDialogView
	model  TxTransferModel
	getPwd PasswordProvider
}

func NewTxTransferDialogController(
	view *view.TxTransferDialogView,
	model TxTransferModel,
	getPwd PasswordProvider,
) *TxTransferDialogController {
	return &TxTransferDialogController{view: view, model: model, getPwd: getPwd}
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
	for _, ai := range c.model.ListAddresses(wallet.OnlyAccountAddresses()) {
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
	gtkutil.RunDialog(c.view.Dialog)
}

func (c *TxTransferDialogController) onSenderChanged() {
	sender := c.view.SenderCombo.GetActiveID()
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
	sender := c.view.SenderCombo.GetActiveID()
	receiver := gtkutil.GetEntryText(c.view.ReceiverEntry)
	amountStr := gtkutil.GetEntryText(c.view.AmountEntry)
	memo := gtkutil.GetEntryText(c.view.MemoEntry)

	amt, err := amount.FromString(amountStr)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	feeStr := gtkutil.GetEntryText(c.view.FeeEntry)
	opts := []wallet.TxOption{wallet.OptionMemo(memo), wallet.OptionFee(feeStr)}

	trx, err := c.model.MakeTransferTx(sender, receiver, amt, opts...)
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

	if !gtkutil.ShowQuestionDialog(c.view.Dialog, msg) {
		return
	}

	password, ok := c.getPwd()
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

func (c *TxTransferDialogController) onCancel() {
	c.view.Dialog.Close()
}
