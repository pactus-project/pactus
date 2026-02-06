//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/crypto"
)

type TxUnbondDialogController struct {
	view  *view.TxUnbondDialogView
	model *model.WalletModel
}

func NewTxUnbondDialogController(
	view *view.TxUnbondDialogView,
	model *model.WalletModel,
) *TxUnbondDialogController {
	return &TxUnbondDialogController{view: view, model: model}
}

func (c *TxUnbondDialogController) Run() {
	for _, ai := range c.model.ListAddresses(crypto.AddressTypeValidator) {
		c.view.ValidatorCombo.Append(ai.Address, ai.Address)
	}

	c.view.ConnectSignals(map[string]any{
		"on_validator_changed": c.onValidatorChanged,
		"on_send":              c.onSend,
		"on_cancel":            c.onCancel,
	})

	c.onValidatorChanged()
	gtkutil.RunDialog(c.view.Dialog)
}

func (c *TxUnbondDialogController) onValidatorChanged() {
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

func (c *TxUnbondDialogController) onSend() {
	validatorEntry, _ := c.view.ValidatorCombo.GetEntry()
	validatorAddr := gtkutil.GetEntryText(validatorEntry)
	memo := gtkutil.GetEntryText(c.view.MemoEntry)

	trx, err := c.model.MakeUnbondTx(validatorAddr, memo)
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
Do you want to continue with this transaction?`, validatorAddr, trx.Fee(), trx.Memo())

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

func (c *TxUnbondDialogController) onCancel() {
	c.view.Dialog.Close()
}
