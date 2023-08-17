package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// Generate creates a new wallet.
func buildGenerateCmd(parentCmd *cobra.Command) {
	generateCmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new wallet",
	}
	parentCmd.AddCommand(generateCmd)

	testnetOpt := generateCmd.Flags().Bool("testnet", false, "creating wallet for testnet")
	entropyOpt := generateCmd.Flags().Int("entropy", 128, "Entropy bit length")

	generateCmd.Run = func(_ *cobra.Command, _ []string) {
		password := cmd.PromptPassword("Password", true)
		mnemonic := wallet.GenerateMnemonic(*entropyOpt)

		network := genesis.Mainnet
		if *testnetOpt {
			network = genesis.Testnet
		}
		wallet, err := wallet.Create(*pathArg, mnemonic, password, network)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsg("Wallet created successfully at: %s", wallet.Path())
		cmd.PrintInfoMsg("Seed: \"%v\"", mnemonic)
		cmd.PrintWarnMsg("Please keep your seed in a safe place; " +
			"if you lose it, you will not be able to restore your wallet.")
	}
}

// ChangePassword updates the wallet password.
func buildChangePasswordCmd(parentCmd *cobra.Command) {
	changePasswordCmd := &cobra.Command{
		Use:   "password",
		Short: "Change wallet password",
	}
	parentCmd.AddCommand(changePasswordCmd)
	passOpt := addPasswordOption(parentCmd)

	changePasswordCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		oldPassword := getPassword(wallet, *passOpt)
		newPassword := cmd.PromptPassword("New Password", true)

		err = wallet.UpdatePassword(oldPassword, newPassword)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintWarnMsg("Wallet password updated")
	}
}
