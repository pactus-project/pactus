package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/dialog_about.ui
var uiAboutDialog []byte

func showAboutDialog(parent gtk.IWidget) {
	builder, err := gtk.BuilderNewFromString(string(uiAboutDialog))
	errorCheck(err)

	objDlg, err := builder.GetObject("id_dialog_about")
	errorCheck(err)

	dlg, err := isAboutDialog(objDlg)
	errorCheck(err)

	dlg.SetModal(true)
	// Show the dialog
	dlg.Show()
}
