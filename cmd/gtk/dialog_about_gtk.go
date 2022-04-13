package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/dialog_about_gtk.ui
var uiAboutGtkDialog []byte

//go:embed assets/images/gtk.svg
var imgGtkIcon []byte

func showAboutGTKDialog(parent gtk.IWidget) {
	builder, err := gtk.BuilderNewFromString(string(uiAboutGtkDialog))
	errorCheck(err)

	objDlg, err := builder.GetObject("id_dialog_about_gtk")
	errorCheck(err)

	dlg, err := isAboutDialog(objDlg)
	errorCheck(err)

	pixbuf, err := gdk.PixbufNewFromDataOnly(imgGtkIcon)
	errorCheck(err)

	dlg.SetLogo(pixbuf)

	dlg.SetModal(true)
	// Show the dialog
	dlg.Show()
}
