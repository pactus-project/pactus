//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

func PasswordProvider(model *model.WalletModel, onPassword func(string, bool)) {
	if model == nil || !model.IsEncrypted() {
		onPassword("", true)

		return
	}

	PromptWalletPassword(onPassword)
}

// PromptWalletPassword shows the wallet password dialog and returns the entered password.
// It is used by the boot sequence (before models are constructed).
func PromptWalletPassword(onPassword func(string, bool)) {
	view := view.NewWalletPasswordDialogView()
	ctrl := NewWalletPasswordDialogController(view)

	ctrl.Show(onPassword)
}
