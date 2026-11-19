//go:build gtk

package controller

import (
	"runtime"
	"testing"

	"github.com/diamondburned/gotk4/pkg/core/gioutil"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/stretchr/testify/assert"
)

// Without a selected address, the actions on the selected address do nothing.
func TestWalletActionsWithoutSelectedAddress(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if !gtk.InitCheck() {
		t.Skip("no display to run GTK")
	}

	addresses := gioutil.NewListModel[addressRow]()
	colView := gtk.NewColumnView(gtk.NewSingleSelection(addresses))
	ctrl := &WalletWidgetController{
		view:        &view.WalletWidgetView{ColViewAddresses: colView},
		lsAddresses: addresses,
	}

	assert.Nil(t, ctrl.getSelectedAddress())
	assert.NotPanics(t, ctrl.ShowUpdateLabel)
	assert.NotPanics(t, ctrl.ShowAddressDetails)
	assert.NotPanics(t, ctrl.ShowPrivateKey)
}
