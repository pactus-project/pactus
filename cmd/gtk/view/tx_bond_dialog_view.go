//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type TxBondDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	SenderCombo    *gtk.ComboBoxText
	SenderHint     *gtk.Label
	ReceiverCombo  *gtk.ComboBoxText
	ReceiverHint   *gtk.Label
	PublicKeyEntry *gtk.Entry
	AmountEntry    *gtk.Entry
	AmountHint     *gtk.Label
	FeeEntry       *gtk.Entry
	FeeHint        *gtk.Label
	MemoEntry      *gtk.Entry

	ButtonCancel *gtk.Button
	ButtonSend   *gtk.Button
}

func NewTxBondDialogView() *TxBondDialogView {
	builder := NewViewBuilder(assets.TxBondDialogUI)

	view := &TxBondDialogView{
		Dialog: builder.GetDialogObj("id_dialog_transaction_bond"),

		SenderCombo:    builder.GetComboBoxTextObj("id_combo_sender"),
		SenderHint:     builder.GetLabelObj("id_hint_sender"),
		ReceiverCombo:  builder.GetComboBoxTextObj("id_combo_receiver"),
		ReceiverHint:   builder.GetLabelObj("id_hint_receiver"),
		PublicKeyEntry: builder.GetEntryObj("id_entry_public_key"),
		AmountEntry:    builder.GetEntryObj("id_entry_amount"),
		AmountHint:     builder.GetLabelObj("id_hint_amount"),
		FeeEntry:       builder.GetEntryObj("id_entry_fee"),
		FeeHint:        builder.GetLabelObj("id_hint_fee"),
		MemoEntry:      builder.GetEntryObj("id_entry_memo"),

		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ButtonSend:   builder.GetButtonObj("id_button_send"),

		ViewBuilder: builder,
	}

	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))
	view.ButtonSend.SetImage(gtkutil.ImageFromPixbuf(assets.IconSendPixbuf16))

	return view
}
