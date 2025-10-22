//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/version"
)

var (
	//go:embed assets/ui/dialog_about.ui
	uiAboutDialog []byte

	//go:embed assets/images/logo.png
	pactusLogo []byte
)

func aboutDialog() *gtk.AboutDialog {
	builder, err := gtk.BuilderNewFromString(string(uiAboutDialog))
	fatalErrorCheck(err)

	dlg := getAboutDialogObj(builder, "id_dialog_about")

	pxLogo, err := gdk.PixbufNewFromBytesOnly(pactusLogo)
	if err != nil {
		terminal.PrintErrorMsgf("Failed to load Logo Pixbuf: %v", err)
	} else {
		dlg.SetLogo(pxLogo)
	}

	dlg.SetVersion(version.NodeVersion().StringWithAlias())

	return dlg
}
