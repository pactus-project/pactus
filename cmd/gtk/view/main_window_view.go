//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
)

type MainWindowView struct {
	ViewBuilder

	Window *gtk.ApplicationWindow

	BoxNode          *gtk.Box
	BoxDefaultWallet *gtk.Box
	BoxValidators    *gtk.Box
}

func NewMainWindowView() (*MainWindowView, error) {
	builder := NewViewBuilder(assets.MainWindowUI)

	boxNode := builder.GetBoxObj("id_box_node")
	boxDefaultWallet := builder.GetBoxObj("id_box_default_wallet")
	boxValidators := builder.GetBoxObj("id_box_validators")

	view := &MainWindowView{
		ViewBuilder: builder,
		Window:      builder.GetApplicationWindowObj("id_main_window"),

		BoxNode:          boxNode,
		BoxDefaultWallet: boxDefaultWallet,
		BoxValidators:    boxValidators,
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

func (v *MainWindowView) Cleanup() {
	v.Window.Destroy()
}
