//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type WalletSeedDialogController struct {
	view  *view.WalletSeedDialogView
	model *model.WalletModel
}

func NewWalletSeedDialogController(
	view *view.WalletSeedDialogView,
	model *model.WalletModel,
) *WalletSeedDialogController {
	return &WalletSeedDialogController{view: view, model: model}
}

func (c *WalletSeedDialogController) Run() {
	PasswordProvider(c.model, func(password string, ok bool) {
		if !ok {
			return
		}

		seed, err := c.model.Mnemonic(password)
		if err != nil {
			gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

			return
		}
		gtkutil.SetTextViewContent(c.view.TextView, seed)

		onClose := func() {
			c.view.Window.Close()
		}

		gtkutil.ConnectButtonSignal(c.view.ButtonClose, onClose)

		gtkutil.ShowModalWindow(c.view.Window)
	})
}
