//go:build gtk

package gtkutil

import (
	"runtime"
	"testing"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/stretchr/testify/assert"
)

// An empty drop-down list has no selected item: its text is empty.
func TestDropDownGetSelectedTextOfEmptyList(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if !gtk.InitCheck() {
		t.Skip("no display to run GTK")
	}

	assert.Empty(t, DropDownGetSelectedText(gtk.NewDropDownFromStrings([]string{})))
	assert.Empty(t, DropDownGetSelectedAddress(gtk.NewDropDownFromStrings([]string{})))
	assert.Equal(t, "first", DropDownGetSelectedText(gtk.NewDropDownFromStrings([]string{"first", "second"})))

	// No position is selected in an empty list, and a position past the values is not used.
	assert.Empty(t, DropDownGetSelectedItem(gtk.NewDropDownFromStrings([]string{}), []string{}))
	drop := gtk.NewDropDownFromStrings([]string{"first", "second"})
	drop.SetSelected(1)
	assert.Empty(t, DropDownGetSelectedItem(drop, []string{"only one value"}))
	assert.Equal(t, "second", DropDownGetSelectedItem(drop, []string{"first", "second"}))
}
