//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type AddressPrivateKeyDialogView struct {
	ViewBuilder

	Dialog *gtk.Dialog

	AddressEntry *gtk.Entry
	PrvKeyEntry  *gtk.Entry
	ButtonClose  *gtk.Button
}

func NewAddressPrivateKeyDialogView() *AddressPrivateKeyDialogView {
	builder := NewViewBuilder(assets.AddressPrivateKeyDialogUI)

	view := &AddressPrivateKeyDialogView{
		Dialog:       builder.GetDialogObj("id_dialog_address_private_key"),
		AddressEntry: builder.GetEntryObj("id_entry_address"),
		PrvKeyEntry:  builder.GetEntryObj("id_entry_private_key"),
		ButtonClose:  builder.GetButtonObj("id_button_close"),
		ViewBuilder:  builder,
	}

	view.ButtonClose.SetImage(gtkutil.ImageFromPixbuf(assets.IconClosePixbuf16))

	return view
}
