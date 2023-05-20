//go:build gtk

package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/node"
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
		cmd.PrintErrorMsg("Aborted! %v", err)
		return
	}

	// If node is not initialized yet
	if !util.PathExists(cmd.PactusDefaultWalletPath(workingDir)) {
		if !startupAssistant(workingDir, *testnetOpt) {
			return
		}
	}

	// Create a new app.
	// When using GtkApplication, it is not necessary to call gtk_init() manually.
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
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

func startingNode(workingDir string) (*node.Node, *config.Config, *time.Time, *wallet.Wallet, error) {
	passwordFetcher := func(wallet *wallet.Wallet) (string, bool) {
		if *passwordOpt != "" {
			return *passwordOpt, true
		}
		return getWalletPassword(wallet)
	}

	gen, conf, signers, rewardAddrs, wallet, err := cmd.GetKeys(workingDir, passwordFetcher)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	node, err := node.NewNode(gen, conf, signers, rewardAddrs)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	//TODO: log to file

	err = node.Start()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	genTime := gen.GenesisTime()
	return node, conf, &genTime, wallet, nil
}

func start(workingDir string, app *gtk.Application) {
	// change working directory
	if err := os.Chdir(workingDir); err != nil {
		log.Println("Aborted! Unable to changes working directory. " + err.Error())
		return
	}

	// TODO: Get genTime from the node or state
	node, conf, genTime, wallet, err := startingNode(workingDir)
	fatalErrorCheck(err)

	// TODO
	// No showing the main window
	if err != nil {
		return
	}

	err = wallet.Connect(conf.GRPC.Listen)
	fatalErrorCheck(err)

	nodeModel := newNodeModel(node)
	walletModel := newWalletModel(wallet)

	// building main window
	win := buildMainWindow(nodeModel, walletModel, *genTime)

	// Show the Window and all of its components.
	win.Show()

	walletModel.rebuildModel()

	app.AddWindow(win)
}
