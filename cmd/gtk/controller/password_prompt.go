//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

// Controller represents a UI controller that can be executed.
// All controllers should have a Run method to execute their primary function.

func PasswordProvider(model *model.WalletModel) (string, bool) {
	if model == nil || !model.IsEncrypted() {
		return "", true
	}

	return PromptWalletPassword()
}

// PromptWalletPassword shows the wallet password dialog and returns the entered password.
// It is used by the boot sequence (before models are constructed).
func PromptWalletPassword() (string, bool) {
	view := view.NewWalletPasswordDialogView()
	ctrl := NewWalletPasswordDialogController(view)

	return ctrl.Run()
}
