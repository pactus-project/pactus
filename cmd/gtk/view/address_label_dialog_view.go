//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type AddressLabelDialogView struct {
	ViewBuilder

	Dialog       *gtk.Dialog
	LabelEntry   *gtk.Entry
	ButtonOK     *gtk.Button
	ButtonCancel *gtk.Button
}

func NewAddressLabelDialogView() *AddressLabelDialogView {
	builder := NewViewBuilder(assets.AddressLabelDialogUI)

	view := &AddressLabelDialogView{
		Dialog:       builder.GetDialogObj("id_dialog_address_label"),
		LabelEntry:   builder.GetEntryObj("id_entry_label"),
		ButtonOK:     builder.GetButtonObj("id_button_ok"),
		ButtonCancel: builder.GetButtonObj("id_button_cancel"),
		ViewBuilder:  builder,
	}

	view.ButtonOK.SetImage(gtkutil.ImageFromPixbuf(assets.IconOkPixbuf16))
	view.ButtonCancel.SetImage(gtkutil.ImageFromPixbuf(assets.IconCancelPixbuf16))

	return view
}
