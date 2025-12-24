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

type TxBondModel interface {
	WalletInfo() (types.WalletInfo, error)
	ListAddresses(opts ...wallet.ListAddressOption) []types.AddressInfo
	AddressInfo(addr string) *types.AddressInfo
	Balance(addr string) (amount.Amount, error)
	Stake(addr string) (amount.Amount, error)

	MakeBondTx(sender, receiver, publicKey string, amt amount.Amount, opts ...wallet.TxOption) (*tx.Tx, error)
	SignTransaction(password string, trx *tx.Tx) error
	BroadcastTransaction(trx *tx.Tx) (string, error)
}

type TxBondDialogController struct {
	view   *view.TxBondDialogView
	model  TxBondModel
	getPwd PasswordProvider
}

func NewTxBondDialogController(
	view *view.TxBondDialogView,
	model TxBondModel,
	getPwd PasswordProvider,
) *TxBondDialogController {
	return &TxBondDialogController{view: view, model: model, getPwd: getPwd}
}

func setHint(lbl *gtk.Label, hint string) {
	if hint == "" {
		lbl.SetMarkup("")

		return
	}
	lbl.SetMarkup(gtkutil.SmallGray(hint))
}

func (c *TxBondDialogController) Run() {
	if info, err := c.model.WalletInfo(); err == nil {
		c.view.FeeEntry.SetText(fmt.Sprintf("%g", info.DefaultFee.ToPAC()))
	}

	for _, ai := range c.model.ListAddresses(wallet.OnlyAccountAddresses()) {
		c.view.SenderCombo.Append(ai.Address, ai.Address)
	}
	for _, vi := range c.model.ListAddresses(wallet.OnlyValidatorAddresses()) {
		c.view.ReceiverCombo.Append(vi.Address, vi.Address)
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

func (c *TxBondDialogController) onSenderChanged() {
	sender := c.view.SenderCombo.GetActiveID()
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
	receiverEntry, _ := c.view.ReceiverCombo.GetEntry()
	receiver := gtkutil.GetEntryText(receiverEntry)

	stake, err := c.model.Stake(receiver)
	hint := ""
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
	sender := c.view.SenderCombo.GetActiveID()
	receiverEntry, _ := c.view.ReceiverCombo.GetEntry()
	receiver := gtkutil.GetEntryText(receiverEntry)
	publicKey := gtkutil.GetEntryText(c.view.PublicKeyEntry)
	amountStr := gtkutil.GetEntryText(c.view.AmountEntry)
	memo := gtkutil.GetEntryText(c.view.MemoEntry)

	amt, err := amount.FromString(amountStr)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	feeStr := gtkutil.GetEntryText(c.view.FeeEntry)
	opts := []wallet.TxOption{wallet.OptionMemo(memo), wallet.OptionFee(feeStr)}

	trx, err := c.model.MakeBondTx(sender, receiver, publicKey, amt, opts...)
	if err != nil {
		gtkutil.ShowError(err)

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

func (c *TxBondDialogController) onCancel() {
	c.view.Dialog.Close()
}
