//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletPasswordDialogView struct {
	ViewBuilder

	Window *gtk.Window

	PasswordEntry *gtk.Entry
	ButtonOK      *gtk.Button
	ButtonCancel  *gtk.Button
}

func NewWalletPasswordDialogView() *WalletPasswordDialogView {
	builder := NewViewBuilder(assets.WalletPasswordDialogUI)

	view := &WalletPasswordDialogView{
		ViewBuilder:   builder,
		Window:        builder.GetWindowObj("id_dialog_wallet_password"),
		PasswordEntry: builder.GetEntryObj("id_entry_password"),
		ButtonOK:      builder.GetButtonObj("id_button_ok"),
		ButtonCancel:  builder.GetButtonObj("id_button_cancel"),
	}

	gtkutil.UpdateOKButton(view.Window, view.ButtonOK)
	gtkutil.UpdateCancelButton(view.ButtonCancel)

	gtkutil.DialogSetup(view.Window)

	return view
}
