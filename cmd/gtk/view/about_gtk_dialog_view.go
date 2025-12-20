//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

func NewAboutGTKDialog() *gtk.AboutDialog {
	builder := NewViewBuilder(assets.DialogAboutGTKUI)
	dlg := builder.GetAboutDialogObj("id_dialog_about_gtk")
	dlg.SetLogo(assets.ImageGTKLogoPixbuf)

	return dlg
}
