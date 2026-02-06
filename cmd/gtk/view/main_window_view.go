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

	BoxNode          *gtk.Box
	BoxDefaultWallet *gtk.Box
	BoxValidators    *gtk.Box
}

func NewMainWindowView() *MainWindowView {
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
