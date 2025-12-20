//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/vault"
)

type UnbondTxModel interface {
	AllValidatorAddresses() []vault.AddressInfo
	AddressInfo(addr string) *vault.AddressInfo
	Stake(addr string) (amount.Amount, error)

	MakeUnbondTx(validator string, opts ...wallet.TxOption) (*tx.Tx, error)
	SignTransaction(password string, trx *tx.Tx) error
	BroadcastTransaction(trx *tx.Tx) (string, error)
}

type UnbondTxController struct {
	view   *view.TxUnbondDialogView
	model  UnbondTxModel
	getPwd TxPasswordProvider
}

func NewUnbondTxController(
	view *view.TxUnbondDialogView,
	model UnbondTxModel,
	getPassword TxPasswordProvider,
) *UnbondTxController {
	return &UnbondTxController{view: view, model: model, getPwd: getPassword}
}

func (c *UnbondTxController) BindAndRun() {
	for _, ai := range c.model.AllValidatorAddresses() {
		c.view.ValidatorCombo.Append(ai.Address, ai.Address)
	}

	onValidatorChanged := func() {
		receiverEntry, _ := c.view.ValidatorCombo.GetEntry()
		validator := gtkutil.GetEntryText(receiverEntry)

		stake, err := c.model.Stake(validator)
		hint := ""
		if err == nil {
			hint = fmt.Sprintf("stake: %s", stake)
		}
		if info := c.model.AddressInfo(validator); info != nil && info.Label != "" {
			if hint != "" {
				hint += ", "
			}
			hint += "label: " + info.Label
		}
		if hint == "" {
			c.view.ValidatorHint.SetMarkup("")
		} else {
			c.view.ValidatorHint.SetMarkup(gtkutil.SmallGray(hint))
		}
	}

	onSend := func() {
		validatorEntry, _ := c.view.ValidatorCombo.GetEntry()
		validator := gtkutil.GetEntryText(validatorEntry)
		memo := gtkutil.GetEntryText(c.view.MemoEntry)

		trx, err := c.model.MakeUnbondTx(validator, wallet.OptionMemo(memo))
		if err != nil {
			gtkutil.ShowError(err)

			return
		}

		msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:     Unbond
Validator: %s
Fee:       %s
Memo:      %s
</tt>

You are going to sign and broadcast this transaction.
<b>⚠️ This action cannot be undone.</b>
Do you want to continue with this transaction?`, validator, trx.Fee(), trx.Memo())

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

	onCancel := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_validator_changed": onValidatorChanged,
		"on_send":              onSend,
		"on_cancel":            onCancel,
	})

	onValidatorChanged()
	gtkutil.RunDialog(c.view.Dialog)
}
