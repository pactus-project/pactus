//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletCreateAddressDialogView struct {
	ViewBuilder

	Window *gtk.Window

	LabelEntry       *gtk.Entry
	AddressTypeCombo *gtk.ComboBoxText
	ButtonOK         *gtk.Button
	ButtonCancel     *gtk.Button
}

func NewWalletCreateAddressDialogView() *WalletCreateAddressDialogView {
	builder := NewViewBuilder(assets.WalletCreateAddressDialogUI)

	view := &WalletCreateAddressDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_wallet_create_address"),

		LabelEntry:       builder.GetEntryObj("id_entry_account_label"),
		AddressTypeCombo: builder.GetComboBoxTextObj("id_combo_address_type"),
		ButtonOK:         builder.GetButtonObj("id_button_ok"),
		ButtonCancel:     builder.GetButtonObj("id_button_cancel"),
	}

	gtkutil.AddImageToButton(view.ButtonOK, assets.IconOk16)
	gtkutil.AddImageToButton(view.ButtonCancel, assets.IconCancel16)

	return view
}
