//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
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

func (c *TxUnbondDialogController) Show() {
	c.applyDefaults()
	c.populateCombos()

	gtkutil.DropDownOnChanged(c.view.ValidatorDrop, c.onValidatorChanged)
	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	c.view.ValidatorDrop.SetSelected(0)
	c.onValidatorChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func (*TxUnbondDialogController) applyDefaults() {
}

func (c *TxUnbondDialogController) populateCombos() {
	gtkutil.DropDownFromAddressList(c.view.ValidatorDrop, c.model.ListValidatorAddresses())
}

func (c *TxUnbondDialogController) onValidatorChanged() {
	hint := ""
	validator := gtkutil.DropDownGetSelectedText(c.view.ValidatorDrop)
	stake, err := c.model.Stake(validator)
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
	validator := gtkutil.DropDownGetSelectedText(c.view.ValidatorDrop)
	memo := gtkutil.EntryGetText(c.view.MemoEntry)

	trx, err := c.model.MakeUnbondTx(validator, memo)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:      Unbond
Validator: %s
Fee:       %s
Memo:      %s
</tt>

You are going to sign and broadcast this transaction.
<b>⚠️ This action cannot be undone.</b>
Do you want to continue with this transaction?`, validator, trx.Fee(), trx.Memo())

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxUnbondDialogController) onCancel() {
	c.view.Window.Close()
}
