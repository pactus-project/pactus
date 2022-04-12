package main

import (
	_ "embed"
	"errors"
	"flag"
	"log"
	"os"
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

	if util.IsDirNotExistsOrEmpty(*workingDir) {
		if !startupAssistant() {
			return
		}
	}

	// Create a new application.
	// When using GtkApplication, it is not necessary to call gtk_init() manually.
	application, err := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	errorCheck(err)

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

	// Connect function to application startup event, this is not required.
	application.Connect("startup", func() {
		log.Println("application startup")
	})

	// Connect function to application activate event
	application.Connect("activate", func() {
		log.Println("application activate")

		// Get the GtkBuilder UI definition in the glade file.
		builder, err := gtk.BuilderNewFromString(string(uiMainWindow))
		errorCheck(err)

		// Map the handlers to callback functions, and connect the signals
		// to the Builder.
		signals := map[string]interface{}{
			"on_main_window_destroy": onMainWindowDestroy,
		}
		builder.ConnectSignals(signals)

		// Get the object with the id of "main_window".
		obj, err := builder.GetObject("main_window")
		errorCheck(err)

		// Verify that the object is a pointer to a gtk.Window.
		win, err := isWindow(obj)
		errorCheck(err)

		// Show the Window and all of its components.
		win.Show()
		application.AddWindow(win)
	})

	// Connect function to application shutdown event, this is not required.
	application.Connect("shutdown", func() {
		log.Println("application shutdown")
	})

	// Launch the application
	os.Exit(application.Run(nil))
}

func isWindow(obj glib.IObject) (*gtk.Window, error) {
	// Make type assertion (as per gtk.go).
	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func isDialog(obj glib.IObject) (*gtk.Dialog, error) {
	// Make type assertion (as per gtk.go).
	if dlg, ok := obj.(*gtk.Dialog); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.Dialog")
}

func isEntry(obj glib.IObject) (*gtk.Entry, error) {
	// Make type assertion (as per gtk.go).
	if dlg, ok := obj.(*gtk.Entry); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.Entry")
}

// onMainWindowDestory is the callback that is linked to the
// on_main_window_destroy handler. It is not required to map this,
// and is here to simply demo how to hook-up custom callbacks.
func onMainWindowDestroy() {
	log.Println("onMainWindowDestroy")
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
