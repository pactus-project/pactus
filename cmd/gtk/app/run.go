//go:build gtk

package app

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/cmd/gtk/controller"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/node"
)

type GUI struct {
	MainWindow *view.MainWindowView
	NodeCtrl   *controller.NodeWidgetController
	WalletCtrl *controller.WalletWidgetController
}

// Run builds and shows the main window, wiring views/controllers.
// It returns a cleanup function that closes the window and stops timers.
func Run(n *node.Node, gtkApp *gtk.Application) (*GUI, error) {
	mwView, err := view.NewMainWindowView()
	if err != nil {
		return nil, err
	}

	nodeView, err := view.NewNodeWidgetView()
	if err != nil {
		return nil, err
	}
	nodeCtrl := controller.NewNodeWidgetController(nodeView, n)
	if err := nodeCtrl.Bind(); err != nil {
		return nil, err
	}

	walletModel, err := model.NewWalletModel(n, cmd.DefaultWalletName)
	if err != nil {
		return nil, err
	}

	nav := NewNavigator(walletModel)

	walletView, err := view.NewWalletWidgetView()
	if err != nil {
		return nil, err
	}

	walletCtrl := controller.NewWalletWidgetController(walletView, walletModel)

	walletCtrl.Bind(controller.WalletWidgetHandlers{
		OnNewAddress: func() {
			nav.ShowCreateAddress()
			walletCtrl.RefreshAddresses()
		},
		OnSetDefaultFee: func() {
			nav.ShowSetDefaultFee()
			walletCtrl.RefreshInfo()
		},
		OnChangePassword: func() {
			nav.ShowChangePassword()
			walletCtrl.RefreshInfo()
		},
		OnShowSeed: func() {
			nav.ShowSeed()
		},
		OnUpdateLabel: func(address string) {
			nav.ShowUpdateLabel(address)
			walletCtrl.RefreshAddresses()
		},
		OnShowDetails: func(address string) {
			nav.ShowAddressDetails(address)
		},
		OnShowPrivateKey: func(address string) {
			nav.ShowPrivateKey(address)
		},
	})

	mwView.BoxNode.Add(nodeView.Box)
	mwView.BoxDefaultWallet.Add(walletView.Box)

	mwCtrl := controller.NewMainWindowController(mwView)
	mwCtrl.Bind(&controller.MainWindowHandlers{
		OnAboutGtk: func() {
			nav.ShowAboutGTK()
		},
		OnAbout: func() {
			nav.ShowAbout()
		},
		OnQuit: func() {
			gtkApp.Quit()
		},
		OnTransactionTransfer: func() {
			nav.ShowTransferTx()
			walletCtrl.RefreshTransactions()
		},
		OnTransactionBond: func() {
			nav.ShowBondTx()
			walletCtrl.RefreshTransactions()
		},
		OnTransactionUnbond: func() {
			nav.ShowUnbondTx()
			walletCtrl.RefreshTransactions()
		},
		OnTransactionWithdraw: func() {
			nav.ShowWithdrawTx()
			walletCtrl.RefreshTransactions()
		},
		OnWalletNewAddress: func() {
			nav.ShowCreateAddress()
			walletCtrl.RefreshAddresses()
		},
		OnWalletChangePassword: func() {
			nav.ShowChangePassword()
			walletCtrl.RefreshInfo()
		},
		OnWalletShowSeed: func() {
			nav.ShowSeed()
		},
		OnWalletSetDefaultFee: func() {
			nav.ShowSetDefaultFee()
			walletCtrl.RefreshInfo()
		},
		OnMenuActivateWebsite: func() {
			_ = gtkutil.OpenURLInBrowser("https://pactus.org/")
		},
		OnMenuActivateExplorer: func() {
			_ = gtkutil.OpenURLInBrowser("https://pacviewer.com/")
		},
		OnMenuActivateDocs: func() {
			_ = gtkutil.OpenURLInBrowser("https://docs.pactus.org/")
		},
	})

	mwView.Window.ShowAll()
	gtkApp.AddWindow(mwView.Window)

	return &GUI{
		MainWindow: mwView,
		NodeCtrl:   nodeCtrl,
		WalletCtrl: walletCtrl,
	}, nil
}

func (g *GUI) Cleanup() {
	g.NodeCtrl.Cleanup()
	g.WalletCtrl.Cleanup()
	g.MainWindow.Cleanup()
}
