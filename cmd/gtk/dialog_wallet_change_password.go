//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/dialog_wallet_change_password.ui
var uiWalletChangePasswordDialog []byte

func changePassword(wdgWallet *widgetWallet) {
	builder, err := gtk.BuilderNewFromString(string(uiWalletChangePasswordDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_wallet_change_password")
	oldPasswordEntry := getEntryObj(builder, "id_entry_old_password")
	oldPasswordLabel := getLabelObj(builder, "id_label_old_password")
	newPasswordEntry := getEntryObj(builder, "id_entry_new_password")
	repeatPasswordEntry := getEntryObj(builder, "id_entry_repeat_password")

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	if !wdgWallet.model.IsEncrypted() {
		oldPasswordEntry.SetVisible(false)
		oldPasswordLabel.SetVisible(false)
	}

	onOk := func() {
		oldPassword := getEntryText(oldPasswordEntry)

		newPassword := getEntryText(newPasswordEntry)

		repeatPassword := getEntryText(repeatPasswordEntry)

		if newPassword != repeatPassword {
			showWarningDialog(dlg, "Passwords do not match")

			return
		}

		err = wdgWallet.model.UpdatePassword(oldPassword, newPassword)
		if err != nil {
			showError(err)

			return
		}

		dlg.Close()

		wdgWallet.rebuild()
	}

	onCancel := func() {
		dlg.Close()
	}

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	runDialog(dlg)
}
