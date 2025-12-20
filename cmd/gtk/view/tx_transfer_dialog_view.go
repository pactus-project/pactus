//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type TxTransferDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	SenderCombo   *gtk.ComboBoxText
	SenderHint    *gtk.Label
	ReceiverEntry *gtk.Entry
	ReceiverHint  *gtk.Label
	AmountEntry   *gtk.Entry
	AmountHint    *gtk.Label
	FeeEntry      *gtk.Entry
	FeeHint       *gtk.Label
	MemoEntry     *gtk.Entry

	ButtonCancel *gtk.Button
	ButtonSend   *gtk.Button
}

func NewTxTransferDialogView() *TxTransferDialogView {
	builder := NewViewBuilder(assets.TxTransferDialogUI)

	view := &TxTransferDialogView{
		Dialog: builder.GetDialogObj("id_dialog_transaction_transfer"),

		SenderCombo:   builder.GetComboBoxTextObj("id_combo_sender"),
		SenderHint:    builder.GetLabelObj("id_hint_sender"),
		ReceiverEntry: builder.GetEntryObj("id_entry_receiver"),
		ReceiverHint:  builder.GetLabelObj("id_hint_receiver"),
		AmountEntry:   builder.GetEntryObj("id_entry_amount"),
		AmountHint:    builder.GetLabelObj("id_hint_amount"),
		FeeEntry:      builder.GetEntryObj("id_entry_fee"),
		FeeHint:       builder.GetLabelObj("id_hint_fee"),
		MemoEntry:     builder.GetEntryObj("id_entry_memo"),

		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ButtonSend:   builder.GetButtonObj("id_button_send"),

		ViewBuilder: builder,
	}

	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))
	view.ButtonSend.SetImage(gtkutil.ImageFromPixbuf(assets.IconSendPixbuf16))

	return view
}
