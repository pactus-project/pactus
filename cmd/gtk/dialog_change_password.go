//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"

	"github.com/zarbchain/zarb-go/wallet"
)

//go:embed assets/ui/dialog_change_password.ui
var uiChangePasswordDialog []byte

func changePassword(wallet *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiChangePasswordDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_change_password")
	oldPasswordEntry := getEntryObj(builder, "id_entry_old_password")
	oldPasswordLabel := getLabelObj(builder, "id_label_old_password")
	newPasswordEntry := getEntryObj(builder, "id_entry_new_password")
	repeatPasswordEntry := getEntryObj(builder, "id_entry_repeat_password")

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	if !wallet.IsEncrypted() {
		oldPasswordEntry.SetVisible(false)
		oldPasswordLabel.SetVisible(false)
	}

	onOk := func() {
		oldPassword, err := oldPasswordEntry.GetText()
		fatalErrorCheck(err)

		newPassword, err := newPasswordEntry.GetText()
		fatalErrorCheck(err)

		repeatPassword, err := repeatPasswordEntry.GetText()
		fatalErrorCheck(err)

		if newPassword != repeatPassword {
			showWarningDialog(dlg, "Passwords do not match")
			return
		}

		err = wallet.UpdatePassword(oldPassword, newPassword)
		if err != nil {
			showErrorDialog(dlg, err.Error())
			return
		}

		err = wallet.Save()
		if err != nil {
			showErrorDialog(dlg, err.Error())
			return
		}

		dlg.Close()
	}

	onCancel := func() {
		dlg.Close()
	}

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	dlg.Run()
}
