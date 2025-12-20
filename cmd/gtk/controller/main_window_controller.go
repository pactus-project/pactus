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
	if h == nil {
		h = &MainWindowHandlers{}
	}

	// Top menu items.
	connectMenuItem(c.view.ExplorerMenuItem, h.OnMenuActivateExplorer)
	connectMenuItem(c.view.WebsiteMenuItem, h.OnMenuActivateWebsite)
	connectMenuItem(c.view.DocumentationMenuItem, h.OnMenuActivateDocs)

	signals := map[string]any{
		"on_about_gtk": func() {
			callIfSet(h.OnAboutGtk)
		},
		"on_about": func() {
			callIfSet(h.OnAbout)
		},
		"on_quit": func() {
			callIfSet(h.OnQuit)
		},
		"on_transaction_transfer": func() {
			callIfSet(h.OnTransactionTransfer)
		},
		"on_transaction_bond": func() {
			callIfSet(h.OnTransactionBond)
		},
		"on_transaction_unbond": func() {
			callIfSet(h.OnTransactionUnbond)
		},
		"on_transaction_withdraw": func() {
			callIfSet(h.OnTransactionWithdraw)
		},
		"on_wallet_new_address": func() {
			callIfSet(h.OnWalletNewAddress)
		},
		"on_wallet_change_password": func() {
			callIfSet(h.OnWalletChangePassword)
		},
		"on_wallet_show_seed": func() {
			callIfSet(h.OnWalletShowSeed)
		},
		"on_wallet_set_default_fee": func() {
			callIfSet(h.OnWalletSetDefaultFee)
		},
	}

	c.view.ConnectSignals(signals)
}

func connectMenuItem(item *gtk.MenuItem, onActivate func()) {
	if item == nil {
		return
	}

	item.Connect("activate", func(_ *gtk.MenuItem) {
		callIfSet(onActivate)
	})
}

func callIfSet(fn func()) {
	if fn == nil {
		return
	}

	fn()
}
