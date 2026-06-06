//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletSeedDialogView struct {
	ViewBuilder

	Window *gtk.Window

	TextView    *gtk.TextView
	ButtonClose *gtk.Button
}

func NewWalletSeedDialogView() *WalletSeedDialogView {
	builder := NewViewBuilder(assets.WalletShowSeedDialogUI)

	view := &WalletSeedDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_wallet_show_seed"),

		TextView:    builder.GetTextViewObj("id_textview_seed"),
		ButtonClose: builder.GetButtonObj("id_button_close"),
	}

	gtkutil.UpdateCloseButton(view.ButtonClose)

	gtkutil.DialogSetup(view.Window)

	return view
}
