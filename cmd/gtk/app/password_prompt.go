//go:build gtk

package app

import (
	"github.com/pactus-project/pactus/cmd/gtk/controller"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

// PromptWalletPassword shows the wallet password dialog and returns the entered password.
// It is used by the boot sequence (before models are constructed).
func PromptWalletPassword() (string, bool) {
	v := view.NewWalletPasswordDialogView()

	ctrl := controller.NewWalletPasswordDialogController(v)

	pwd, ok := ctrl.Run()

	return pwd, ok
}
