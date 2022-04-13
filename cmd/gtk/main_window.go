package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

//go:embed ui/main_window.ui
var uiMainWindow []byte

type mainWindow struct {
	*gtk.ApplicationWindow

	overviewWidget    *gtk.Widget
	addressesTreeView *addressesTreeView
}

func buildMainWindow() *mainWindow {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromString(string(uiMainWindow))
	errorCheck(err)

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{}
	builder.ConnectSignals(signals)

	objMainWindow, err := builder.GetObject("main_window")
	errorCheck(err)

	appWindow, err := isApplicationWindow(objMainWindow)
	errorCheck(err)

	return &mainWindow{
		ApplicationWindow: appWindow,
		addressesTreeView: buildAddressesTreeView(builder),
	}
}

func (mw *mainWindow) SetAddressesModel(model gtk.ITreeModel) {
	mw.addressesTreeView.SetModel(model)
}
