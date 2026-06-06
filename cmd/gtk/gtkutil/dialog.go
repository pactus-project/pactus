//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

func DialogSetup(dlg *gtk.Window) {
	keyController := gtk.NewEventControllerKey()
	keyController.ConnectKeyPressed(func(keyval uint, _ uint, _ gdk.ModifierType) bool {
		if keyval == gdk.KEY_Escape {
			dlg.Close()

			return true
		}

		return false
	})
	dlg.AddController(keyController)
}
