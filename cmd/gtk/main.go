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
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/node/config"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

const appID = "com.github.pactus-project.pactus.pactus-gui"

var (
	workingDirOpt *string
	testnetOpt    *bool
)

func init() {
	workingDirOpt = flag.String("working-dir", cmd.PactusHomeDir(), "working directory")
	testnetOpt = flag.Bool("testnet", true, "working directory") // TODO: make it false after mainnet launch

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

func startingNode(workingDir string,
	wallet *wallet.Wallet, password string) (*node.Node, *time.Time, error) {
	gen, err := genesis.LoadFromFile(cmd.PactusGenesisPath(workingDir))
	if err != nil {
		return nil, nil, err
	}
	if gen.Params().IsTestnet() {
		crypto.AddressHRP = "tpc"
		crypto.PublicKeyHRP = "tpublic"
		crypto.PrivateKeyHRP = "tsecret"
		crypto.XPublicKeyHRP = "txpublic"
		crypto.XPrivateKeyHRP = "txsecret"
	}

	conf, err := config.LoadFromFile(cmd.PactusConfigPath(workingDir))
	if err != nil {
		return nil, nil, err
	}

	addrInfos := wallet.AddressLabels()
	if len(addrInfos) == 0 {
		return nil, nil, fmt.Errorf("validator address is not defined")
	}

	signers := make([]crypto.Signer, conf.NumValidators)
	rewardAddrs := make([]crypto.Address, conf.NumValidators)
	for i := 0; i < conf.NumValidators; i++ {
		prvKey, err := wallet.PrivateKey(password, addrInfos[i*2].Address)
		fatalErrorCheck(err)

		addr, err := crypto.AddressFromString(addrInfos[(i*2)+1].Address)
		fatalErrorCheck(err)

		signers[i] = crypto.NewSigner(prvKey)
		rewardAddrs[i] = addr
	}

	node, err := node.NewNode(gen, conf, signers, rewardAddrs)
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

func start(workingDir string, app *gtk.Application) {
	// change working directory
	if err := os.Chdir(workingDir); err != nil {
		log.Println("Aborted! Unable to changes working directory. " + err.Error())
		return
	}

	path := cmd.PactusDefaultWalletPath(workingDir)
	wallet, err := wallet.OpenWallet(path, false)
	fatalErrorCheck(err)

	password, ok := getWalletPassword(wallet)
	if !ok {
		showInfoDialog(nil, "Canceled!")
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

	walletModel.rebuildModel()

	app.AddWindow(win)
}
