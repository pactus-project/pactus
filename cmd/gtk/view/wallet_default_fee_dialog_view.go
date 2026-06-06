//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type WalletDefaultFeeDialogView struct {
	ViewBuilder

	Window *gtk.Window

	FeeEntry        *gtk.Entry
	CurrentFeeLabel *gtk.Label
	ButtonOK        *gtk.Button
	ButtonCancel    *gtk.Button
}

func NewWalletDefaultFeeDialogView() *WalletDefaultFeeDialogView {
	builder := NewViewBuilder(assets.WalletSetDefaultFeeDialogUI)

	view := &WalletDefaultFeeDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_wallet_set_default_fee"),

		FeeEntry:        builder.GetEntryObj("id_entry_default_fee"),
		CurrentFeeLabel: builder.GetLabelObj("id_label_current_fee_value"),
		ButtonOK:        builder.GetButtonObj("id_button_ok"),
		ButtonCancel:    builder.GetButtonObj("id_button_cancel"),
	}

	gtkutil.UpdateOKButton(view.Window, view.ButtonOK)
	gtkutil.UpdateCancelButton(view.ButtonCancel)

	gtkutil.DialogSetup(view.Window)

	return view
}
