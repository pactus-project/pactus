package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/node"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/wallet"
)

var workingDir *string

const appId = "com.github.zarb"

func init() {
	workingDir = flag.String("working-dir", cmd.ZarbHomeDir(), "working directory")
}
func main() {
	flag.Parse()

	var err error
	workspacePath, err := filepath.Abs(*workingDir)
	if err != nil {
		cmd.PrintErrorMsg("Aborted! %v", err)
		return
	}

	// If node is not initialized yet
	if util.IsDirNotExistsOrEmpty(workspacePath) {
		if !startupAssistant(workspacePath) {
			return
		}
	}

	// Create a new app.
	// When using GtkApplication, it is not necessary to call gtk_init() manually.
	app, err := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	errorCheck(err)

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		log.Println("application activate")
		start(workspacePath, app)
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(nil))
}

func startingNode(workspacePath string, wallet *wallet.Wallet, password string) (*node.Node, *time.Time, error) {
	addresses := wallet.Addresses()
	valPrvKeyStr, err := wallet.PrivateKey(password, addresses[0].Address)
	if err != nil {
		return nil, nil, err
	}
	prv, err := bls.PrivateKeyFromString(valPrvKeyStr)
	if err != nil {
		return nil, nil, err
	}
	signer := crypto.NewSigner(prv)

	gen, err := genesis.LoadFromFile(cmd.ZarbGenesisPath(workspacePath))
	if err != nil {
		return nil, nil, err
	}
	conf, err := config.LoadFromFile(cmd.ZarbConfigPath(workspacePath))
	if err != nil {
		return nil, nil, err
	}
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

func start(workspacePath string, app *gtk.Application) {
	// change working directory
	if err := os.Chdir(workspacePath); err != nil {
		log.Println("Aborted! Unable to changes working directory. " + err.Error())
		return
	}

	time.Sleep(1 * time.Second)

	path := cmd.ZarbDefaultWalletPath(workspacePath)
	wallet, err := wallet.OpenWallet(path)
	errorCheck(err)

	password, ok := getWalletPassword(nil, wallet)
	if !ok {
		showInfoDialog("Canceled!")
		return
	}
	node, genTime, err := startingNode(workspacePath, wallet, password)
	errorCheck(err)

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

	walletModel.rebuildModel()

	app.AddWindow(win)
}
