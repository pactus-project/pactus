//go:build gtk

//nolint:staticcheck // Using depreciated widgets
package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/types/amount"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func ComboBoxFromAddressList(combo *gtk.ComboBoxText, infos []*pactus.AddressInfo) {
	addrs := make([]string, 0, len(infos))
	for _, info := range infos {
		addrs = append(addrs, info.Address)
	}
	ComboBoxSetup(combo, addrs)
}

// ComboBoxFromAddressListWithAmount fills the combo box with "address - amount"
// items. Use ComboBoxGetSelectedAddress to read the selected or typed address.
func ComboBoxFromAddressListWithAmount(combo *gtk.ComboBoxText, infos []*pactus.AddressInfo,
	getAmount func(addr string) (amount.Amount, error),
) {
	ComboBoxSetup(combo, AddressItems(infos, getAmount))
}

func ComboBoxSetup(combo *gtk.ComboBoxText, names []string) {
	for _, name := range names {
		combo.AppendText(name)
	}
}

func ComboBoxGetSelectedText(combo *gtk.ComboBoxText) string {
	return combo.ActiveText()
}

// ComboBoxGetSelectedAddress returns the address of the selected item, or the
// address typed in the combo box entry.
func ComboBoxGetSelectedAddress(combo *gtk.ComboBoxText) string {
	return AddressFromItem(combo.ActiveText())
}

func ComboBoxOnChanged(combo *gtk.ComboBoxText, callback func()) {
	combo.ConnectChanged(callback)
}
