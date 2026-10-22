//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/types/amount"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func DropDownFromAddressList(drop *gtk.DropDown, infos []*pactus.AddressInfo) {
	addrs := make([]string, 0, len(infos))
	for _, info := range infos {
		addrs = append(addrs, info.Address)
	}
	DropDownSetup(drop, addrs)
}

// DropDownFromAddressListWithAmount fills the drop-down with "address - amount"
// items. Use DropDownGetSelectedAddress to read the selected address.
func DropDownFromAddressListWithAmount(drop *gtk.DropDown, infos []*pactus.AddressInfo,
	getAmount func(addr string) (amount.Amount, error),
) {
	DropDownSetup(drop, AddressItems(infos, getAmount))
}

func DropDownSetup(drop *gtk.DropDown, names []string) {
	model := gtk.NewStringList(names)
	drop.SetModel(model)
}

func DropDownGetSelectedItem[T any](drop *gtk.DropDown, values []T) T {
	idx := drop.Selected()
	if idx > uint(len(values)) {
		var zero T

		return zero
	}

	return values[idx]
}

func DropDownGetSelectedText(drop *gtk.DropDown) string {
	item := drop.SelectedItem()
	str := item.Cast().(*gtk.StringObject)

	return str.String()
}

// DropDownGetSelectedAddress returns the address of the selected address list item.
func DropDownGetSelectedAddress(drop *gtk.DropDown) string {
	return AddressFromItem(DropDownGetSelectedText(drop))
}

func DropDownOnChanged(drop *gtk.DropDown, callback func()) {
	drop.NotifyProperty("selected", callback)
}
