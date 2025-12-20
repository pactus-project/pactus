//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type AddressDetailsDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	AddressEntry *gtk.Entry
	PubKeyEntry  *gtk.Entry
	PathEntry    *gtk.Entry

	ButtonClose *gtk.Button
}

func NewAddressDetailsDialogView() *AddressDetailsDialogView {
	builder := NewViewBuilder(assets.AddressDetailsDialogUI)

	view := &AddressDetailsDialogView{
		Dialog: builder.GetDialogObj("id_dialog_address_details"),

		AddressEntry: builder.BuildExtendedEntry("id_overlay_address"),
		PubKeyEntry:  builder.BuildExtendedEntry("id_overlay_public_key"),
		PathEntry:    builder.GetEntryObj("id_entry_path"),

		ButtonClose: builder.GetButtonObj("id_button_close"),

		ViewBuilder: builder,
	}

	view.ButtonClose.SetImage(gtkutil.ImageFromPixbuf(assets.IconClosePixbuf16))

	return view
}
