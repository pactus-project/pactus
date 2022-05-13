//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/dialog_label.ui
var uiLabelDialog []byte

func getAddressLabel(oldLabel string) (string, bool) {
	builder, err := gtk.BuilderNewFromString(string(uiLabelDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_label")
	labelEntry := getEntryObj(builder, "id_entry_label")
	labelEntry.SetText(oldLabel)

	newLabel := ""
	ok := false
	onOk := func() {
		newLabel, err = labelEntry.GetText()
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

	return newLabel, ok
}
