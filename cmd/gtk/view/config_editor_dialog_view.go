//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type ConfigEditorDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	TextView      *gtk.TextView
	ButtonSave    *gtk.Button
	ButtonRestore *gtk.Button
}

func NewConfigEditorDialogView() *ConfigEditorDialogView {
	builder := NewViewBuilder(assets.ConfigEditorDialogUI)

	view := &ConfigEditorDialogView{
		ViewBuilder: builder,
		Dialog:      builder.GetDialogObj("id_dialog_config_editor"),

		TextView:      builder.GetTextViewObj("id_textview_config"),
		ButtonSave:    builder.GetButtonObj("id_button_save_config"),
		ButtonRestore: builder.GetButtonObj("id_button_restore_default"),
	}

	view.ButtonSave.SetImage(gtkutil.ImageFromPixbuf(assets.IconSavePixbuf16))
	view.ButtonRestore.SetImage(gtkutil.ImageFromPixbuf(assets.IconRefreshPixbuf16))

	return view
}
