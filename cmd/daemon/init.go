package main

import (
	"path/filepath"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// BuildInitCmd builds the init command for the Pactus blockchain.
func buildInitCmd(parentCmd *cobra.Command) {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize the Pactus blockchain",
	}
	parentCmd.AddCommand(initCmd)
	workingDirOpt := initCmd.Flags().StringP("working-dir", "w",
		cmd.PactusHomeDir(), "A path to the working directory to save the wallet and node files")

	testnetOpt := initCmd.Flags().Bool("testnet", true,
		"Initialize working directory for joining the testnet") // TODO: make it false after mainnet launch

	localnetOpt := initCmd.Flags().Bool("localnet", false,
		"Initialize working directory for localnet (for developers)")

	restoreOpt := initCmd.Flags().String("restore", "", "Restore the default_wallet using a mnemonic (seed phrase)")

	initCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		if !util.IsDirNotExistsOrEmpty(workingDir) {
			cmd.PrintErrorMsg("The working directory is not empty: %s", workingDir)
			return
		}
		mnemonic := ""
		if len(*restoreOpt) == 0 {
			mnemonic = wallet.GenerateMnemonic(128)
			cmd.PrintLine()
			cmd.PrintInfoMsg("Your wallet seed is:")
			cmd.PrintInfoMsgBold("   " + mnemonic)
			cmd.PrintLine()
			cmd.PrintWarnMsg("Write down this seed on a piece of paper to recover your validator key in future.")
			cmd.PrintLine()
			confirmed := cmd.PromptConfirm("Do you want to continue")
			if !confirmed {
				return
			}
		} else {
			mnemonic = *restoreOpt
			err := wallet.CheckMnemonic(*restoreOpt)
			cmd.FatalErrorCheck(err)
		}
		cmd.PrintLine()
		cmd.PrintInfoMsg("Enter a password for wallet")
		password := cmd.PromptPassword("Password", true)

		cmd.PrintLine()
		cmd.PrintInfoMsgBold("How many validators do you want to create?")
		cmd.PrintInfoMsg("Each node can run up to 32 validators, and each validator can hold up to 1000 staked coins.")
		cmd.PrintInfoMsg("You can define validators based on the amount of coins you want to stake.")
		numValidators := cmd.PromptInputWithRange("Number of Validators", 7, 1, 32)

		chain := genesis.Mainnet
		// The order of checking the network (chain type) matters here.
		if *testnetOpt {
			chain = genesis.Testnet
		}
		if *localnetOpt {
			chain = genesis.Localnet
		}
		validatorAddrs, rewardAddrs, err := cmd.CreateNode(numValidators, chain, workingDir, mnemonic, password)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgBold("Validator addresses:")
		for i, addr := range validatorAddrs {
			cmd.PrintInfoMsg("%v- %s", i+1, addr)
		}
		cmd.PrintLine()

		cmd.PrintInfoMsgBold("Reward addresses:")
		for i, addr := range rewardAddrs {
			cmd.PrintInfoMsg("%v- %s", i+1, addr)
		}

		cmd.PrintLine()
		cmd.PrintInfoMsgBold("Network: %v", chain.String())
		cmd.PrintLine()
		cmd.PrintSuccessMsg("A pactus node is successfully initialized at %v", workingDir)
		cmd.PrintLine()
		cmd.PrintInfoMsg("You can start the node by running this command:")
		cmd.PrintInfoMsg("./pactus-daemon start -w %v", workingDir)
	}
}
