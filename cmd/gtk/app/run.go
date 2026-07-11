//go:build gtk

package app

import (
	"context"

	adw "github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/gogpu/systray"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
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
	NetworkCtrl   *controller.NetworkWidgetController

	tray *systray.SystemTray
}

// Run builds and shows the main window, wiring views/controllers.
// It accepts a gRPC connection to the node (standard grpc.ClientConn or gRPC-Web).
// connectionLabel is "Remote address" or "Working directory"; connectionValue is the address or path.
// workingDir is the local node data directory; isLocal is true when running a local node.
// It returns a cleanup function that closes the window and stops timers.
func Run(ctx context.Context, conn grpc.ClientConnInterface,
	gtkApp *gtk.Application, notify func(string),
	connectionLabel, connectionValue, workingDir string, isLocal bool,
) (*GUI, error) {
	blockchainClient := pactus.NewBlockchainClient(conn)
	transactionClient := pactus.NewTransactionClient(conn)
	networkClient := pactus.NewNetworkClient(conn)
	walletClient := pactus.NewWalletClient(conn)

	nodeModel := model.NewNodeModel(ctx, blockchainClient, networkClient)
	validatorModel := model.NewValidatorModel(ctx, blockchainClient)
	walletModel := model.NewWalletModel(ctx, walletClient, transactionClient, blockchainClient, cmd.DefaultWalletName)
	committeeModel := model.NewCommitteeModel(ctx, blockchainClient)
	networkModel := model.NewNetworkModel(ctx, networkClient)

	nodeView := gtkutil.IdleAddSyncT(view.NewNodeWidgetView)
	walletView := gtkutil.IdleAddSyncT(view.NewWalletWidgetView)
	validatorView := gtkutil.IdleAddSyncT(view.NewValidatorWidgetView)
	committeeView := gtkutil.IdleAddSyncT(view.NewCommitteeWidgetView)
	networkView := gtkutil.IdleAddSyncT(view.NewNetworkWidgetView)

	nodeCtrl := controller.NewNodeWidgetController(nodeView, nodeModel)
	walletCtrl := controller.NewWalletWidgetController(walletView, walletModel)
	validatorCtrl := controller.NewValidatorWidgetController(validatorView, validatorModel)
	committeeCtrl := controller.NewCommitteeWidgetController(committeeView, committeeModel)
	networkCtrl := controller.NewNetworkWidgetController(networkView, networkModel)
	configModel, err := model.NewConfigModel(workingDir, isLocal)
	if err != nil {
		return nil, err
	}

	nav := controller.NewNavigator(gtkApp, walletModel, walletCtrl, configModel)

	notify("Fetching Node info...")
	err = nodeCtrl.BuildView(ctx, connectionLabel, connectionValue)
	if err != nil {
		return nil, err
	}

	notify("Fetching Validators info...")
	err = validatorCtrl.BuildView(ctx)
	if err != nil {
		return nil, err
	}

	notify("Fetching Committee info...")
	err = committeeCtrl.BuildView(ctx)
	if err != nil {
		return nil, err
	}

	notify("Fetching Network info...")
	err = networkCtrl.BuildView(ctx)
	if err != nil {
		return nil, err
	}

	notify("Fetching Wallet info...")
	err = walletCtrl.BuildView(ctx, nav)
	if err != nil {
		return nil, err
	}

	mwView := gtkutil.IdleAddSyncT(func() *view.MainWindowView {
		mwView := view.NewMainWindowView()

		// Register the main window so dialogs open centered over it.
		gtkutil.SetMainWindow(&mwView.Window.Window)

		// Custom title bar: small Pactus logo next to the app name on the left,
		// and a round light/dark toggle on the right, beside the window buttons.
		header := gtk.NewHeaderBar()

		brand := gtk.NewBox(gtk.OrientationHorizontal, 8)
		brand.SetVAlign(gtk.AlignCenter)
		logo := gtk.NewImageFromPaintable(assets.ImagePactusLogoTexture)
		logo.SetPixelSize(20)
		logo.AddCSSClass("app-logo")
		title := gtk.NewLabel("Pactus GUI")
		title.AddCSSClass("app-title")
		brand.Append(logo)
		brand.Append(title)
		header.PackStart(brand)
		// Suppress the centered window title so the brand stays left-aligned.
		header.SetTitleWidget(gtk.NewLabel(""))

		themeToggle := gtk.NewToggleButton()
		themeToggle.AddCSSClass("theme-toggle")
		themeToggle.AddCSSClass("circular")
		themeToggle.SetVAlign(gtk.AlignCenter)
		themeToggle.SetTooltipText("Toggle light / dark mode")
		applyThemeIcon := func(dark bool) {
			if dark {
				themeToggle.SetLabel("🌙")
			} else {
				themeToggle.SetLabel("☀")
			}
		}
		// Restore the persisted light/dark choice, falling back to the system
		// theme when the user has not chosen yet.
		startDark := adw.StyleManagerGetDefault().Dark()
		if saved, ok := gtkutil.LoadDarkMode(); ok {
			startDark = saved
			nav.SetDarkMode(saved)
		}
		themeToggle.SetActive(startDark)
		applyThemeIcon(startDark)
		themeToggle.ConnectToggled(func() {
			active := themeToggle.Active()
			nav.SetDarkMode(active)
			applyThemeIcon(active)
		})
		header.PackEnd(themeToggle)

		mwView.Window.SetTitlebar(header)

		walletCtrl.SetupMenu(mwView.Window)

		menu := nav.CreateMenu(isLocal)
		gtkApp.SetMenubar(menu)
		mwView.Window.SetShowMenubar(true)

		mwView.BoxNode.Append(nodeView.Box)
		mwView.BoxWallet.Append(walletView.Box)
		mwView.BoxValidators.Append(validatorView.BoxValidators)
		mwView.BoxCommittee.Append(committeeView.Box)
		mwView.BoxNetwork.Append(networkView.Box)

		// Build controller
		mwCtrl := controller.NewMainWindowController(mwView)
		mwCtrl.BuildView()

		gtkApp.AddWindow(&mwView.Window.Window)
		mwView.Window.Present()

		// Center the main window on first show; GTK4 cannot position it, so do
		// it natively once mapped (no-op on Linux/macOS).
		glib.TimeoutAdd(80, func() bool {
			gtkutil.CenterActiveWindow()

			return false
		})

		return mwView
	})

	// Create the system tray icon.
	trayIcon := createTray(mwView, gtkApp)

	return &GUI{
		MainWindow:    mwView,
		NodeCtrl:      nodeCtrl,
		WalletCtrl:    walletCtrl,
		ValidatorCtrl: validatorCtrl,
		CommitteeCtrl: committeeCtrl,
		NetworkCtrl:   networkCtrl,
		tray:          trayIcon,
	}, nil
}

func (g *GUI) Cleanup() {
	g.MainWindow.Cleanup()
	if g.tray != nil {
		g.tray.Remove()
	}
}

// createTray creates a system tray icon with Show and Exit menu items.
// The tray runs its event loop in a background goroutine.
func createTray(mwView *view.MainWindowView, gtkApp *gtk.Application) *systray.SystemTray {
	trayIcon := systray.New()

	menu := systray.NewMenu()
	menu.Add("Show", func() {
		gtkutil.IdleAddSync(func() {
			mwView.Window.Present()
		})
	})
	menu.AddSeparator()
	menu.Add("Exit", func() {
		gtkutil.IdleAddSync(func() {
			mwView.HideOnClose = false
			gtkApp.Quit()
		})
	})

	trayIcon.
		SetIcon(assets.ImagePactusTrayDark32Data).
		SetTooltip("Pactus GUI").
		SetMenu(menu).
		OnClick(func() {
			gtkutil.IdleAddSync(func() {
				mwView.Window.Present()
			})
		})

	trayIcon.Show()

	// Run the tray event loop in a background goroutine.
	go func() {
		if err := trayIcon.Run(); err != nil {
			gtkutil.Logf("System tray error: %v", err)
		}
	}()

	return trayIcon
}
