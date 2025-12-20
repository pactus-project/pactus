//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type WalletChangePasswordModel interface {
	IsEncrypted() bool
	UpdatePassword(oldPassword, newPassword string) error
}

type WalletChangePasswordDialogController struct {
	view *view.WalletChangePasswordDialogView
}

func NewWalletChangePasswordDialogController(
	view *view.WalletChangePasswordDialogView,
) *WalletChangePasswordDialogController {
	return &WalletChangePasswordDialogController{view: view}
}

func (c *WalletChangePasswordDialogController) Run(model WalletChangePasswordModel, afterSave func()) {
	if !model.IsEncrypted() {
		c.view.OldPasswordEntry.SetVisible(false)
		c.view.OldPasswordLabel.SetVisible(false)
	}

	onOk := func() {
		oldPassword := gtkutil.GetEntryText(c.view.OldPasswordEntry)
		newPassword := gtkutil.GetEntryText(c.view.NewPasswordEntry)
		repeatPassword := gtkutil.GetEntryText(c.view.RepeatEntry)

		if newPassword != repeatPassword {
			gtkutil.ShowWarningDialog(c.view.Dialog, "Passwords do not match")

			return
		}
		if err := model.UpdatePassword(oldPassword, newPassword); err != nil {
			gtkutil.ShowError(err)

			return
		}
		c.view.Dialog.Close()
		if afterSave != nil {
			afterSave()
		}
	}

	onCancel := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
