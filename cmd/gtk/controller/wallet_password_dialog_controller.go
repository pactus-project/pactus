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

func (c *WalletPasswordDialogController) Show(onPassword func(string, bool)) {
	password := ""
	ok := false

	onOk := func() {
		password = gtkutil.EntryGetText(c.view.PasswordEntry)
		ok = true

		c.view.Window.Close()
	}
	onCancel := func() {
		c.view.Window.Close()
	}

	// Window close button (X)
	c.view.Window.Connect("close-request", func() bool {
		onPassword(password, ok)

		return false
	})

	gtkutil.ConnectButtonSignal(c.view.ButtonOK, onOk)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, onCancel)

	gtkutil.ShowModalWindow(c.view.Window)
}
