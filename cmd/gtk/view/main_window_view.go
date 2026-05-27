//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type MainWindowView struct {
	ViewBuilder

	Window *gtk.ApplicationWindow

	// MenuEditConfig   *gtk.Button
	BoxNode          *gtk.Box
	BoxDefaultWallet *gtk.Box
	BoxValidators    *gtk.Box
	BoxCommittee     *gtk.Box
	BoxNetwork       *gtk.Box
}

func NewMainWindowView() *MainWindowView {
	builder := NewViewBuilder(assets.MainWindowUI)

	view := &MainWindowView{
		ViewBuilder: builder,
		Window:      builder.GetApplicationWindowObj("id_main_window"),

		// MenuEditConfig:   builder.GetButtonObj("id_menu_edit_config"),
		BoxNode:          builder.GetBoxObj("id_box_node"),
		BoxDefaultWallet: builder.GetBoxObj("id_box_default_wallet"),
		BoxValidators:    builder.GetBoxObj("id_box_validators"),
		BoxCommittee:     builder.GetBoxObj("id_box_committee"),
		BoxNetwork:       builder.GetBoxObj("id_box_network"),
	}

	// apply custom css
	provider := gtk.NewCSSProvider()
	provider.LoadFromData(assets.MainWindowCSS)
	display := gdk.DisplayGetDefault()

	gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return view
}

func (v *MainWindowView) Cleanup() {
	v.Window.Destroy()
}
