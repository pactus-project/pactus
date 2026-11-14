//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/diamondburned/gotk4/pkg/pango"
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

// DropDownSetup fills the drop-down with the given items. The selected item is
// ellipsized in the middle, so a long address does not set the minimum width of
// the dialog; the popup list still shows each item in full.
func DropDownSetup(drop *gtk.DropDown, names []string) {
	model := gtk.NewStringList(names)
	drop.SetModel(model)
	drop.SetFactory(&dropDownLabelFactory(pango.EllipsizeMiddle).ListItemFactory)
	drop.SetListFactory(&dropDownLabelFactory(pango.EllipsizeNone).ListItemFactory)
}

// dropDownLabelFactory makes a factory that shows each string item in a label.
func dropDownLabelFactory(ellipsize pango.EllipsizeMode) *gtk.SignalListItemFactory {
	factory := gtk.NewSignalListItemFactory()
	factory.ConnectSetup(func(obj *glib.Object) {
		item := obj.Cast().(*gtk.ListItem)
		label := gtk.NewLabel("")
		label.SetXAlign(0)
		label.SetEllipsize(ellipsize)
		item.SetChild(label)
	})
	factory.ConnectBind(func(obj *glib.Object) {
		item := obj.Cast().(*gtk.ListItem)
		str, ok := item.Item().Cast().(*gtk.StringObject)
		if !ok {
			return
		}
		label := item.Child().(*gtk.Label)
		label.SetText(str.String())
		label.SetTooltipText(str.String())
	})

	return factory
}

func DropDownGetSelectedItem[T any](drop *gtk.DropDown, values []T) T {
	// No item is selected when the list is empty: the position is then invalid.
	idx := drop.Selected()
	if idx >= uint(len(values)) {
		var zero T

		return zero
	}

	return values[idx]
}

// DropDownGetSelectedText returns the text of the selected item,
// or an empty string if no item is selected, for example when the list is empty.
func DropDownGetSelectedText(drop *gtk.DropDown) string {
	item := drop.SelectedItem()
	if item == nil {
		return ""
	}

	str, ok := item.Cast().(*gtk.StringObject)
	if !ok {
		return ""
	}

	return str.String()
}

// DropDownGetSelectedAddress returns the address of the selected address list item.
func DropDownGetSelectedAddress(drop *gtk.DropDown) string {
	return AddressFromItem(DropDownGetSelectedText(drop))
}

func DropDownOnChanged(drop *gtk.DropDown, callback func()) {
	drop.NotifyProperty("selected", callback)
}
