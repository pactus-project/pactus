//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/dialog_about.ui
var uiAboutDialog []byte

func showAboutDialog(parent gtk.IWindow) {
	builder, err := gtk.BuilderNewFromString(string(uiAboutDialog))
	errorCheck(parent, err)

	dlg := getAboutDialogObj(builder, "id_dialog_about")

	dlg.SetModal(true)

	dlg.Show()
}
