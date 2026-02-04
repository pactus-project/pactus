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
	"time"

	"github.com/ezex-io/gopkg/signal"
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
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/version"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const appID = "com.github.pactus-project.pactus.pactus-gui"

var (
	workingDirOpt *string
	passwordOpt   *string
	testnetOpt    *bool
	grpcAddrOpt   *string
)

func init() {
	workingDirOpt = flag.String("working-dir", cmd.PactusDefaultHomeDir(), "working directory path")
	passwordOpt = flag.String("password", "", "wallet password")
	testnetOpt = flag.Bool("testnet", false, "initializing for the testnet")
	grpcAddrOpt = flag.String("grpc-addr", "", "connect to remote gRPC server instead of starting local node")
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

	ctx, cancel := context.WithCancel(context.Background())

	var fileLock *flock.Flock
	if *grpcAddrOpt == "" {
		// If node is not initialized yet
		if util.IsDirNotExistsOrEmpty(workingDir) {
			network := genesis.Mainnet
			if *testnetOpt {
				network = genesis.Testnet
			}
			if !startupAssistant(ctx, workingDir, network) {
				return
			}
		}

		// Define the lock file path
		lockFilePath := filepath.Join(workingDir, ".pactus.lock")
		fileLock = flock.New(lockFilePath)

		locked, err := fileLock.TryLock()
		gtkutil.FatalErrorCheck(err)

		if !locked {
			terminal.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
		}
	}
	var guiNode *node.Node

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		log.Println("application startup")
	})

	var gui *gtkapp.GUI
	var grpcConn *grpc.ClientConn
	activateOnce := new(sync.Once)
	shutdownOnce := new(sync.Once)

	shutdown := func() {
		shutdownOnce.Do(func() {
			cancel()
			if grpcConn != nil {
				_ = grpcConn.Close()
			}
			if gui != nil {
				gui.Cleanup()
			}
			if guiNode != nil {
				guiNode.Stop()
			}
			if fileLock != nil {
				_ = fileLock.Unlock()
			}
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
				grpcAddr := *grpcAddrOpt
				if grpcAddr == "" {
					time.Sleep(1 * time.Second)
					reportStatus(notify, "Starting local node...")
					guiNode, err = newNode(ctx, workingDir, notify)

					grpcAddr = guiNode.GRPC().Address()
				} else {
					time.Sleep(1 * time.Second)
					reportStatus(notify, "Connecting to remote node...")
				}

				glib.IdleAdd(func() bool {
					grpcConn, err = grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
					gtkutil.FatalErrorCheck(err)

					gui, err = gtkapp.Run(ctx, grpcConn, app)
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

func newNode(ctx context.Context, workingDir string, statusCb statusReporter) (*node.Node, error) {
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
	n, err := cmd.StartNode(ctx, workingDir, passwordFetcher, nil)
	if err != nil {
		return nil, err
	}

	return n, nil
}
