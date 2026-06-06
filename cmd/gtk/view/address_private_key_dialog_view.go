//go:build gtk

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
		ViewBuilder:  builder,
		Window:       builder.GetWindowObj("id_dialog_address_private_key"),
		AddressEntry: builder.BuildExtendedEntry("id_overlay_address"),
		PrvKeyEntry:  builder.BuildExtendedEntry("id_overlay_private_key"),
		ButtonClose:  builder.GetButtonObj("id_button_close"),
	}

	gtkutil.UpdateCloseButton(view.ButtonClose)

	gtkutil.DialogSetup(view.Window)

	return view
}
