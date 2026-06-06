//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func DropDownFromAddressList(drop *gtk.DropDown, infos []*pactus.AddressInfo) {
	addrs := make([]string, 0, len(infos))
	for _, info := range infos {
		addrs = append(addrs, info.Address)
	}
	DropDownSetup(drop, addrs)
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

func DropDownOnChanged(drop *gtk.DropDown, callback func()) {
	drop.NotifyProperty("selected", callback)
}
