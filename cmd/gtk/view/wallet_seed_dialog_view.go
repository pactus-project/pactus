//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletSeedDialogView struct {
	ViewBuilder

	Dialog      *gtk.Dialog
	TextView    *gtk.TextView
	Image       *gtk.Image
	ButtonClose *gtk.Button
}

func NewWalletSeedDialogView() *WalletSeedDialogView {
	builder := NewViewBuilder(assets.WalletShowSeedDialogUI)

	view := &WalletSeedDialogView{
		Dialog:      builder.GetDialogObj("id_dialog_wallet_show_seed"),
		TextView:    builder.GetTextViewObj("id_textview_seed"),
		Image:       builder.GetImageObj("id_image_seed"),
		ButtonClose: builder.GetButtonObj("id_button_close"),
		ViewBuilder: builder,
	}

	view.Image.SetFromPixbuf(assets.ImageSeedPixbuf)
	view.ButtonClose.SetImage(gtkutil.ImageFromPixbuf(assets.IconClosePixbuf16))

	return view
}
