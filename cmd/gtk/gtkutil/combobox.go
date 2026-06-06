//go:build gtk

//nolint:staticcheck // Using depreciated widgets
package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func ComboBoxFromAddressList(combo *gtk.ComboBoxText, infos []*pactus.AddressInfo) {
	addrs := make([]string, 0, len(infos))
	for _, info := range infos {
		addrs = append(addrs, info.Address)
	}
	ComboBoxSetup(combo, addrs)
}

func ComboBoxSetup(combo *gtk.ComboBoxText, names []string) {
	for _, name := range names {
		combo.AppendText(name)
	}
}

func ComboBoxGetSelectedText(combo *gtk.ComboBoxText) string {
	return combo.ActiveText()
}

func ComboBoxOnChanged(combo *gtk.ComboBoxText, callback func()) {
	combo.ConnectChanged(callback)
}
