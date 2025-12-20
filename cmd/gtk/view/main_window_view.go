//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type MainWindowView struct {
	ViewBuilder

	Window           *gtk.ApplicationWindow
	BoxNode          *gtk.Box
	BoxDefaultWallet *gtk.Box

	ExplorerMenuItem      *gtk.MenuItem
	WebsiteMenuItem       *gtk.MenuItem
	DocumentationMenuItem *gtk.MenuItem
}

func NewMainWindowView() (*MainWindowView, error) {
	builder := NewViewBuilder(assets.MainWindowUI)

	appWindow := builder.GetApplicationWindowObj("id_main_window")
	boxNode := builder.GetBoxObj("id_box_node")
	boxDefaultWallet := builder.GetBoxObj("id_box_default_wallet")

	view := &MainWindowView{
		Window:           appWindow,
		BoxNode:          boxNode,
		BoxDefaultWallet: boxDefaultWallet,

		ExplorerMenuItem:      builder.GetMenuItem("id_explorer_menu"),
		WebsiteMenuItem:       builder.GetMenuItem("id_website_menu"),
		DocumentationMenuItem: builder.GetMenuItem("id_documentation_menu"),

		ViewBuilder: builder,
	}

	// apply custom css
	provider, err := gtk.CssProviderNew()
	if err != nil {
		return nil, err
	}
	if err := provider.LoadFromData(assets.MainWindowCSS); err != nil {
		return nil, err
	}

	screen, err := gdk.ScreenGetDefault()
	if err != nil {
		return nil, err
	}
	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return view, nil
}
