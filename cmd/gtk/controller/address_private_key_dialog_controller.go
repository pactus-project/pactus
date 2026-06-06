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

func (c *AddressPrivateKeyDialogController) Show(addr string) {
	PasswordProvider(c.model, func(password string, ok bool) {
		if !ok {
			return
		}

		prv, err := c.model.PrivateKey(password, addr)
		if err != nil {
			gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

			return
		}

		c.view.AddressEntry.SetText(addr)
		c.view.PrvKeyEntry.SetText(prv.String())

		onClose := func() {
			c.view.Window.Close()
		}

		gtkutil.ConnectButtonSignal(c.view.ButtonClose, onClose)

		gtkutil.ShowModalWindow(c.view.Window)
	})
}
