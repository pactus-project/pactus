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
	"github.com/pactus-project/pactus/wallet/vault"
)

type TransferTxModel interface {
	WalletInfo() (wallet.Info, error)
	AllAccountAddresses() []vault.AddressInfo
	AddressInfo(addr string) *vault.AddressInfo
	Balance(addr string) (amount.Amount, error)

	MakeTransferTx(sender, receiver string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	SignTransaction(password string, trx *tx.Tx) error
	BroadcastTransaction(trx *tx.Tx) (string, error)
}

type TxPasswordProvider func() (string, bool)

type TransferTxController struct {
	view   *view.TxTransferDialogView
	model  TransferTxModel
	getPwd TxPasswordProvider
}

func NewTransferTxController(
	view *view.TxTransferDialogView,
	model TransferTxModel,
	getPassword TxPasswordProvider,
) *TransferTxController {
	return &TransferTxController{view: view, model: model, getPwd: getPassword}
}

func setHintLabel(lbl *gtk.Label, hint string) {
	if hint == "" {
		lbl.SetMarkup("")

		return
	}
	lbl.SetMarkup(gtkutil.SmallGray(hint))
}

func (c *TransferTxController) BindAndRun() {
	// Defaults
	if info, err := c.model.WalletInfo(); err == nil {
		c.view.FeeEntry.SetText(fmt.Sprintf("%g", info.DefaultFee.ToPAC()))
	}

	// Fill sender accounts
	for _, ai := range c.model.AllAccountAddresses() {
		c.view.SenderCombo.Append(ai.Address, ai.Address)
	}
	c.view.SenderCombo.SetActive(0)

	onSenderChanged := func() {
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

	onReceiverChanged := func() {
		receiver := gtkutil.GetEntryText(c.view.ReceiverEntry)
		if info := c.model.AddressInfo(receiver); info != nil && info.Label != "" {
			setHintLabel(c.view.ReceiverHint, fmt.Sprintf("label: %s", info.Label))
		} else {
			setHintLabel(c.view.ReceiverHint, "")
		}
	}

	onFeeChanged := func() {
		// Placeholder (confirmation time estimation)
		_ = payload.TypeTransfer
		setHintLabel(c.view.FeeHint, "")
	}

	onSend := func() {
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
Type:   Transfer
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

	onCancel := func() {
		c.view.Dialog.Close()
	}

	c.view.ConnectSignals(map[string]any{
		"on_sender_changed":   onSenderChanged,
		"on_receiver_changed": onReceiverChanged,
		"on_fee_changed":      onFeeChanged,
		"on_send":             onSend,
		"on_cancel":           onCancel,
	})

	onSenderChanged()

	gtkutil.RunDialog(c.view.Dialog)
}
