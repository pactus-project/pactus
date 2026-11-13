//go:build gtk

package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/gofrs/flock"
	"github.com/pactus-project/gopkg/signal"
	"github.com/pactus-project/pactus/cmd"
	gtkapp "github.com/pactus-project/pactus/cmd/gtk/app"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/autostart"
	"github.com/pactus-project/pactus/cmd/gtk/controller"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
	"github.com/pactus-project/pactus/www/grpc/basicauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const appID = "com.github.pactus-project.pactus.pactus-gui"

var (
	workingDirOpt   *string
	passwordOpt     *string
	testnetOpt      *bool
	grpcAddrOpt     *string
	grpcInsecureOpt *bool
	grpcUserOpt     *string
	grpcPassFileOpt *string
	snapshotURLOpt  *string
)

func init() {
	workingDirOpt = flag.String("working-dir", cmd.PactusDefaultHomeDir(), "working directory path")
	passwordOpt = flag.String("password", "", "wallet password")
	testnetOpt = flag.Bool("testnet", false, "initializing for the testnet")
	grpcAddrOpt = flag.String("grpc-addr", "", "connect to remote gRPC server instead of starting local node")
	grpcInsecureOpt = flag.Bool("grpc-insecure", false, "use insecure connection to gRPC server")
	grpcUserOpt = flag.String("grpc-user", "", "username for the basic authentication of the remote gRPC server")
	grpcPassFileOpt = flag.String("grpc-password-file", "",
		"file containing the password for the basic authentication of the remote gRPC server")
	snapshotURLOpt = flag.String(
		"snapshot-url",
		cmd.DefaultSnapshotURL,
		"snapshot server base URL (network path will be appended automatically)",
	)
	version.NodeAgent.AppType = "gui"
}

// initGTK initializes GTK. It needs a display, so it is called when the GUI starts
// and not when the package is loaded: the tests of the package run without a display.
func initGTK() {
	// The GUI text is English only, so force GTK's built-in widget strings
	// (assistant navigation, about dialog, file choosers ...) to English too
	// instead of following the operating system locale. This affects message
	// translation only, not number or date formatting.
	// Use g_setenv, not os.Setenv: on Windows the latter uses the Win32
	// environment block, which the C runtime that GTK/libintl links against
	// does not read, so the change would be invisible to gettext.
	glib.Setenv("LANGUAGE", "en", true)
	glib.Setenv("LC_MESSAGES", "en", true)
	forceEnglishUILanguage()

	gtk.Init()
	// Initialize libadwaita so its StyleManager controls the light/dark
	// appearance, independent of the operating system theme.
	adw.Init()
	// Apply the user's theme before any window is created.
	gtkutil.ApplySavedTheme()
}

//nolint:gocognit // needs refactoring
func main() {
	initGTK()
	flag.Parse()

	gtkutil.Logf("Starting Pactus GUI")
	// The gtk should run on main thread.
	runtime.UnlockOSThread()
	runtime.LockOSThread()
	gtkutil.Logf("Locking main thread")

	// Create a new app.
	app := gtk.NewApplication(appID, gio.ApplicationNonUnique)
	gtk.WidgetSetDefaultDirection(gtk.TextDirLTR)

	// apply custom css
	provider := gtk.NewCSSProvider()
	provider.LoadFromString(assets.MainWindowCSS)
	display := gdk.DisplayGetDefault()

	gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)

	assets.InitAssets()

	workingDir, err := filepath.Abs(*workingDirOpt)
	if err != nil {
		gtkutil.Logf("Aborted! %v", err)

		return
	}

	// The autostart entry starts the GUI with the same flags, except the wallet
	// password, which must not be stored. The working directory is always set
	// as an absolute path, because the login session may start elsewhere.
	autostartArgs := autostart.ArgsFromFlags(flag.CommandLine, "password", "working-dir")
	autostartEntry := autostart.New(append(autostartArgs, "-working-dir="+workingDir))

	ctx, cancel := context.WithCancel(context.Background())

	var fileLock *flock.Flock
	if *grpcAddrOpt == "" {
		// If node is not initialized yet
		if util.IsDirNotExistsOrEmpty(workingDir) {
			network := genesis.Mainnet
			if *testnetOpt {
				network = genesis.Testnet
			}
			if !startupAssistant(ctx, workingDir, network, *snapshotURLOpt) {
				return
			}
		}
	}
	var guiNode *node.Node

	// Connect function to application startup event, this is not required.
	app.Connect("startup", func() {
		gtkutil.Logf("application startup")
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
		gtkutil.Logf("Application shutdown")
		shutdown()
	})

	// Connect function to application activate event
	app.Connect("activate", func() {
		activateOnce.Do(func() {
			gtkutil.Logf("application activate")

			splash := view.NewSplashWindow(app)
			splash.SetVersion(version.NodeVersion().StringWithAlias())
			gtk.WindowSetAutoStartupNotification(false)
			splash.ShowAll()
			gtk.WindowSetAutoStartupNotification(true)
			app.AddWindow(splash.Window())

			notify := func(msg string) {
				gtkutil.Logf("Splash msg: %s", msg)

				gtkutil.IdleAddAsync(func() {
					splash.UpdateStatus(msg)
				})
			}

			go func() {
				var grpcAddr string
				var grpcInsecure bool

				if *grpcAddrOpt == "" {
					// Define the lock file path
					lockFilePath := filepath.Join(workingDir, ".pactus.lock")
					fileLock = flock.New(lockFilePath)

					locked, err := fileLock.TryLock()
					if err != nil {
						gtkutil.ShowErrorDialog(splash.Window(),
							fmt.Sprintf("Aborted! Unable to acquire file lock. %v", err),
							func(_ int) {
								splash.Destroy()
							})

						return
					}

					if !locked {
						gtkutil.ShowErrorDialog(splash.Window(),
							fmt.Sprintf("Could not lock '%s', another instance is running?", lockFilePath),
							func(_ int) {
								splash.Destroy()
							})

						return
					}

					reportStatus(notify, "Starting local node...")
					time.Sleep(1 * time.Second)

					guiNode, err = newNode(ctx, workingDir, notify)
					if err != nil {
						gtkutil.ShowWarningDialog(splash.Window(), err.Error(),
							func(_ int) {
								splash.Destroy()
							})

						return
					}

					// The GUI talks to its own gRPC server, on the loopback interface.
					grpcAddr, err = guiNode.StartGUIGRPC()
					if err != nil {
						gtkutil.ShowErrorDialog(splash.Window(), err.Error(),
							func(_ int) {
								splash.Destroy()
							})

						return
					}
					grpcInsecure = true
				} else {
					reportStatus(notify, "Connecting to remote node...")
					time.Sleep(1 * time.Second)

					grpcAddr = *grpcAddrOpt
					grpcInsecure = *grpcInsecureOpt
				}

				grpcConn, err = openGRPCConn(grpcAddr, grpcInsecure, *grpcAddrOpt != "")
				if err != nil {
					gtkutil.ShowErrorDialog(splash.Window(), err.Error(),
						func(_ int) {
							splash.Destroy()
						})

					return
				}

				var connectionLabel, connectionValue string
				var isLocal bool
				if *grpcAddrOpt == "" {
					isLocal = true
					connectionLabel = "Working Directory"
					connectionValue = workingDir
				} else {
					isLocal = false
					connectionLabel = "Remote Address"
					connectionValue = *grpcAddrOpt
				}
				gui, err = gtkapp.Run(ctx, grpcConn, app, notify,
					connectionLabel, connectionValue, workingDir, isLocal, autostartEntry)
				if err != nil {
					gtkutil.ShowErrorDialog(splash.Window(), err.Error(), func(_ int) {
						splash.Destroy()
					})

					return
				}

				gtkutil.IdleAddSync(func() {
					splash.Destroy()
				})
			}()
		})
	})

	signal.HandleInterrupt(func() {
		gtkutil.Logf("Exiting...")
		gtkutil.IdleAddSync(func() {
			shutdown()
		})
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

// openGRPCConn connects to the gRPC server of the node. The local node has its own gRPC
// server for the GUI, without authentication. A remote server may require it.
func openGRPCConn(addr string, insecureCredentials, remote bool) (*grpc.ClientConn, error) {
	var auth *basicauth.BasicAuth
	if remote {
		var err error
		auth, err = remoteBasicAuth(*grpcUserOpt, *grpcPassFileOpt)
		if err != nil {
			return nil, err
		}
	}

	return newRemoteGRPCConn(addr, insecureCredentials, auth)
}

// remoteBasicAuth returns the basic authentication for the remote gRPC server, or nil if no user is set.
// The password is read from a file: a password on the command line is visible to the other users.
func remoteBasicAuth(user, passwordFile string) (*basicauth.BasicAuth, error) {
	if user == "" {
		return nil, nil //nolint:nilnil // no authentication
	}
	if passwordFile == "" {
		return nil, errors.New("-grpc-user needs -grpc-password-file")
	}

	data, err := util.ReadFile(passwordFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read the gRPC password file: %w", err)
	}

	return basicauth.New(user, strings.TrimRight(string(data), "\r\n")), nil
}

// newRemoteGRPCConn creates a gRPC client for a remote URL.
func newRemoteGRPCConn(addr string, insecureCredentials bool, auth *basicauth.BasicAuth) (*grpc.ClientConn, error) {
	target, prefix, err := util.ParseGRPCAddr(addr, insecureCredentials)
	if err != nil {
		return nil, err
	}
	var creds credentials.TransportCredentials
	if insecureCredentials {
		gtkutil.Logf("Connecting to node using insecure credentials. target: %s, prefix: %s", target, prefix)
		creds = insecure.NewCredentials()
	} else {
		gtkutil.Logf("Connecting to node using TLS. target: %s, prefix: %s", target, prefix)
		creds = credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(methodPrefixUnaryInterceptor(prefix)),
		grpc.WithStreamInterceptor(methodPrefixStreamInterceptor(prefix)),
	}
	if auth != nil {
		if insecureCredentials {
			gtkutil.Logf("The password of the gRPC server is sent without encryption")
		}
		opts = append(opts, grpc.WithPerRPCCredentials(auth))
	}

	return grpc.NewClient(target, opts...)
}

func methodPrefixUnaryInterceptor(prefix string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any,
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {
		return invoker(ctx, prefix+method, req, reply, cc, opts...)
	}
}

func methodPrefixStreamInterceptor(prefix string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc,
		cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		return streamer(ctx, desc, cc, prefix+method, opts...)
	}
}

func newNode(ctx context.Context, workingDir string, statusCb statusReporter) (*node.Node, error) {
	// change working directory
	if err := os.Chdir(workingDir); err != nil {
		gtkutil.Logf("Aborted! Unable to changes working directory. %v", err)

		return nil, err
	}

	reportStatus(statusCb, "Opening wallet...")
	passwordFetcher := func() (string, bool) {
		gtkutil.Logf("Fetching wallet password")

		if *passwordOpt != "" {
			return *passwordOpt, true
		}

		password := ""
		ok := false
		resultChan := make(chan bool)

		gtkutil.IdleAddSync(func() {
			controller.PromptWalletPassword(func(_password string, _ok bool) {
				password = _password
				ok = _ok

				resultChan <- true
			})
		})

		<-resultChan

		return password, ok
	}

	reportStatus(statusCb, "Starting node services...")
	n, err := cmd.StartNode(ctx, workingDir, passwordFetcher, guiConfig)
	if err != nil {
		return nil, err
	}

	return n, nil
}

// guiConfig is the config of the node that the GUI runs. The gRPC server of the config is left
// as the operator set it: the GUI starts its own gRPC server, see Node.StartGUIGRPC.
func guiConfig(cfg *config.Config) *config.Config {
	return cfg
}
