package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/wallet"
)

//go:embed assets/ui/dialog_password.ui
var uiPasswordDialog []byte

func getWalletPassword(parent gtk.IWidget, wallet *wallet.Wallet) (string, bool) {
	password := ""
	if !wallet.IsEncrypted() {
		return password, true
	}

	builder, err := gtk.BuilderNewFromString(string(uiPasswordDialog))
	errorCheck(err)

	dlg := getAboutDialogObj(builder, "password_dialog")
	passwordEntry := getEntryObj(builder, "password_entry")

	password, err = passwordEntry.GetText()
	errorCheck(err)

	ok := false
	onOk := func() {
		password, err = passwordEntry.GetText()
		errorCheck(err)

		ok = true
		dlg.Close()
		dlg.Destroy()
	}

	onCancel := func() {
		dlg.Close()
		dlg.Destroy()
	}

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	// Show the dialog
	dlg.Run()

	return password, ok
}
