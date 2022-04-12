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

func errorCheck(err error) {
	if err != nil {
		showErrorDialog(nil, err.Error())
		gtk.MainQuit()
	}
}

func getWalletPassword(parent *gtk.Widget, wallet *wallet.Wallet) (string, bool) {
	if !wallet.IsEncrypted() {
		return "", true
	}

	gtk.Init(nil)

	var ok = new(bool)
	var password = new(string)
	var dlg *gtk.Dialog

	builder, err := gtk.BuilderNewFromString(string(uiPasswordDialog))
	errorCheck(err)

	onOk := func() {
		obj, err := builder.GetObject("password_entry")
		errorCheck(err)

		passwordEntry, err := isEntry(obj)
		errorCheck(err)

		pwd, err := passwordEntry.GetText()
		errorCheck(err)

		*password = pwd
		*ok = true

		dlg.Close()
		// TODO
		gtk.MainQuit()
	}

	onCancel := func() {
		*ok = false
		dlg.Close()

		// TODO
		gtk.MainQuit()
	}

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	// Get the object with the id of "password_dialog".
	obj, err := builder.GetObject("password_dialog")
	errorCheck(err)

	// Verify that the object is a pointer to a gtk.Dialog.
	dlg, err = isDialog(obj)
	errorCheck(err)

	dlg.SetModal(true)

	// Show the dialog
	dlg.Show()

	gtk.Main()

	return *password, *ok
}
