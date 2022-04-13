package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/main_window.ui
var uiMainWindow []byte

type mainWindow struct {
	*gtk.ApplicationWindow

	widgetNode   *widgetNode
	widgetWallet *widgetWallet
}

func buildMainWindow() *mainWindow {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromString(string(uiMainWindow))
	errorCheck(err)

	objMainWindow, err := builder.GetObject("id_main_window")
	errorCheck(err)

	appWindow, err := isApplicationWindow(objMainWindow)
	errorCheck(err)

	objBoxNode, err := builder.GetObject("id_box_node")
	errorCheck(err)

	boxNode, err := isBox(objBoxNode)
	errorCheck(err)

	objBoxDefaultWallet, err := builder.GetObject("id_box_default_wallet")
	errorCheck(err)

	boxDefaultWallet, err := isBox(objBoxDefaultWallet)
	errorCheck(err)

	widgetNode := buildWidgetNode()
	widgetWallet := buildWidgetWallet()

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

func (mw *mainWindow) SetWalletModel(model *walletModel) {
	mw.widgetWallet.SetModel(model)
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
