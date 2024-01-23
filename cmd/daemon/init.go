package main

import (
	"path/filepath"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildInitCmd builds a sub-command to initialize the Pactus blockchain node.
func buildInitCmd(parentCmd *cobra.Command) {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "initialize the Pactus Blockchain node",
	}
	parentCmd.AddCommand(initCmd)
	workingDirOpt := initCmd.Flags().StringP("working-dir", "w", cmd.PactusDefaultHomeDir(),
		"a path to the working directory to save the wallet and node files")

	testnetOpt := initCmd.Flags().Bool("testnet", false,
		"initialize working directory for joining the testnet")

	localnetOpt := initCmd.Flags().Bool("localnet", false,
		"initialize working directory for localnet (for development)")

	restoreOpt := initCmd.Flags().String("restore", "",
		"restore the 'default_wallet' using a mnemonic or seed phrase")

	passwordOpt := initCmd.Flags().StringP("password", "p", "",
		"the wallet password")

	entropyOpt := initCmd.Flags().IntP("entropy", "e", 128,
		"entropy bits for seed generation. range: 128 to 256")

	valNumOpt := initCmd.Flags().IntP("val-num", "", 0,
		"number of validators to be created. range: 1 to 32")

	initCmd.Run = func(_ *cobra.Command, _ []string) {
		workingDir, _ := filepath.Abs(*workingDirOpt)
		if !util.IsDirNotExistsOrEmpty(workingDir) {
			cmd.PrintErrorMsgf("The working directory is not empty: %s", workingDir)

			return
		}
		var mnemonic string
		if *restoreOpt == "" {
			mnemonic, _ = wallet.GenerateMnemonic(*entropyOpt)
			cmd.PrintLine()
			cmd.PrintInfoMsgf("Your wallet seed is:")
			cmd.PrintInfoMsgBoldf("   " + mnemonic)
			cmd.PrintLine()
			cmd.PrintWarnMsgf("Write down this seed on a piece of paper to recover your validator key in the future.")
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

		var password string
		if *passwordOpt == "" {
			cmd.PrintLine()
			cmd.PrintInfoMsgf("Enter a password for wallet")
			password = cmd.PromptPassword("Password", true)
		} else {
			password = *passwordOpt
		}

		var valNum int
		if *valNumOpt == 0 {
			cmd.PrintLine()
			cmd.PrintInfoMsgBoldf("How many validators do you want to create?")
			cmd.PrintInfoMsgf("Each node can run up to 32 validators, and each validator can hold up to 1000 staked coins.")
			cmd.PrintInfoMsgf("You can define validators based on the amount of coins you want to stake.")
			valNum = cmd.PromptInputWithRange("Number of Validators", 7, 1, 32)
		} else {
			if *valNumOpt < 1 || *valNumOpt > 32 {
				cmd.PrintErrorMsgf("%v is not in valid range of validator number, it should be between 1 and 32", *valNumOpt)

				return
			}
			valNum = *valNumOpt
		}

		chain := genesis.Mainnet
		// The order of checking the network (chain type) matters here.
		if *testnetOpt {
			chain = genesis.Testnet
		}
		if *localnetOpt {
			chain = genesis.Localnet
		}
		validatorAddrs, rewardAddrs, err := cmd.CreateNode(valNum, chain, workingDir, mnemonic, password)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgBoldf("Validator addresses:")
		for i, addr := range validatorAddrs {
			cmd.PrintInfoMsgf("%v- %s", i+1, addr)
		}
		cmd.PrintLine()

		cmd.PrintInfoMsgBoldf("Reward addresses:")
		for i, addr := range rewardAddrs {
			cmd.PrintInfoMsgf("%v- %s", i+1, addr)
		}

		cmd.PrintLine()
		cmd.PrintInfoMsgBoldf("Network: %v", chain.String())
		cmd.PrintLine()
		cmd.PrintSuccessMsgf("A pactus node is successfully initialized at %v", workingDir)
		cmd.PrintLine()
		cmd.PrintInfoMsgf("You can start the node by running this command:")
		cmd.PrintInfoMsgf("./pactus-daemon start -w %v", workingDir)
	}
}
