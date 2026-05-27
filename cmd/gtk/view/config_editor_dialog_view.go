//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type ConfigEditorDialogView struct {
	ViewBuilder

	Window *gtk.Window

	TextView      *gtk.TextView
	ButtonSave    *gtk.Button
	ButtonRestore *gtk.Button
}

func NewConfigEditorDialogView() *ConfigEditorDialogView {
	builder := NewViewBuilder(assets.ConfigEditorDialogUI)

	view := &ConfigEditorDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_config_editor"),

		TextView:      builder.GetTextViewObj("id_textview_config"),
		ButtonSave:    builder.GetButtonObj("id_button_save_config"),
		ButtonRestore: builder.GetButtonObj("id_button_restore_default"),
	}

	gtkutil.AddImageToButton(view.ButtonSave, assets.IconSave16)
	gtkutil.AddImageToButton(view.ButtonRestore, assets.IconRefresh16)

	return view
}
