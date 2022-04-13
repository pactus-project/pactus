package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/wallet"
)

//go:embed ui/password_dialog.ui
var uiPasswordDialog []byte

func showInfoDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", msg)
	dlg.Run()
	dlg.Destroy()
}

func showErrorDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "%s", msg)
	dlg.Run()
	dlg.Destroy()
}

func getWalletPassword(parent *gtk.Widget, wallet *wallet.Wallet) (string, bool) {
	password := ""
	if !wallet.IsEncrypted() {
		return password, true
	}

	builder, err := gtk.BuilderNewFromString(string(uiPasswordDialog))
	errorCheck(err)

	// Get the object with the id of "password_dialog".
	obj, err := builder.GetObject("password_dialog")
	errorCheck(err)

	// Verify that the object is a pointer to a gtk.Dialog.
	dlg, err := isDialog(obj)
	errorCheck(err)

	onOk := func() {
		obj, err := builder.GetObject("password_entry")
		errorCheck(err)

		passwordEntry, err := isEntry(obj)
		errorCheck(err)

		password, err = passwordEntry.GetText()
		errorCheck(err)

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

	// Show the dialog
	code := dlg.Run()

	return password, code == gtk.RESPONSE_ACCEPT
}
