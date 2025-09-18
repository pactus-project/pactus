//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_address_details.ui
var uiAddressDetailsDialog []byte

func showAddressDetails(wlt *wallet.Wallet, addr string) {
	builder, err := gtk.BuilderNewFromString(string(uiAddressDetailsDialog))
	fatalErrorCheck(err)

	info := wlt.AddressInfo(addr)
	if info == nil {
		showErrorDialog(nil, "address not found")

		return
	}

	dlg := getDialogObj(builder, "id_dialog_address_details")
	addressEntry := buildExtendedEntry(builder, "id_overlay_address")
	pubKeyEntry := buildExtendedEntry(builder, "id_overlay_public_key")
	pathEntry := getEntryObj(builder, "id_entry_path")

	addressEntry.SetText(info.Address)
	pubKeyEntry.SetText(info.PublicKey)
	pathEntry.SetText(info.Path)

	getButtonObj(builder, "id_button_close").SetImage(CloseIcon())

	onClose := func() {
		dlg.Close()
	}

	signals := map[string]any{
		"on_close": onClose,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	runDialog(dlg)
}
