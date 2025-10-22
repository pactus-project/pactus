package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildCreateCmd builds a command to create a new wallet.
func buildCreateCmd(parentCmd *cobra.Command) {
	generateCmd := &cobra.Command{
		Use:   "create",
		Short: "creating a new wallet",
	}
	parentCmd.AddCommand(generateCmd)

	testnetOpt := generateCmd.Flags().Bool("testnet", false,
		"create a wallet for the testnet environment")
	entropyOpt := generateCmd.Flags().Int("entropy", 128,
		"specify the entropy bit length")

	generateCmd.Run = func(_ *cobra.Command, _ []string) {
		password := prompt.PromptPassword("Password", true)
		mnemonic, err := wallet.GenerateMnemonic(*entropyOpt)
		cmd.FatalErrorCheck(err)

		network := genesis.Mainnet
		if *testnetOpt {
			network = genesis.Testnet
		}
		wlt, err := wallet.Create(*pathOpt, mnemonic, password, network)
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsgf("Your wallet was successfully created at: %s", wlt.Path())
		cmd.PrintInfoMsgf("Seed phrase: \"%v\"", mnemonic)
		cmd.PrintWarnMsgf("Please keep your seed in a safe place; " +
			"if you lose it, you will not be able to restore your wallet.")
	}
}

// buildChangePasswordCmd builds a command to update the wallet's password.
func buildChangePasswordCmd(parentCmd *cobra.Command) {
	changePasswordCmd := &cobra.Command{
		Use:   "password",
		Short: "changes the wallet's password",
	}
	parentCmd.AddCommand(changePasswordCmd)
	passOpt := addPasswordOption(changePasswordCmd)

	changePasswordCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		oldPassword := getPassword(wlt, *passOpt)
		newPassword := prompt.PromptPassword("New Password", true)

		err = wlt.UpdatePassword(oldPassword, newPassword)
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintWarnMsgf("Your wallet password successfully updated.")
	}
}
