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
	view  *view.WalletChangePasswordDialogView
	model WalletChangePasswordModel
}

func NewWalletChangePasswordDialogController(
	view *view.WalletChangePasswordDialogView,
	model WalletChangePasswordModel,
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
			gtkutil.ShowWarningDialog(c.view.Dialog, "Passwords do not match")

			return
		}
		if err := c.model.UpdatePassword(oldPassword, newPassword); err != nil {
			gtkutil.ShowError(err)

			return
		}
		c.view.Dialog.Close()
	}

	onCancel := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	})

	c.view.Dialog.SetModal(true)
	gtkutil.RunDialog(c.view.Dialog)
}
