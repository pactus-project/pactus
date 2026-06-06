//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type AddressLabelDialogView struct {
	ViewBuilder

	Window *gtk.Window

	LabelEntry   *gtk.Entry
	ButtonOK     *gtk.Button
	ButtonCancel *gtk.Button
}

func NewAddressLabelDialogView() *AddressLabelDialogView {
	builder := NewViewBuilder(assets.AddressLabelDialogUI)

	view := &AddressLabelDialogView{
		ViewBuilder:  builder,
		Window:       builder.GetWindowObj("id_dialog_address_label"),
		LabelEntry:   builder.GetEntryObj("id_entry_label"),
		ButtonOK:     builder.GetButtonObj("id_button_ok"),
		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
	}

	gtkutil.UpdateOKButton(view.Window, view.ButtonOK)
	gtkutil.UpdateCancelButton(view.ButtonCancel)

	gtkutil.DialogSetup(view.Window)

	return view
}
