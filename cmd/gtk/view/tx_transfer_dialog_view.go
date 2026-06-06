//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type TxTransferDialogView struct {
	ViewBuilder

	Window *gtk.Window

	SenderDrop    *gtk.DropDown
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
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_transaction_transfer"),

		SenderDrop:    builder.GetDropDownObj("id_drop_sender"),
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
	}

	gtkutil.UpdateCancelButton(view.ButtonCancel)
	gtkutil.UpdateSendButton(view.ButtonSend)

	gtkutil.DialogSetup(view.Window)

	return view
}
