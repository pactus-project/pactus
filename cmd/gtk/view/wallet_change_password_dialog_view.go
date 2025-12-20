//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletChangePasswordDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	OldPasswordEntry *gtk.Entry
	OldPasswordLabel *gtk.Label
	NewPasswordEntry *gtk.Entry
	RepeatEntry      *gtk.Entry

	ButtonOK     *gtk.Button
	ButtonCancel *gtk.Button
}

func NewWalletChangePasswordDialogView() *WalletChangePasswordDialogView {
	builder := NewViewBuilder(assets.WalletChangePasswordDialogUI)

	view := &WalletChangePasswordDialogView{
		Dialog: builder.GetDialogObj("id_dialog_wallet_change_password"),

		OldPasswordEntry: builder.GetEntryObj("id_entry_old_password"),
		OldPasswordLabel: builder.GetLabelObj("id_label_old_password"),
		NewPasswordEntry: builder.GetEntryObj("id_entry_new_password"),
		RepeatEntry:      builder.GetEntryObj("id_entry_repeat_password"),

		ButtonOK:     builder.GetButtonObj("id_button_ok"),
		ButtonCancel: builder.GetButtonObj("id_button_cancel"),

		ViewBuilder: builder,
	}

	view.ButtonOK.SetImage(gtkutil.ImageFromPixbuf(assets.IconOkPixbuf16))
	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))

	return view
}
