//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

func NewAboutDialog() *gtk.AboutDialog {
	builder := NewViewBuilder(assets.DialogAboutUI)
	dlg := builder.GetAboutDialogObj("id_dialog_about")

	dlg.SetLogo(assets.ImagePactusLogoPixbuf)

	return dlg
}
