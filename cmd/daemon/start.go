package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/ezex-io/gopkg/signal"
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

	zmqBlockInfoOpt := startCmd.Flags().String("zmq-block-info", "",
		"enable zeromq block info publisher, example: tcp://127.0.0.1:28332")

	zmqTxInfoOpt := startCmd.Flags().String("zmq-tx-info", "",
		"enable zeromq transaction info publisher, example: tcp://127.0.0.1:28332")

	zmqRawBlockOpt := startCmd.Flags().String("zmq-raw-block", "",
		"enable zeromq raw block publisher, example: tcp://127.0.0.1:28332")

	zmqRawTxOpt := startCmd.Flags().String("zmq-raw-tx", "",
		"enable zeromq raw transaction publisher, example: tcp://127.0.0.1:28332")

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

			if *zmqBlockInfoOpt != "" {
				cfg.ZeroMq.ZmqPubBlockInfo = *zmqBlockInfoOpt
			}

			if *zmqTxInfoOpt != "" {
				cfg.ZeroMq.ZmqPubTxInfo = *zmqTxInfoOpt
			}

			if *zmqRawBlockOpt != "" {
				cfg.ZeroMq.ZmqPubRawBlock = *zmqRawBlockOpt
			}

			if *zmqRawTxOpt != "" {
				cfg.ZeroMq.ZmqPubRawTx = *zmqRawTxOpt
			}

			return cfg
		}

		ctx, cancel := context.WithCancel(context.Background())
		node, err := cmd.StartNode(ctx, workingDir, passwordFetcher, configModifier)
		terminal.FatalErrorCheck(err)

		signal.HandleInterrupt(func() {
			terminal.PrintInfoMsgf("Exiting...")
			cancel()

			_ = fileLock.Unlock()
			node.Stop()
		})

		// run forever (the node will not be returned)
		select {}
	}
}
