//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/crypto"
)

type AddressPrivateKeyModel interface {
	PrivateKey(password, addr string) (crypto.PrivateKey, error)
}

type AddressPrivateKeyDialogController struct {
	view   *view.AddressPrivateKeyDialogView
	model  AddressPrivateKeyModel
	getPwd TxPasswordProvider
}

func NewAddressPrivateKeyDialogController(
	view *view.AddressPrivateKeyDialogView,
	model AddressPrivateKeyModel,
	getPassword TxPasswordProvider,
) *AddressPrivateKeyDialogController {
	return &AddressPrivateKeyDialogController{view: view, model: model, getPwd: getPassword}
}

func (c *AddressPrivateKeyDialogController) Run(addr string) {
	password, ok := c.getPwd()
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
