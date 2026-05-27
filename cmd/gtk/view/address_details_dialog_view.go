//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type AddressDetailsDialogView struct {
	ViewBuilder

	Window *gtk.Window

	AddressEntry *gtk.Entry
	PubKeyEntry  *gtk.Entry
	PathEntry    *gtk.Entry

	ButtonClose *gtk.Button
}

func NewAddressDetailsDialogView() *AddressDetailsDialogView {
	builder := NewViewBuilder(assets.AddressDetailsDialogUI)

	view := &AddressDetailsDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_address_details"),

		AddressEntry: builder.BuildExtendedEntry("id_overlay_address"),
		PubKeyEntry:  builder.BuildExtendedEntry("id_overlay_public_key"),
		PathEntry:    builder.GetEntryObj("id_entry_path"),

		ButtonClose: builder.GetButtonObj("id_button_close"),
	}

	gtkutil.AddImageToButton(view.ButtonClose, assets.IconClose16)

	return view
}
