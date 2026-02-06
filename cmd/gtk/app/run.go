//go:build gtk

package app

import (
	"context"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/cmd/gtk/controller"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"google.golang.org/grpc"
)

type GUI struct {
	MainWindow    *view.MainWindowView
	NodeCtrl      *controller.NodeWidgetController
	WalletCtrl    *controller.WalletWidgetController
	ValidatorCtrl *controller.ValidatorWidgetController
	grpcConn      grpc.ClientConnInterface
}

// Run builds and shows the main window, wiring views/controllers.
// It accepts a gRPC connection to the node (standard grpc.ClientConn or gRPC-Web).
// It returns a cleanup function that closes the window and stops timers.
func Run(ctx context.Context, conn grpc.ClientConnInterface,
	gtkApp *gtk.Application, notify func(string),
) (*GUI, error) {
	mwView := view.NewMainWindowView()

	blockchainClient := pactus.NewBlockchainClient(conn)
	transactionClient := pactus.NewTransactionClient(conn)
	networkClient := pactus.NewNetworkClient(conn)
	walletClient := pactus.NewWalletClient(conn)

	nodeModel := model.NewNodeModel(ctx, blockchainClient, networkClient)
	validatorModel := model.NewValidatorModel(ctx, blockchainClient)
	walletModel := model.NewWalletModel(ctx, walletClient, transactionClient, blockchainClient, cmd.DefaultWalletName)

	nodeView := view.NewNodeWidgetView()
	walletView := view.NewWalletWidgetView()
	validatorView := view.NewValidatorWidgetView()

	nodeCtrl := controller.NewNodeWidgetController(nodeView, nodeModel)
	walletCtrl := controller.NewWalletWidgetController(walletView, walletModel)
	validatorCtrl := controller.NewValidatorWidgetController(validatorView, validatorModel)

	nav := controller.NewNavigator(gtkApp, walletModel, walletCtrl)

	notify("Gathering Node info...")
	if err := nodeCtrl.BuildView(ctx); err != nil {
		return nil, err
	}

	notify("Gathering Validators info...")
	if err := validatorCtrl.BuildView(ctx); err != nil {
		return nil, err
	}

	notify("Gathering Wallet info...")
	if err := walletCtrl.BuildView(ctx, nav); err != nil {
		return nil, err
	}

	mwView.BoxNode.Add(nodeView.Box)
	mwView.BoxDefaultWallet.Add(walletView.Box)
	mwView.BoxValidators.Add(validatorView.Box)

	mwCtrl := controller.NewMainWindowController(mwView)
	mwCtrl.BuildView(nav)

	gtkutil.IdleAddSync(func() {
		mwView.Window.ShowAll()
		gtkApp.AddWindow(mwView.Window)
	})

	return &GUI{
		MainWindow:    mwView,
		NodeCtrl:      nodeCtrl,
		WalletCtrl:    walletCtrl,
		ValidatorCtrl: validatorCtrl,
		grpcConn:      conn,
	}, nil
}

func (g *GUI) Cleanup() {
	g.MainWindow.Cleanup()
}
