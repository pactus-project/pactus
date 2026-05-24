//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type MainWindowView struct {
	ViewBuilder

	Window *gtk.ApplicationWindow

	MenuEditConfig   *gtk.MenuItem
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

		MenuEditConfig:   builder.GetMenuItem("id_menu_edit_config"),
		BoxNode:          builder.GetBoxObj("id_box_node"),
		BoxDefaultWallet: builder.GetBoxObj("id_box_default_wallet"),
		BoxValidators:    builder.GetBoxObj("id_box_validators"),
		BoxCommittee:     builder.GetBoxObj("id_box_committee"),
		BoxNetwork:       builder.GetBoxObj("id_box_network"),
	}

	// apply custom css
	provider, err := gtk.CssProviderNew()
	gtkutil.FatalErrorCheck(err)

	err = provider.LoadFromData(assets.MainWindowCSS)
	gtkutil.FatalErrorCheck(err)

	screen, err := gdk.ScreenGetDefault()
	gtkutil.FatalErrorCheck(err)

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return view
}

func (v *MainWindowView) Cleanup() {
	v.Window.Destroy()
}
