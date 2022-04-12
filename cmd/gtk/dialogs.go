package main

import (
	"github.com/gotk3/gotk3/gtk"
)

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

// func showErrDialog(parent widgets.QWidget_ITF, msg string) {
// 	widgets.QMessageBox_Critical(
// 		parent,
// 		"error",
// 		msg,
// 		widgets.QMessageBox__Ok,
// 		widgets.QMessageBox__NoButton)
// }

// func getWalletPassword(parent widgets.QWidget_ITF, wallet *wallet.Wallet) (string, bool) {
// 	if !wallet.IsEncrypted() {
// 		return "", true
// 	}

// 	// ok := new(bool)
// 	// password := widgets.QInputDialog_GetText(parent,
// 	// 	"Enter password",
// 	// 	"Please enter your wallet password:",
// 	// 	widgets.QLineEdit__Password,
// 	// 	"",
// 	// 	ok,
// 	// 	core.Qt__Dialog,
// 	// 	core.Qt__ImhNone)

// 	ok := new(bool)
// 	password := new(string)

// 	dialog := widgets.NewQInputDialog(parent, core.Qt__Tool)
// 	dialog.SetWindowTitle("Enter password")
// 	dialog.SetLabelText("Please enter your wallet password:")
// 	dialog.SetTextEchoMode(widgets.QLineEdit__Password)
// 	dialog.SetInputMethodHints(core.Qt__ImhNone)

// 	dialog.ConnectAccept(func() {
// 		*password = dialog.TextValue()
// 		*ok = true
// 		dialog.AcceptDefault()
// 	})

// 	dialog.ConnectReject(func() {
// 		*ok = false
// 		dialog.RejectDefault()
// 	})

// 	dialog.Exec()

// 	return *password, *ok
// }
