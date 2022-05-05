//go:build gtk

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

func showAboutGTKDialog(parent gtk.IWindow) {
	builder, err := gtk.BuilderNewFromString(string(uiAboutGtkDialog))
	errorCheck(parent, err)

	dlg := getAboutDialogObj(builder, "id_dialog_about_gtk")

	pixbuf, err := gdk.PixbufNewFromDataOnly(imgGtkIcon)
	errorCheck(parent, err)

	dlg.SetLogo(pixbuf)

	dlg.SetModal(true)

	dlg.Show()
}
