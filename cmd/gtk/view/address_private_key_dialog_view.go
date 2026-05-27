//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type AddressPrivateKeyDialogView struct {
	ViewBuilder

	Window *gtk.Window

	AddressEntry *gtk.Entry
	PrvKeyEntry  *gtk.Entry
	ButtonClose  *gtk.Button
}

func NewAddressPrivateKeyDialogView() *AddressPrivateKeyDialogView {
	builder := NewViewBuilder(assets.AddressPrivateKeyDialogUI)

	view := &AddressPrivateKeyDialogView{
		ViewBuilder: builder,
		Window:      builder.GetWindowObj("id_dialog_address_private_key"),

		AddressEntry: builder.BuildExtendedEntry("id_entry_address"),
		PrvKeyEntry:  builder.BuildExtendedEntry("id_entry_private_key"),
		ButtonClose:  builder.GetButtonObj("id_button_close"),
	}

	gtkutil.AddImageToButton(view.ButtonClose, assets.IconClose16)

	return view
}
