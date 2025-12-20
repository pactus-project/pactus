//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletPasswordDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	PasswordEntry *gtk.Entry
	ButtonOK      *gtk.Button
	ButtonCancel  *gtk.Button
}

func NewWalletPasswordDialogView() *WalletPasswordDialogView {
	builder := NewViewBuilder(assets.WalletPasswordDialogUI)

	view := &WalletPasswordDialogView{
		ViewBuilder: builder,
		Dialog:      builder.GetDialogObj("id_dialog_wallet_password"),

		PasswordEntry: builder.GetEntryObj("id_entry_password"),
		ButtonOK:      builder.GetButtonObj("id_button_ok"),
		ButtonCancel:  builder.GetButtonObj("id_button_cancel"),
	}

	view.ButtonOK.SetImage(gtkutil.ImageFromPixbuf(assets.IconOkPixbuf16))
	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))

	return view
}
