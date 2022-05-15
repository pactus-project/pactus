//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"

	"github.com/zarbchain/zarb-go/wallet"
)

//go:embed assets/ui/dialog_password.ui
var uiPasswordDialog []byte

func getWalletPassword(wallet *wallet.Wallet) (string, bool) {
	password := ""
	if !wallet.IsEncrypted() {
		return password, true
	}

	builder, err := gtk.BuilderNewFromString(string(uiPasswordDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_password")
	passwordEntry := getEntryObj(builder, "id_entry_password")

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
	signals := map[string]interface{}{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	dlg.Run()

	return password, ok
}
