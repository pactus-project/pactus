//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/dialog_address_label.ui
var uiAddressLabelDialog []byte

func getAddressLabel(oldLabel string) (string, bool) {
	builder, err := gtk.BuilderNewFromString(string(uiAddressLabelDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_address_label")
	labelEntry := getEntryObj(builder, "id_entry_label")
	labelEntry.SetText(oldLabel)

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	newLabel := ""
	ok := false
	onOk := func() {
		newLabel = getEntryText(labelEntry)

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

	return newLabel, ok
}
