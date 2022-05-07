//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"

	"github.com/zarbchain/zarb-go/wallet"
)

//go:embed assets/ui/dialog_password.ui
var uiPasswordDialog []byte

func getWalletPassword(parent gtk.IWindow, wallet *wallet.Wallet) (string, bool) {
	password := ""
	if !wallet.IsEncrypted() {
		return password, true
	}

	builder, err := gtk.BuilderNewFromString(string(uiPasswordDialog))
	errorCheck(parent, err)

	dlg := getDialogObj(builder, "id_dialog_password")
	passwordEntry := getEntryObj(builder, "id_entry_password")

	password, err = passwordEntry.GetText()
	errorCheck(parent, err)

	ok := false
	onOk := func() {
		password, err = passwordEntry.GetText()
		errorCheck(parent, err)

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

	dlg.Run()

	return password, ok
}
