//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_address_private_key.ui
var uiAddressPrivateKeyDialog []byte

func showAddressPrivateKey(wallet *wallet.Wallet, addr string) {
	builder, err := gtk.BuilderNewFromString(string(uiAddressPrivateKeyDialog))
	fatalErrorCheck(err)

	password, ok := getWalletPassword(wallet)
	if !ok {
		return
	}

	prv, err := wallet.PrivateKey(password, addr)
	errorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_address_private_key")
	addressEntry := getEntryObj(builder, "id_entry_address")
	prvKeyEntry := getEntryObj(builder, "id_entry_private_key")

	addressEntry.SetText(addr)
	prvKeyEntry.SetText(prv.String())

	getButtonObj(builder, "id_button_close").SetImage(CloseIcon())

	onClose := func() {
		dlg.Close()
	}

	signals := map[string]interface{}{
		"on_close": onClose,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	dlg.Run()
}
