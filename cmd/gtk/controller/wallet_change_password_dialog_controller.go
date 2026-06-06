//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type WalletChangePasswordDialogController struct {
	view  *view.WalletChangePasswordDialogView
	model *model.WalletModel
}

func NewWalletChangePasswordDialogController(
	view *view.WalletChangePasswordDialogView,
	model *model.WalletModel,
) *WalletChangePasswordDialogController {
	return &WalletChangePasswordDialogController{view: view, model: model}
}

func (c *WalletChangePasswordDialogController) Run() {
	if !c.model.IsEncrypted() {
		c.view.OldPasswordEntry.SetVisible(false)
		c.view.OldPasswordLabel.SetVisible(false)
	}

	onOk := func() {
		oldPassword := gtkutil.EntryGetText(c.view.OldPasswordEntry)
		newPassword := gtkutil.EntryGetText(c.view.NewPasswordEntry)
		repeatPassword := gtkutil.EntryGetText(c.view.RepeatEntry)

		if newPassword != repeatPassword {
			gtkutil.ShowWarningDialog(c.view.Window, "Passwords do not match", nil)

			return
		}
		if err := c.model.UpdatePassword(oldPassword, newPassword); err != nil {
			gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

			return
		}
		c.view.Window.Close()
	}

	onCancel := func() { c.view.Window.Close() }

	gtkutil.ConnectButtonSignal(c.view.ButtonOK, onOk)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, onCancel)

	gtkutil.ShowModalWindow(c.view.Window)
}
