//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type MainWindowView struct {
	ViewBuilder

	Window *gtk.ApplicationWindow

	// MenuEditConfig   *gtk.Button
	BoxNode       *gtk.Box
	BoxWallet     *gtk.Box
	BoxValidators *gtk.Box
	BoxCommittee  *gtk.Box
	BoxNetwork    *gtk.Box

	// HideOnClose controls the window close behavior.
	// When true, the window is hidden instead of destroyed on close-request.
	HideOnClose bool
}

func NewMainWindowView() *MainWindowView {
	builder := NewViewBuilder(assets.MainWindowUI)

	view := &MainWindowView{
		ViewBuilder:   builder,
		Window:        builder.GetApplicationWindowObj("id_main_window"),
		BoxNode:       builder.GetBoxObj("id_box_node"),
		BoxWallet:     builder.GetBoxObj("id_box_wallet"),
		BoxValidators: builder.GetBoxObj("id_box_validators"),
		BoxCommittee:  builder.GetBoxObj("id_box_committee"),
		BoxNetwork:    builder.GetBoxObj("id_box_network"),
		HideOnClose:   true,
	}

	// Intercept the close-request signal to hide instead of destroy.
	view.Window.ConnectCloseRequest(func() (ok bool) {
		if view.HideOnClose {
			view.Window.SetVisible(false)

			return true // prevent default close/destroy
		}

		return false // allow close
	})

	return view
}

func (v *MainWindowView) Cleanup() {
	v.Window.Destroy()
}
