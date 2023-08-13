package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// Recover recovers a wallet from mnemonic (seed phrase).
func buildRecoverCmd(parentCmd *cobra.Command) {
	recoverCmd := &cobra.Command{
		Use:   "recover",
		Short: "Recover waller from the seed phrase (mnemonic)",
	}
	parentCmd.AddCommand(recoverCmd)

	passwordOpt := addPasswordOption(parentCmd)

	seedOpt := recoverCmd.Flags().StringP("seed", "s", "", "wallet seed phrase (mnemonic)")

	recoverCmd.Run = func(_ *cobra.Command, _ []string) {
		mnemonic := *seedOpt
		if mnemonic == "" {
			mnemonic = cmd.PromptInput("Seed")
		}
		password := *passwordOpt

		wallet, err := wallet.Create(*pathArg, mnemonic, password, 0)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsg("Wallet recovered successfully at: %s", wallet.Path())
	}
}

// GetSeed prints the seed phrase (mnemonics).
func buildGetSeedCmd(parentCmd *cobra.Command) {
	getSeedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Show secret seed phrase (mnemonic) that can be used to recover this wallet",
	}
	parentCmd.AddCommand(getSeedCmd)

	passOpt := addPasswordOption(parentCmd)

	getSeedCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		password := getPassword(wallet, *passOpt)
		mnemonic, err := wallet.Mnemonic(password)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
	}
}
