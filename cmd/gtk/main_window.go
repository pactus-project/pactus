//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

//go:embed assets/ui/main_window.ui
var uiMainWindow []byte

//go:embed assets/css/style.css
var cssCustom string

type mainWindow struct {
	*gtk.ApplicationWindow

	widgetNode   *widgetNode
	widgetWallet *widgetWallet
}

func buildMainWindow(nodeModel *nodeModel, walletModel *walletModel) *mainWindow {
	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromString(string(uiMainWindow))
	fatalErrorCheck(err)

	appWindow := getApplicationWindowObj(builder, "id_main_window")
	boxNode := getBoxObj(builder, "id_box_node")
	boxDefaultWallet := getBoxObj(builder, "id_box_default_wallet")

	widgetNode, err := buildWidgetNode(nodeModel)
	fatalErrorCheck(err)

	widgetWallet, err := buildWidgetWallet(walletModel)
	fatalErrorCheck(err)

	boxNode.Add(widgetNode)
	boxDefaultWallet.Add(widgetWallet)

	mainWnd := &mainWindow{
		ApplicationWindow: appWindow,
		widgetNode:        widgetNode,
		widgetWallet:      widgetWallet,
	}

	explorerItemMenu := getMenuItem(builder, "id_explorer_menu")
	explorerItemMenu.Connect("activate", mainWnd.onMenuItemActivateExplorer)

	websiteItemMenu := getMenuItem(builder, "id_website_menu")
	websiteItemMenu.Connect("activate", mainWnd.onMenuItemActivateWebsite)

	documentationItemMenu := getMenuItem(builder, "id_documentation_menu")
	documentationItemMenu.Connect("activate", mainWnd.onMenuItemActivateDocumentation)

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]any{
		"on_about_gtk":            mainWnd.onAboutGtk,
		"on_about":                mainWnd.onAbout,
		"on_quit":                 mainWnd.onQuit,
		"on_transaction_transfer": mainWnd.OnTransactionTransfer,
		"on_transaction_bond":     mainWnd.OnTransactionBond,
		"on_transaction_unbond":   mainWnd.OnTransactionUnbond,
		"on_transaction_withdraw": mainWnd.OnTransactionWithdraw,
	}
	builder.ConnectSignals(signals)

	// apply custom css
	provider, err := gtk.CssProviderNew()
	fatalErrorCheck(err)

	err = provider.LoadFromData(cssCustom)
	fatalErrorCheck(err)

	screen, err := gdk.ScreenGetDefault()
	fatalErrorCheck(err)

	gtk.AddProviderForScreen(screen, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	return mainWnd
}

func (mw *mainWindow) onQuit() {
	mw.Close()
}

func (*mainWindow) onAboutGtk() {
	showAboutGTKDialog()
}

func (*mainWindow) onAbout() {
	dlg := aboutDialog()

	runDialog(&dlg.Dialog)
}

func (mw *mainWindow) OnTransactionTransfer() {
	broadcastTransactionTransfer(mw.widgetWallet.model.wallet)
}

func (mw *mainWindow) OnTransactionBond() {
	broadcastTransactionBond(mw.widgetWallet.model.wallet)
}

func (mw *mainWindow) OnTransactionUnbond() {
	broadcastTransactionUnbond(mw.widgetWallet.model.wallet)
}

func (mw *mainWindow) OnTransactionWithdraw() {
	broadcastTransactionWithdraw(mw.widgetWallet.model.wallet)
}

func (*mainWindow) onMenuItemActivateWebsite(_ *gtk.MenuItem) {
	if err := openURLInBrowser("https://pactus.org/"); err != nil {
		fatalErrorCheck(err)
	}
}

func (*mainWindow) onMenuItemActivateExplorer(_ *gtk.MenuItem) {
	if err := openURLInBrowser("https://pacviewer.com/"); err != nil {
		fatalErrorCheck(err)
	}
}

func (*mainWindow) onMenuItemActivateDocumentation(_ *gtk.MenuItem) {
	if err := openURLInBrowser("https://docs.pactus.org/"); err != nil {
		fatalErrorCheck(err)
	}
}
