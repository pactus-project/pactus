package main

import (
	_ "embed"
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

//go:embed ui/main_window.ui
var uiMainWindow []byte

const appId = "com.github.zarb"

func init() {
	workingDir = flag.String("working-dir", cmd.ZarbHomeDir(), "working directory")
}
func main() {
	flag.Parse()

	var err error
	*workingDir, err = filepath.Abs(*workingDir)
	if err != nil {
		cmd.PrintErrorMsg("Aborted! %v", err)
		return
	}

	// If node is not initialized yet
	if util.IsDirNotExistsOrEmpty(*workingDir) {
		if !startupAssistant(*workingDir) {
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
		start(app)
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		log.Println("application activate")
	})

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(app.Run(nil))
}

func startingNode(wallet *wallet.Wallet, password string) (*node.Node, error) {
	addresses := wallet.Addresses()
	valPrvKeyStr, err := wallet.PrivateKey(password, addresses[0].Address)
	if err != nil {
		return nil, err
	}
	prv, err := bls.PrivateKeyFromString(valPrvKeyStr)
	if err != nil {
		return nil, err
	}
	signer := crypto.NewSigner(prv)

	gen, err := genesis.LoadFromFile(cmd.ZarbGenesisPath(*workingDir))
	if err != nil {
		return nil, err
	}
	conf := config.DefaultConfig()
	node, err := node.NewNode(gen, conf, signer)
	if err != nil {
		return nil, err
	}

	err = node.Start()
	if err != nil {
		return nil, err
	}

	return node, nil
}

func start(app *gtk.Application) {

	// change working directory
	if err := os.Chdir(*workingDir); err != nil {
		log.Println("Aborted! Unable to changes working directory. " + err.Error())
		return
	}

	time.Sleep(1 * time.Second)

	path := cmd.ZarbDefaultWalletPath(*workingDir)
	wallet, err := wallet.OpenWallet(path)
	errorCheck(err)

	password, ok := getWalletPassword(nil, wallet)
	if !ok {
		showInfoDialog(nil, "Canceled!")
		return
	}
	_, err = startingNode(wallet, password)
	errorCheck(err)

	// Get the GtkBuilder UI definition in the glade file.
	builder, err := gtk.BuilderNewFromString(string(uiMainWindow))
	errorCheck(err)

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{}
	builder.ConnectSignals(signals)

	// Get the object with the id of "main_window".
	obj, err := builder.GetObject("main_window")
	errorCheck(err)

	// Verify that the object is a pointer to a gtk.Window.
	win, err := isApplicationWindow(obj)
	errorCheck(err)

	// Show the Window and all of its components.
	win.Show()
	app.AddWindow(win)
}
