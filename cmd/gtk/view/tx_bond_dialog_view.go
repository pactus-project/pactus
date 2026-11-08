//go:build gtk

//nolint:staticcheck // Using depreciated widgets
package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type TxBondDialogView struct {
	ViewBuilder

	Window *gtk.Window

	SenderDrop      *gtk.DropDown
	SenderHint      *gtk.Label
	ReceiverCombo   *gtk.ComboBoxText
	ReceiverHint    *gtk.Label
	PublicKeyEntry  *gtk.Entry
	AmountEntry     *gtk.Entry
	AmountHint      *gtk.Label
	AmountMaxButton *gtk.Button
	FeeEntry        *gtk.Entry
	FeeHint         *gtk.Label
	MemoEntry       *gtk.Entry

	// Delegation (PIP-49)
	DelegateCheck *gtk.CheckButton
	DelegateBox   *gtk.Box
	OwnerLabel    *gtk.Label
	ShareEntry    *gtk.Entry
	ShareHint     *gtk.Label
	TermDrop      *gtk.DropDown
	ExpiryEntry   *gtk.Entry
	ExpiryHint    *gtk.Label

	ButtonCancel *gtk.Button
	ButtonSend   *gtk.Button
}

func NewTxBondDialogView() *TxBondDialogView {
	builder := NewViewBuilder(assets.TxBondDialogUI)

	view := &TxBondDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_transaction_bond"),

		SenderDrop:      builder.GetDropDownObj("id_drop_sender"),
		SenderHint:      builder.GetLabelObj("id_hint_sender"),
		ReceiverCombo:   builder.GetComboBoxTextObj("id_combo_receiver"),
		ReceiverHint:    builder.GetLabelObj("id_hint_receiver"),
		PublicKeyEntry:  builder.GetEntryObj("id_entry_public_key"),
		AmountEntry:     builder.GetEntryObj("id_entry_amount"),
		AmountHint:      builder.GetLabelObj("id_hint_amount"),
		AmountMaxButton: builder.GetButtonObj("id_button_amount_max"),
		FeeEntry:        builder.GetEntryObj("id_entry_fee"),
		FeeHint:         builder.GetLabelObj("id_hint_fee"),
		MemoEntry:       builder.GetEntryObj("id_entry_memo"),

		DelegateCheck: GetObj[*gtk.CheckButton](builder.Builder(), "id_check_delegate"),
		DelegateBox:   builder.GetBoxObj("id_box_delegate"),
		OwnerLabel:    builder.GetLabelObj("id_label_owner"),
		ShareEntry:    builder.GetEntryObj("id_entry_share"),
		ShareHint:     builder.GetLabelObj("id_hint_share"),
		TermDrop:      builder.GetDropDownObj("id_drop_term"),
		ExpiryEntry:   builder.GetEntryObj("id_entry_expiry"),
		ExpiryHint:    builder.GetLabelObj("id_hint_expiry"),

		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ButtonSend:   builder.GetButtonObj("id_button_send"),
	}

	gtkutil.UpdateCancelButton(view.ButtonCancel)
	gtkutil.UpdateSendButton(view.ButtonSend)

	gtkutil.DialogSetup(view.Window)

	return view
}
