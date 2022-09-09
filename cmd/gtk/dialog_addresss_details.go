//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_address_details.ui
var uiAddressDetailsDialog []byte

func showAddressDetails(wallet *wallet.Wallet, addr string) {
	builder, err := gtk.BuilderNewFromString(string(uiAddressDetailsDialog))
	fatalErrorCheck(err)

	info := wallet.AddressInfo(addr)
	if info == nil {
		showErrorDialog(nil, "address not found")
		return
	}

	dlg := getDialogObj(builder, "id_dialog_address_details")
	addressEntry := getEntryObj(builder, "id_entry_address")
	pubKeyEntry := getEntryObj(builder, "id_entry_public_key")
	pathEntry := getEntryObj(builder, "id_entry_path")

	addressEntry.SetText(info.Address)
	pubKeyEntry.SetText(info.Pub.String())
	pathEntry.SetText(info.Path.String())

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
