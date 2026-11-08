//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

// TxAnchorDialogView is the dialog that creates, updates or deletes the
// anchor of an account (PIP-50).
type TxAnchorDialogView struct {
	ViewBuilder

	Window *gtk.Window

	LinkEntry       *gtk.Entry
	LinkApplyButton *gtk.Button
	LinkHint        *gtk.Label

	SenderDrop *gtk.DropDown
	SenderHint *gtk.Label

	NoAnchorLabel *gtk.Label
	CurrentGrid   *gtk.Grid
	CurrentType   *gtk.Label
	CurrentRoot   *gtk.Label
	CurrentURI    *gtk.Label
	CurrentLocked *gtk.Label
	CurrentUpdate *gtk.Label

	SetToggle    *gtk.ToggleButton
	DeleteToggle *gtk.ToggleButton

	SetBox       *gtk.Box
	TypeDrop     *gtk.DropDown
	TypeSpin     *gtk.SpinButton
	TypeHint     *gtk.Label
	RootEntry    *gtk.Entry
	RootHint     *gtk.Label
	URIEntry     *gtk.Entry
	URIHint      *gtk.Label
	DepositEntry *gtk.Entry
	DepositHint  *gtk.Label

	DeleteInfo   *gtk.Label
	WarningLabel *gtk.Label

	FeeEntry  *gtk.Entry
	FeeHint   *gtk.Label
	MemoEntry *gtk.Entry

	ButtonCancel *gtk.Button
	ButtonSend   *gtk.Button
}

func NewTxAnchorDialogView() *TxAnchorDialogView {
	builder := NewViewBuilder(assets.TxAnchorDialogUI)
	gtkBuilder := builder.Builder()

	view := &TxAnchorDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_transaction_anchor"),

		LinkEntry:       builder.GetEntryObj("id_entry_link"),
		LinkApplyButton: builder.GetButtonObj("id_button_link_apply"),
		LinkHint:        builder.GetLabelObj("id_hint_link"),

		SenderDrop: builder.GetDropDownObj("id_drop_sender"),
		SenderHint: builder.GetLabelObj("id_hint_sender"),

		NoAnchorLabel: builder.GetLabelObj("id_label_no_anchor"),
		CurrentGrid:   GetObj[*gtk.Grid](gtkBuilder, "id_grid_current"),
		CurrentType:   builder.GetLabelObj("id_val_type"),
		CurrentRoot:   builder.GetLabelObj("id_val_root"),
		CurrentURI:    builder.GetLabelObj("id_val_uri"),
		CurrentLocked: builder.GetLabelObj("id_val_locked"),
		CurrentUpdate: builder.GetLabelObj("id_val_updated"),

		SetToggle:    GetObj[*gtk.ToggleButton](gtkBuilder, "id_toggle_set"),
		DeleteToggle: GetObj[*gtk.ToggleButton](gtkBuilder, "id_toggle_delete"),

		SetBox:       builder.GetBoxObj("id_box_set"),
		TypeDrop:     builder.GetDropDownObj("id_drop_type"),
		TypeSpin:     GetObj[*gtk.SpinButton](gtkBuilder, "id_spin_type"),
		TypeHint:     builder.GetLabelObj("id_hint_type"),
		RootEntry:    builder.GetEntryObj("id_entry_root"),
		RootHint:     builder.GetLabelObj("id_hint_root"),
		URIEntry:     builder.GetEntryObj("id_entry_uri"),
		URIHint:      builder.GetLabelObj("id_hint_uri"),
		DepositEntry: builder.GetEntryObj("id_entry_deposit"),
		DepositHint:  builder.GetLabelObj("id_hint_deposit"),

		DeleteInfo:   builder.GetLabelObj("id_label_delete_info"),
		WarningLabel: builder.GetLabelObj("id_label_warning"),

		FeeEntry:  builder.GetEntryObj("id_entry_fee"),
		FeeHint:   builder.GetLabelObj("id_hint_fee"),
		MemoEntry: builder.GetEntryObj("id_entry_memo"),

		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ButtonSend:   builder.GetButtonObj("id_button_send"),
	}

	gtkutil.UpdateCancelButton(view.ButtonCancel)
	gtkutil.UpdateSendButton(view.ButtonSend)

	gtkutil.DialogSetup(view.Window)

	return view
}
