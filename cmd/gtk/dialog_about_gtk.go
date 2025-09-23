//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
)

var (
	//go:embed assets/ui/dialog_about_gtk.ui
	uiAboutGtkDialog []byte

	//go:embed assets/images/gtk.png
	imageGtk []byte
)

func showAboutGTKDialog() {
	builder, err := gtk.BuilderNewFromString(string(uiAboutGtkDialog))
	fatalErrorCheck(err)

	dlg := getAboutDialogObj(builder, "id_dialog_about_gtk")

	pixbuf, err := gdk.PixbufNewFromBytesOnly(imageGtk)
	if err != nil {
		// Handle error gracefully instead of fatal
		cmd.PrintErrorMsgf("Failed to load Logo Pixbuf: %v", err)
	} else {
		dlg.SetLogo(pixbuf)
	}

	dlg.SetModal(true)

	runDialog(&dlg.Dialog)
}
