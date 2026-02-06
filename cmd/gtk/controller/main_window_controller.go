//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type MainWindowController struct {
	view *view.MainWindowView
}

func NewMainWindowController(view *view.MainWindowView) *MainWindowController {
	return &MainWindowController{view: view}
}

func (c *MainWindowController) BuildView(nav *Navigator) {
	c.view.ConnectSignals(map[string]any{
		"on_quit":                   nav.Quit,
		"on_transaction_transfer":   nav.ShowTransactionTransfer,
		"on_transaction_bond":       nav.ShowTransactionBond,
		"on_transaction_unbond":     nav.ShowTransactionUnbond,
		"on_transaction_withdraw":   nav.ShowTransactionWithdraw,
		"on_wallet_new_address":     nav.ShowWalletNewAddress,
		"on_wallet_change_password": nav.ShowWalletChangePassword,
		"on_wallet_show_seed":       nav.ShowWalletShowSeed,
		"on_wallet_set_default_fee": nav.ShowWalletSetDefaultFee,
		"on_about_gtk":              nav.ShowAboutGtk,
		"on_about":                  nav.ShowAbout,
		"on_open_explorer":          nav.BrowseExplorer,
		"on_open_website":           nav.BrowseWebsite,
		"on_open_docs":              nav.BrowseDocs,
	})
}
