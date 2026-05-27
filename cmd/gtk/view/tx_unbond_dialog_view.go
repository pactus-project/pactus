//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type TxUnbondDialogView struct {
	ViewBuilder

	Window *gtk.Window

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
		Window:      builder.GetWindowObj("id_dialog_transaction_unbond"),

		ValidatorCombo: builder.GetComboBoxTextObj("id_combo_validator"),
		ValidatorHint:  builder.GetLabelObj("id_hint_validator"),
		MemoEntry:      builder.GetEntryObj("id_entry_memo"),

		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ButtonSend:   builder.GetButtonObj("id_button_send"),
	}

	gtkutil.AddImageToButton(view.ButtonCancel, assets.IconCancel16)
	gtkutil.AddImageToButton(view.ButtonSend, assets.IconSend16)

	return view
}
