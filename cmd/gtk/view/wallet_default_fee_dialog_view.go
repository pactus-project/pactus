//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletDefaultFeeDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	FeeEntry        *gtk.Entry
	CurrentFeeLabel *gtk.Label
	ButtonOK        *gtk.Button
	ButtonCancel    *gtk.Button
}

func NewWalletDefaultFeeDialogView() *WalletDefaultFeeDialogView {
	builder := NewViewBuilder(assets.WalletSetDefaultFeeDialogUI)

	view := &WalletDefaultFeeDialogView{
		ViewBuilder: builder,
		Dialog:      builder.GetDialogObj("id_dialog_wallet_set_default_fee"),

		FeeEntry:        builder.GetEntryObj("id_entry_default_fee"),
		CurrentFeeLabel: builder.GetLabelObj("id_label_current_fee_value"),
		ButtonOK:        builder.GetButtonObj("id_button_ok"),
		ButtonCancel:    builder.GetButtonObj("id_button_cancel"),
	}

	view.ButtonOK.SetImage(gtkutil.ImageFromPixbuf(assets.IconOkPixbuf16))
	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))

	return view
}
