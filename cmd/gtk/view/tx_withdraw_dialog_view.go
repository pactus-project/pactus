//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type TxWithdrawDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	ValidatorCombo *gtk.ComboBoxText
	ValidatorHint  *gtk.Label
	ReceiverCombo  *gtk.ComboBoxText
	ReceiverHint   *gtk.Label
	StakeEntry     *gtk.Entry
	StakeHint      *gtk.Label
	FeeEntry       *gtk.Entry
	FeeHint        *gtk.Label
	MemoEntry      *gtk.Entry

	ButtonCancel *gtk.Button
	ButtonSend   *gtk.Button
}

func NewTxWithdrawDialogView() *TxWithdrawDialogView {
	builder := NewViewBuilder(assets.TxWithdrawDialogUI)

	view := &TxWithdrawDialogView{
		ViewBuilder: builder,
		Dialog:      builder.GetDialogObj("id_dialog_transaction_withdraw"),

		ValidatorCombo: builder.GetComboBoxTextObj("id_combo_validator"),
		ValidatorHint:  builder.GetLabelObj("id_hint_validator"),
		ReceiverCombo:  builder.GetComboBoxTextObj("id_combo_receiver"),
		ReceiverHint:   builder.GetLabelObj("id_hint_receiver"),
		StakeEntry:     builder.GetEntryObj("id_entry_stake"),
		StakeHint:      builder.GetLabelObj("id_hint_stake"),
		FeeEntry:       builder.GetEntryObj("id_entry_fee"),
		FeeHint:        builder.GetLabelObj("id_hint_fee"),
		MemoEntry:      builder.GetEntryObj("id_entry_memo"),

		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ButtonSend:   builder.GetButtonObj("id_button_send"),
	}

	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))
	view.ButtonSend.SetImage(gtkutil.ImageFromPixbuf(assets.IconSendPixbuf16))

	return view
}
