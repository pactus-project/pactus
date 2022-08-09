//go:build gtk

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/node"
	"github.com/zarbchain/zarb-go/node/config"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/genesis"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/wallet"
)

const appID = "com.github.zarbchain.zarb-go.zarb-gui"

var (
	workingDirOpt *string
	testnetOpt    *bool
)

func init() {
	workingDirOpt = flag.String("working-dir", cmd.ZarbHomeDir(), "working directory")
	testnetOpt = flag.Bool("testnet", true, "working directory") // TODO: make it false after mainnet launch
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
	if !util.PathExists(cmd.ZarbDefaultWalletPath(workingDir)) {
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
		start(nil, workingDir, app)
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(nil))
}

func startingNode(workingDir string, wallet *wallet.Wallet, password string) (*node.Node, *time.Time, error) {
	gen, err := genesis.LoadFromFile(cmd.ZarbGenesisPath(workingDir))
	if err != nil {
		return nil, nil, err
	}
	if gen.Params().IsTestnet() {
		crypto.AddressHRP = "tzc"
		crypto.PublicKeyHRP = "tpublic"
		crypto.PrivateKeyHRP = "tsecret"
		crypto.XPublicKeyHRP = "txpublic"
		crypto.XPrivateKeyHRP = "txsecret"
	}

	conf, err := config.LoadFromFile(cmd.ZarbConfigPath(workingDir))
	if err != nil {
		return nil, nil, err
	}

	addrInfos := wallet.AddressLabels()
	if len(addrInfos) == 0 {
		return nil, nil, fmt.Errorf("validator address is not defined")
	}
	prv, err := wallet.PrivateKey(password, addrInfos[0].Address)
	if err != nil {
		return nil, nil, err
	}
	signer := crypto.NewSigner(prv)

	node, err := node.NewNode(gen, conf, signer)
	if err != nil {
		return nil, nil, err
	}
	//TODO: log to file

	err = node.Start()
	if err != nil {
		return nil, nil, err
	}

	genTime := gen.GenesisTime()
	return node, &genTime, nil
}

func start(parent gtk.IWindow, workingDir string, app *gtk.Application) {
	//theme, _ := gtk.IconThemeGetDefault()
	//theme.AddResourcePath("icons")

	// change working directory
	if err := os.Chdir(workingDir); err != nil {
		log.Println("Aborted! Unable to changes working directory. " + err.Error())
		return
	}

	time.Sleep(1 * time.Second)

	path := cmd.ZarbDefaultWalletPath(workingDir)
	wallet, err := wallet.OpenWallet(path, false)
	fatalErrorCheck(err)

	password, ok := getWalletPassword(wallet)
	if !ok {
		showInfoDialog(parent, "Canceled!")
		return
	}
	// TODO: Get genTime from the node or state
	node, genTime, err := startingNode(workingDir, wallet, password)
	fatalErrorCheck(err)

	// TODO
	// No showing the main window
	if err != nil {
		return
	}

	nodeModel := newNodeModel(node)
	walletModel := newWalletModel(wallet)

	// building main window
	win := buildMainWindow(nodeModel, walletModel, *genTime)

	// Show the Window and all of its components.
	win.Show()

	err = walletModel.rebuildModel()
	fatalErrorCheck(err)

	app.AddWindow(win)
}
