package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildRecoverCmd builds a command to recover a wallet using a mnemonic (seed phrase).
func buildRecoverCmd(parentCmd *cobra.Command) {
	recoverCmd := &cobra.Command{
		Use:   "recover",
		Short: "recover wallet from the seed phrase or mnemonic",
	}
	parentCmd.AddCommand(recoverCmd)

	passOpt := addPasswordOption(recoverCmd)
	testnetOpt := recoverCmd.Flags().Bool("testnet", false,
		"recover the wallet for the testnet environment")
	seedOpt := recoverCmd.Flags().StringP("seed", "s", "", "mnemonic or seed phrase used for wallet recovery")

	recoverCmd.Run = func(_ *cobra.Command, _ []string) {
		mnemonic := *seedOpt
		if mnemonic == "" {
			mnemonic = prompt.PromptInput("Seed")
		}
		chainType := genesis.Mainnet
		if *testnetOpt {
			chainType = genesis.Testnet
		}
		wlt, err := wallet.Create(*pathOpt, mnemonic, *passOpt, chainType)
		terminal.FatalErrorCheck(err)

		ctx, cancel := context.WithCancel(context.Background())
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			cancel()
		}()

		terminal.PrintInfoMsgf("🔄 Recovering wallet addresses...")
		terminal.PrintInfoMsgf("   Press Ctrl+C to abort if needed")
		terminal.PrintLine()

		index := 0
		err = wlt.RecoveryAddresses(ctx, *passOpt, func(addr string) {
			terminal.PrintInfoMsgf("%d. %s", index+1, addr)
			index++
		})

		// Check if context was cancelled
		wasInterrupted := ctx.Err() != nil

		if err != nil {
			if wasInterrupted || errors.Is(err, context.Canceled) {
				terminal.PrintLine()
				terminal.PrintWarnMsgf("⚠️  Recovery aborted by user")
			} else {
				terminal.PrintLine()
				terminal.PrintWarnMsgf("Recovery addresses failed: %v", err)
			}
		}

		// Always save the wallet before exiting
		terminal.PrintLine()
		terminal.PrintInfoMsgf("💾 Saving wallet...")
		err = wlt.Save()
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintSuccessMsgf("✅ Wallet successfully recovered and saved at: %s", wlt.Path())
	}
}

// buildGetSeedCmd builds a command to display the wallet's mnemonic (seed phrase).
func buildGetSeedCmd(parentCmd *cobra.Command) {
	getSeedCmd := &cobra.Command{
		Use:   "seed",
		Short: "displays the seed phrase that can be used to recover this wallet",
	}
	parentCmd.AddCommand(getSeedCmd)

	passOpt := addPasswordOption(getSeedCmd)

	getSeedCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet()
		terminal.FatalErrorCheck(err)

		password := getPassword(wlt, *passOpt)
		mnemonic, err := wlt.Mnemonic(password)
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintInfoMsgf("🌱 Your wallet seed phrase:")
		terminal.PrintInfoMsgBoldf("   %v", mnemonic)
		terminal.PrintLine()
		terminal.PrintWarnMsgf("⚠️  CRITICAL: Write down this seed phrase and store it safely!")
		terminal.PrintWarnMsgf("   This is the ONLY way to recover your wallet if needed.")
		terminal.PrintWarnMsgf("   Never share it with anyone or store it electronically.")
	}
}
