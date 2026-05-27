//go111:build gtk

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
		oldPassword := gtkutil.GetEntryText(c.view.OldPasswordEntry)
		newPassword := gtkutil.GetEntryText(c.view.NewPasswordEntry)
		repeatPassword := gtkutil.GetEntryText(c.view.RepeatEntry)

		if newPassword != repeatPassword {
			gtkutil.ShowWarningDialog(c.view.Window, "Passwords do not match")

			return
		}
		if err := c.model.UpdatePassword(oldPassword, newPassword); err != nil {
			gtkutil.ShowError(err)

			return
		}
		c.view.Window.Close()
	}

	onCancel := func() { c.view.Window.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	})

	gtkutil.ShowModalDialog(c.view.Window)
}
