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
	CommitteeCtrl *controller.CommitteeWidgetController
}

// Run builds and shows the main window, wiring views/controllers.
// It accepts a gRPC connection to the node (standard grpc.ClientConn or gRPC-Web).
// connectionLabel is "Remote address" or "Working directory"; connectionValue is the address or path.
// It returns a cleanup function that closes the window and stops timers.
func Run(ctx context.Context, conn grpc.ClientConnInterface,
	gtkApp *gtk.Application, notify func(string),
	connectionLabel, connectionValue string,
) (*GUI, error) {
	blockchainClient := pactus.NewBlockchainClient(conn)
	transactionClient := pactus.NewTransactionClient(conn)
	networkClient := pactus.NewNetworkClient(conn)
	walletClient := pactus.NewWalletClient(conn)

	nodeModel := model.NewNodeModel(ctx, blockchainClient, networkClient)
	validatorModel := model.NewValidatorModel(ctx, blockchainClient)
	walletModel := model.NewWalletModel(ctx, walletClient, transactionClient, blockchainClient, cmd.DefaultWalletName)
	committeeModel := model.NewCommitteeModel(ctx, blockchainClient)

	nodeView := gtkutil.IdleAddSyncT(view.NewNodeWidgetView)
	walletView := gtkutil.IdleAddSyncT(view.NewWalletWidgetView)
	validatorView := gtkutil.IdleAddSyncT(view.NewValidatorWidgetView)
	committeeView := gtkutil.IdleAddSyncT(view.NewCommitteeWidgetView)

	nodeCtrl := controller.NewNodeWidgetController(nodeView, nodeModel)
	walletCtrl := controller.NewWalletWidgetController(walletView, walletModel)
	validatorCtrl := controller.NewValidatorWidgetController(validatorView, validatorModel)
	committeeCtrl := controller.NewCommitteeWidgetController(committeeView, committeeModel)

	nav := controller.NewNavigator(gtkApp, walletModel, walletCtrl)

	notify("Fetching Node info...")
	err := nodeCtrl.BuildView(ctx, connectionLabel, connectionValue)
	if err != nil {
		return nil, err
	}

	notify("Fetching Validators info...")
	err = validatorCtrl.BuildView(ctx)
	if err != nil {
		return nil, err
	}

	notify("Fetching Wallet info...")
	err = walletCtrl.BuildView(ctx, nav)
	if err != nil {
		return nil, err
	}

	notify("Fetching Committee info...")
	err = committeeCtrl.BuildView(ctx)
	if err != nil {
		return nil, err
	}

	mwView := gtkutil.IdleAddSyncT(func() *view.MainWindowView {
		mwView := view.NewMainWindowView()

		mwView.BoxNode.Add(nodeView.Box)
		mwView.BoxDefaultWallet.Add(walletView.Box)
		mwView.BoxValidators.Add(validatorView.Box)
		mwView.BoxCommittee.Add(committeeView.Box)

		mwCtrl := controller.NewMainWindowController(mwView)
		mwCtrl.BuildView(nav)

		mwView.Window.ShowAll()
		gtkApp.AddWindow(mwView.Window)

		return mwView
	})

	return &GUI{
		MainWindow:    mwView,
		NodeCtrl:      nodeCtrl,
		WalletCtrl:    walletCtrl,
		ValidatorCtrl: validatorCtrl,
		CommitteeCtrl: committeeCtrl,
	}, nil
}

func (g *GUI) Cleanup() {
	g.MainWindow.Cleanup()
}
