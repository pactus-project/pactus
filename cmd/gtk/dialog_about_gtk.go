//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

var (
	//go:embed assets/ui/dialog_about_gtk.ui
	uiAboutGtkDialog []byte

	//go:embed assets/images/gtk.svg
	imgGtkIcon []byte
)

func showAboutGTKDialog() {
	builder, err := gtk.BuilderNewFromString(string(uiAboutGtkDialog))
	fatalErrorCheck(err)

	dlg := getAboutDialogObj(builder, "id_dialog_about_gtk")

	pixbuf, err := gdk.PixbufNewFromDataOnly(imgGtkIcon)
	fatalErrorCheck(err)

	dlg.SetLogo(pixbuf)

	dlg.SetModal(true)

	dlg.Show()
}
