//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletCreateAddressDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	LabelEntry       *gtk.Entry
	AddressTypeCombo *gtk.ComboBoxText
	ButtonOK         *gtk.Button
	ButtonCancel     *gtk.Button
}

func NewWalletCreateAddressDialogView() *WalletCreateAddressDialogView {
	builder := NewViewBuilder(assets.WalletCreateAddressDialogUI)

	view := &WalletCreateAddressDialogView{
		Dialog:           builder.GetDialogObj("id_dialog_wallet_create_address"),
		LabelEntry:       builder.GetEntryObj("id_entry_account_label"),
		AddressTypeCombo: builder.GetComboBoxTextObj("id_combo_address_type"),
		ButtonOK:         builder.GetButtonObj("id_button_ok"),
		ButtonCancel:     builder.GetButtonObj("id_button_cancel"),
		ViewBuilder:      builder,
	}

	view.ButtonOK.SetImage(gtkutil.ImageFromPixbuf(assets.IconOkPixbuf16))
	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))

	return view
}
