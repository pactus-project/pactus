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
	password, ok := PasswordProvider(c.model)
	if !ok {
		return
	}
	seed, err := c.model.Mnemonic(password)
	if err != nil {
		gtkutil.ShowError(err)

		return
	}
	gtkutil.SetTextViewContent(c.view.TextView, seed)
	c.view.ConnectSignals(map[string]any{
		"on_close": func() { c.view.Dialog.Close() },
	})
	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
