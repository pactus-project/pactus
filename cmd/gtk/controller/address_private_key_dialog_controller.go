//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type AddressPrivateKeyDialogController struct {
	view  *view.AddressPrivateKeyDialogView
	model *model.WalletModel
}

func NewAddressPrivateKeyDialogController(
	view *view.AddressPrivateKeyDialogView,
	model *model.WalletModel,
) *AddressPrivateKeyDialogController {
	return &AddressPrivateKeyDialogController{view: view, model: model}
}

func (c *AddressPrivateKeyDialogController) Run(addr string) {
	password, ok := PasswordProvider(c.model)
	if !ok {
		return
	}

	prv, err := c.model.PrivateKey(password, addr)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}

	c.view.AddressEntry.SetText(addr)
	c.view.PrvKeyEntry.SetText(prv.String())

	c.view.ConnectSignals(map[string]any{
		"on_close": func() { c.view.Dialog.Close() },
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
