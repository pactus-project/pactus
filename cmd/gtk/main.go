//go:build gtk

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

const appID = "com.github.pactus-project.pactus.pactus-gui"

var (
	workingDirOpt *string
	passwordOpt   *string
	testnetOpt    *bool
)

func init() {
	workingDirOpt = flag.String("working-dir", cmd.PactusHomeDir(), "working directory path")
	passwordOpt = flag.String("password", "", "wallet password")
	testnetOpt = flag.Bool("testnet", true, "initializing for the testnet") // TODO: make it false after mainnet launch

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
		start(workingDir, app)
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
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

	passwordFetcher := func(wallet *wallet.Wallet) (string, bool) {
		if *passwordOpt != "" {
			return *passwordOpt, true
		}
		return getWalletPassword(wallet)
	}

	node, wallet, err := cmd.StartNode(workingDir, passwordFetcher)
	fatalErrorCheck(err)

	grpcAddr := node.GRPC().Address()
	fmt.Printf("connect wallet to grpc server: %s\n", grpcAddr)

	err = wallet.Connect(grpcAddr)
	fatalErrorCheck(err)

	nodeModel := newNodeModel(node)
	walletModel := newWalletModel(wallet)

	// building main window
	win := buildMainWindow(nodeModel, walletModel)

	// Show the Window and all of its components.
	win.Show()

	walletModel.rebuildModel()

	app.AddWindow(win)
}
