//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

func NewAboutDialog() *gtk.AboutDialog {
	builder := NewViewBuilder(assets.DialogAboutUI)
	dlg := builder.GetAboutDialogObj("id_dialog_about")

	pic := gtkutil.NewScaledPictureFromTexture(assets.ImagePactusLogoTexture, 128, 128)
	dlg.SetLogo(pic.Paintable())

	gtkutil.DisableLabelSelection(dlg)

	return dlg
}
