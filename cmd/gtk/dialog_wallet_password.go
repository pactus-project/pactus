//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_wallet_password.ui
var uiPasswordDialog []byte

func getWalletPassword(wlt *wallet.Wallet) (string, bool) {
	password := ""
	if !wlt.IsEncrypted() {
		return password, true
	}

	builder, err := gtk.BuilderNewFromString(string(uiPasswordDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_wallet_password")
	passwordEntry := getEntryObj(builder, "id_entry_password")

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	ok := false
	onOk := func() {
		password, err = passwordEntry.GetText()
		fatalErrorCheck(err)

		ok = true
		dlg.Close()
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

	return password, ok
}
