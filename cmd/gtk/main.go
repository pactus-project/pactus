//go:build gtk

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/gofrs/flock"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	gtkapp "github.com/pactus-project/pactus/cmd/gtk/app"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/signal"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/version"
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
	version.NodeAgent.AppType = "gui"

	if runtime.GOOS == "darwin" {
		// Changing the PANGOCAIRO_BACKEND is necessary on MacOS to render emoji
		_ = os.Setenv("PANGOCAIRO_BACKEND", "fontconfig")
	}

	gtk.Init(nil)
}

func main() {
	flag.Parse()

	// The gtk should run on main thread.
	runtime.UnlockOSThread()
	runtime.LockOSThread()

	// Create a new app.
	app, err := gtk.ApplicationNew(appID, glib.APPLICATION_NON_UNIQUE)
	gtkutil.FatalErrorCheck(err)

	settings, err := gtk.SettingsGetDefault()
	gtkutil.FatalErrorCheck(err)

	err = settings.Object.Set("gtk-application-prefer-dark-theme", true)
	gtkutil.FatalErrorCheck(err)

	assets.InitAssets()

	workingDir, err := filepath.Abs(*workingDirOpt)
	if err != nil {
		terminal.PrintErrorMsgf("Aborted! %v", err)

		return
	}

	// If node is not initialized yet
	if util.IsDirNotExistsOrEmpty(workingDir) {
		network := genesis.Mainnet
		if *testnetOpt {
			network = genesis.Testnet
		}
		if !startupAssistant(context.Background(), workingDir, network) {
			return
		}
	}

	// Define the lock file path
	lockFilePath := filepath.Join(workingDir, ".pactus.lock")
	fileLock := flock.New(lockFilePath)

	locked, err := fileLock.TryLock()
	gtkutil.FatalErrorCheck(err)

	if !locked {
		terminal.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

		return
	}
	var guiNode *node.Node

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	var gui *gtkapp.GUI
	activateOnce := new(sync.Once)
	shutdownOnce := new(sync.Once)

	shutdown := func() {
		shutdownOnce.Do(func() {
			if gui != nil {
				gui.Cleanup()
			}
			if guiNode != nil {
				guiNode.Stop()
			}
			_ = fileLock.Unlock()
		})
	}

	// Connect function to application shutdown event, this is not required.
	app.Connect("shutdown", func() {
		log.Println("Application shutdown")
		shutdown()
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		activateOnce.Do(func() {
			log.Println("application activate")

			splash := view.NewSplashWindow(app)
			splash.SetVersion(version.NodeVersion().StringWithAlias())
			gtk.WindowSetAutoStartupNotification(false)
			splash.ShowAll()
			gtk.WindowSetAutoStartupNotification(true)
			app.AddWindow(splash.Window())

			notify := func(msg string) {
				glib.IdleAdd(func() bool {
					splash.SetStatus(msg)

					return false
				})
			}

			go func() {
				n, err := newNode(workingDir, notify)

				glib.IdleAdd(func() bool {
					if err != nil {
						splash.Destroy()
						shutdown()
						gtkutil.ShowError(err)
						app.Quit()

						return false
					}

					guiNode = n
					splash.SetStatus("Loading wallet interface...")

					gui, err = gtkapp.Run(guiNode, app)
					gtkutil.FatalErrorCheck(err)

					splash.Destroy()

					return false
				})
			}()
		})
	})

	signal.HandleInterrupt(func() {
		terminal.PrintInfoMsgf("Exiting...")
		shutdown()
	})

	// Launch the application
	os.Exit(app.Run(nil))
}

type statusReporter func(string)

func reportStatus(cb statusReporter, msg string) {
	if cb != nil {
		cb(msg)
	}
}

func newNode(workingDir string, statusCb statusReporter) (*node.Node, error) {
	// change working directory
	if err := os.Chdir(workingDir); err != nil {
		log.Println("Aborted! Unable to changes working directory. " + err.Error())

		return nil, err
	}

	reportStatus(statusCb, "Opening wallet...")
	passwordFetcher := func() (string, bool) {
		if *passwordOpt != "" {
			return *passwordOpt, true
		}

		var (
			pwd string
			ok  bool
		)
		done := make(chan struct{})
		glib.IdleAdd(func() bool {
			pwd, ok = gtkapp.PromptWalletPassword()
			close(done)

			return false
		})
		<-done

		return pwd, ok
	}

	reportStatus(statusCb, "Starting node services...")
	n, err := cmd.StartNode(workingDir, passwordFetcher, nil)
	if err != nil {
		return nil, err
	}

	return n, nil
}
