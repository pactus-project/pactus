//go:build gtk

package gtkutil

import "github.com/diamondburned/gotk4/pkg/gtk/v4"

func EntryOnChanged(entry *gtk.Entry, callback func()) {
	entry.ConnectChanged(callback)
}

func EntryGetText(entry *gtk.Entry) string {
	return entry.Text()
}
