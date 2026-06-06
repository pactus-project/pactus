//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletChangePasswordDialogView struct {
	ViewBuilder

	Window *gtk.Window

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
		ViewBuilder:      builder,
		Window:           builder.GetWindowObj("id_dialog_wallet_change_password"),
		OldPasswordEntry: builder.GetEntryObj("id_entry_old_password"),
		OldPasswordLabel: builder.GetLabelObj("id_label_old_password"),
		NewPasswordEntry: builder.GetEntryObj("id_entry_new_password"),
		RepeatEntry:      builder.GetEntryObj("id_entry_repeat_password"),
		ButtonOK:         builder.GetButtonObj("id_button_ok"),
		ButtonCancel:     builder.GetButtonObj("id_button_cancel"),
	}

	gtkutil.UpdateOKButton(view.Window, view.ButtonOK)
	gtkutil.UpdateCancelButton(view.ButtonCancel)

	gtkutil.DialogSetup(view.Window)

	return view
}
