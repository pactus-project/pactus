//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

func NewAboutGTKDialog() *gtk.AboutDialog {
	builder := NewViewBuilder(assets.DialogAboutGTKUI)
	dlg := builder.GetAboutDialogObj("id_dialog_about_gtk")

	logo := gtkutil.ResizeImage(assets.ImageGTKLogo, 96, 96)
	dlg.SetLogo(logo.Paintable())

	return dlg
}
