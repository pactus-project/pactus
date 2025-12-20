//go:build gtk

package controller

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type MainWindowHandlers struct {
	OnAboutGtk             func()
	OnAbout                func()
	OnQuit                 func()
	OnTransactionTransfer  func()
	OnTransactionBond      func()
	OnTransactionUnbond    func()
	OnTransactionWithdraw  func()
	OnWalletNewAddress     func()
	OnWalletChangePassword func()
	OnWalletShowSeed       func()
	OnWalletSetDefaultFee  func()
	OnMenuActivateExplorer func()
	OnMenuActivateWebsite  func()
	OnMenuActivateDocs     func()
}

type MainWindowController struct {
	view *view.MainWindowView
}

func NewMainWindowController(view *view.MainWindowView) *MainWindowController {
	return &MainWindowController{view: view}
}

func (c *MainWindowController) Bind(h *MainWindowHandlers) {
	// Top menu items.
	connectMenuItem(c.view.ExplorerMenuItem, h.OnMenuActivateExplorer)
	connectMenuItem(c.view.WebsiteMenuItem, h.OnMenuActivateWebsite)
	connectMenuItem(c.view.DocumentationMenuItem, h.OnMenuActivateDocs)

	signals := map[string]any{
		"on_about_gtk":              h.OnAboutGtk,
		"on_about":                  h.OnAbout,
		"on_quit":                   h.OnQuit,
		"on_transaction_transfer":   h.OnTransactionTransfer,
		"on_transaction_bond":       h.OnTransactionBond,
		"on_transaction_unbond":     h.OnTransactionUnbond,
		"on_transaction_withdraw":   h.OnTransactionWithdraw,
		"on_wallet_new_address":     h.OnWalletNewAddress,
		"on_wallet_change_password": h.OnWalletChangePassword,
		"on_wallet_show_seed":       h.OnWalletShowSeed,
		"on_wallet_set_default_fee": h.OnWalletSetDefaultFee,
	}

	c.view.ConnectSignals(signals)
}

func connectMenuItem(item *gtk.MenuItem, onActivate func()) {
	item.Connect("activate", func(_ *gtk.MenuItem) {
		onActivate()
	})
}
