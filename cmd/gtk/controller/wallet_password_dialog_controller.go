//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type WalletPasswordDialogController struct {
	view *view.WalletPasswordDialogView
}

func NewWalletPasswordDialogController(
	view *view.WalletPasswordDialogView,
) *WalletPasswordDialogController {
	return &WalletPasswordDialogController{view: view}
}

func (c *WalletPasswordDialogController) Run() (string, bool) {
	password := ""
	ok := false

	onOk := func() {
		password = gtkutil.GetEntryText(c.view.PasswordEntry)
		ok = true
		c.view.Dialog.Close()
	}
	andClose := func() {
		c.view.Dialog.Close()
	}

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": andClose,
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)

	return password, ok
}
