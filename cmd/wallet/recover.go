package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildRecoverCmd builds a command to recover a wallet using a mnemonic (seed phrase).
func buildRecoverCmd(parentCmd *cobra.Command) {
	recoverCmd := &cobra.Command{
		Use:   "recover",
		Short: "Recover waller from the seed phrase (mnemonic)",
	}
	parentCmd.AddCommand(recoverCmd)

	passOpt := addPasswordOption(recoverCmd)
	testnetOpt := recoverCmd.Flags().Bool("testnet", true,
		"Recover the wallet for the testnet environment")
	seedOpt := recoverCmd.Flags().StringP("seed", "s", "", "Mnemonic (seed phrase) used for wallet recovery")

	recoverCmd.Run = func(_ *cobra.Command, _ []string) {
		mnemonic := *seedOpt
		if mnemonic == "" {
			mnemonic = cmd.PromptInput("Seed")
		}
		chainType := genesis.Mainnet
		if *testnetOpt {
			chainType = genesis.Testnet
		}
		wallet, err := wallet.Create(*pathOpt, mnemonic, *passOpt, chainType)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("Wallet successfully recovered and saved at: %s", wallet.Path())
	}
}

// buildGetSeedCmd builds a command to display the wallet's mnemonic (seed phrase).
func buildGetSeedCmd(parentCmd *cobra.Command) {
	getSeedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Display the mnemonic (seed phrase) that can be used to recover this wallet",
	}
	parentCmd.AddCommand(getSeedCmd)

	passOpt := addPasswordOption(getSeedCmd)

	getSeedCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		password := getPassword(wallet, *passOpt)
		mnemonic, err := wallet.Mnemonic(password)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("Your wallet's seed phrase is: \"%v\"", mnemonic)
	}
}
