//go:build gtk

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	flock "github.com/gofrs/flock"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
	"github.com/pactus-project/pactus/wallet"
)

const appID = "com.github.pactus-project.pactus.pactus-gui"

var (
	workingDirOpt *string
	passwordOpt   *string
	testnetOpt    *bool
)

func init() {
	workingDirOpt = flag.String("working-dir", cmd.PactusDefaultHomeDir(), "working directory path")
	passwordOpt = flag.String("password", "", "wallet password")
	testnetOpt = flag.Bool("testnet", false, "initializing for the testnet")
	version.AppType = "gui"
	gtk.Init(nil)
}

func main() {
	flag.Parse()

	var err error
	workingDir, err := filepath.Abs(*workingDirOpt)
	if err != nil {
		cmd.PrintErrorMsgf("Aborted! %v", err)

		return
	}

	// If node is not initialized yet
	if !util.PathExists(workingDir) {
		network := genesis.Mainnet
		if *testnetOpt {
			network = genesis.Testnet
		}
		if !startupAssistant(workingDir, network) {
			return
		}
	}

	// Define the lock file path
	lockFilePath := filepath.Join(workingDir, ".pactus.lock")
	fileLock := flock.New(lockFilePath)

	locked, err := fileLock.TryLock()
	if err != nil {
		// handle unable to attempt to acquire lock
		fatalErrorCheck(err)
	}

	if !locked {
		cmd.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

		return
	}

	// Create a new app.
	// When using GtkApplication, it is not necessary to call gtk_init() manually.
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_NON_UNIQUE)
	fatalErrorCheck(err)

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		log.Println("application activate")

		// Show about dialog as splash screen
		splashDlg := aboutDialog()
		splashDlg.SetDecorated(false)
		splashDlg.SetResizable(false)
		splashDlg.SetPosition(gtk.WIN_POS_CENTER)
		splashDlg.SetTypeHint(gdk.WINDOW_TYPE_HINT_SPLASHSCREEN)

		gtk.WindowSetAutoStartupNotification(false)
		splashDlg.ShowAll()
		gtk.WindowSetAutoStartupNotification(true)

		app.AddWindow(splashDlg)

		// This might also force GTK to draw the splash screen
		for gtk.EventsPending() {
			gtk.MainIteration()
		}

		// Running the start-up logic in a separate goroutine
		glib.TimeoutAdd(uint(100), func() bool {
			start(workingDir, app)
			splashDlg.Destroy()

			// Ensures the function is not called again
			return false
		})
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		_ = fileLock.Unlock()
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(nil))
}

func start(workingDir string, app *gtk.Application) {
	// change working directory
	if err := os.Chdir(workingDir); err != nil {
		log.Println("Aborted! Unable to changes working directory. " + err.Error())

		return
	}

	passwordFetcher := func(wlt *wallet.Wallet) (string, bool) {
		if *passwordOpt != "" {
			return *passwordOpt, true
		}

		return getWalletPassword(wlt)
	}
	node, wlt, err := cmd.StartNode(workingDir, passwordFetcher)
	fatalErrorCheck(err)

	grpcAddr := node.GRPC().Address()
	cmd.PrintInfoMsgf("connect wallet to grpc server: %s\n", grpcAddr)

	err = wlt.Connect(grpcAddr)
	fatalErrorCheck(err)

	nodeModel := newNodeModel(node)
	walletModel := newWalletModel(wlt, node)

	// building main window
	win := buildMainWindow(nodeModel, walletModel)

	// Show the Window and all of its components.
	win.ShowAll()

	walletModel.rebuildModel()

	app.AddWindow(win)
}
