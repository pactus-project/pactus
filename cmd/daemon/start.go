package main

import (
	"context"
	"os"
	ossignal "os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/gofrs/flock"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/config"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/spf13/cobra"
)

// buildStartCmd builds a sub-command to starts the Pactus blockchain node.
func buildStartCmd(parentCmd *cobra.Command) {
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "start the Pactus blockchain node",
	}

	parentCmd.AddCommand(startCmd)

	workingDirOpt := addWorkingDirOption(startCmd)

	passwordOpt := startCmd.Flags().StringP("password", "p", "",
		"the wallet password")

	passwordFromFileOpt := startCmd.Flags().String("password-from-file", "",
		"the file containing the wallet password")

	gRPCOpt := startCmd.Flags().String("grpc", "",
		"enable gRPC transport, example: localhost:50051")

	gRPCWalletOpt := startCmd.Flags().Bool("grpc-wallet", false, "enable gRPC wallet service")

	startCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		// change working directory
		err := os.Chdir(workingDir)
		terminal.FatalErrorCheck(err)

		// Define the lock file path
		lockFilePath := filepath.Join(workingDir, ".pactus.lock")
		fileLock := flock.New(lockFilePath)

		locked, err := fileLock.TryLock()
		terminal.FatalErrorCheck(err)

		if !locked {
			terminal.PrintWarnMsgf("Could not lock '%s', another instance is running?", lockFilePath)

			return
		}

		passwordFetcher := func() (string, bool) {
			var password string

			if *passwordOpt != "" {
				password = *passwordOpt
			} else if *passwordFromFileOpt != "" {
				b, err := util.ReadFile(*passwordFromFileOpt)
				terminal.FatalErrorCheck(err)

				password = strings.TrimSpace(string(b))
			} else {
				password = prompt.PromptPassword("Wallet password", false)
			}

			return password, true
		}

		configModifier := func(cfg *config.Config) *config.Config {
			if *gRPCOpt != "" {
				cfg.GRPC.Enable = true
				cfg.GRPC.EnableWallet = *gRPCWalletOpt
				cfg.GRPC.Listen = *gRPCOpt
			}

			return cfg
		}

		ctx, cancel := context.WithCancel(context.Background())
		shutdown := newShutdown(cancel, func() { _ = fileLock.Unlock() }, os.Exit)

		startedNode, err := cmd.StartNode(ctx, workingDir, passwordFetcher, configModifier)
		if err != nil {
			_ = fileLock.Unlock()
		}
		terminal.FatalErrorCheck(err)

		shutdown.setNode(startedNode)

		// run forever (the node will not be returned)
		select {}
	}
}

// shutdown stops the node on SIGINT or SIGTERM.
//
// The signals are handled from the start: a signal while the node starts cancels the start,
// and the node is stopped once it is started, so the store is closed. A second signal exits
// at once, in case stopping the node blocks.
type shutdown struct {
	lk       sync.Mutex
	node     stopper
	signaled os.Signal
	release  func()
	exit     func(code int)
}

// stopper is the part of the node that the shutdown uses.
type stopper interface {
	Stop()
}

func newShutdown(cancel context.CancelFunc, release func(), exit func(code int)) *shutdown {
	shut := &shutdown{release: release, exit: exit}

	sigChan := make(chan os.Signal, 2)
	ossignal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		terminal.PrintInfoMsgf("Exiting...")
		cancel()

		go func() {
			<-sigChan
			terminal.PrintWarnMsgf("Forced exit")
			exit(exitCode(sig))
		}()

		shut.lk.Lock()
		shut.signaled = sig
		n := shut.node
		shut.lk.Unlock()

		if n != nil {
			shut.stop(n, sig)
		}
		// Otherwise the node is starting: setNode stops it.
	}()

	return shut
}

// setNode records the started node, and stops it if a signal arrived while it was starting.
func (shut *shutdown) setNode(n stopper) {
	shut.lk.Lock()
	shut.node = n
	sig := shut.signaled
	shut.lk.Unlock()

	if sig != nil {
		shut.stop(n, sig)
	}
}

func (shut *shutdown) stop(n stopper, sig os.Signal) {
	n.Stop()
	shut.release()
	shut.exit(exitCode(sig))
}

// exitCode follows the Unix convention: 128 + the signal number.
func exitCode(sig os.Signal) int {
	if s, ok := sig.(syscall.Signal); ok {
		return 128 + int(s)
	}

	return 1
}
