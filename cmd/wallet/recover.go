package main

import (
	"context"

	"github.com/pactus-project/pactus/cmd"
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
		wlt, err := wallet.Create(context.Background(), *pathOpt, mnemonic, *passOpt, chainType)
		terminal.FatalErrorCheck(err)

		cmd.RecoverWalletAddresses(wlt, *passOpt)

		// Always save the wallet before exiting
		terminal.PrintLine()
		terminal.PrintInfoMsgf("💾 Saving wallet...")

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
		wlt, err := openWallet(context.Background())
		terminal.FatalErrorCheck(err)
		defer wlt.Close()

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
