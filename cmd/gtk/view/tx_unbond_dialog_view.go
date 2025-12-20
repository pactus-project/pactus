//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type TxUnbondDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	ValidatorCombo *gtk.ComboBoxText
	ValidatorHint  *gtk.Label
	MemoEntry      *gtk.Entry

	ButtonCancel *gtk.Button
	ButtonSend   *gtk.Button
}

func NewTxUnbondDialogView() *TxUnbondDialogView {
	builder := NewViewBuilder(assets.TxUnbondDialogUI)

	view := &TxUnbondDialogView{
		ViewBuilder: builder,
		Dialog:      builder.GetDialogObj("id_dialog_transaction_unbond"),

		ValidatorCombo: builder.GetComboBoxTextObj("id_combo_validator"),
		ValidatorHint:  builder.GetLabelObj("id_hint_validator"),
		MemoEntry:      builder.GetEntryObj("id_entry_memo"),

		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ButtonSend:   builder.GetButtonObj("id_button_send"),
	}

	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))
	view.ButtonSend.SetImage(gtkutil.ImageFromPixbuf(assets.IconSendPixbuf16))

	return view
}
