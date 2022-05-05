//go:build gtk

package main

import (
	_ "embed"
	"time"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/main_window.ui
var uiMainWindow []byte

type mainWindow struct {
	*gtk.ApplicationWindow

	widgetNode   *widgetNode
	widgetWallet *widgetWallet
}

func buildMainWindow(nodeModel *nodeModel, walletModel *walletModel, genesisTime time.Time) *mainWindow {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromString(string(uiMainWindow))
	errorCheck(nil, err)

	appWindow := getApplicationWindowObj(builder, "id_main_window")
	boxNode := getBoxObj(builder, "id_box_node")
	boxDefaultWallet := getBoxObj(builder, "id_box_default_wallet")

	widgetNode, err := buildWidgetNode(nodeModel, genesisTime)
	errorCheck(appWindow, err)
	widgetWallet, err := buildWidgetWallet(walletModel)
	errorCheck(appWindow, err)

	boxNode.Add(widgetNode)
	boxDefaultWallet.Add(widgetWallet)

	mw := &mainWindow{
		ApplicationWindow: appWindow,
		widgetNode:        widgetNode,
		widgetWallet:      widgetWallet,
	}

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{
		"on_about_gtk": mw.onAboutGtk,
		"on_about":     mw.onAbout,
		"on_quit":      mw.onQuit,
	}
	builder.ConnectSignals(signals)

	return mw
}

func (mw *mainWindow) onQuit() {
	mw.Close()
}

func (mw *mainWindow) onAboutGtk() {
	showAboutGTKDialog(mw)
}

func (mw *mainWindow) onAbout() {
	showAboutDialog(mw)
}
