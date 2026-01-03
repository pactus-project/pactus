package main

import (
	"context"
	"path/filepath"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildInitCmd builds a sub-command to initialize the Pactus blockchain node.
func buildInitCmd(parentCmd *cobra.Command) {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "initialize the Pactus blockchain node",
	}
	parentCmd.AddCommand(initCmd)

	workingDirOpt := addWorkingDirOption(initCmd)

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
			terminal.PrintErrorMsgf("The working directory is not empty: %s", workingDir)

			return
		}

		var mnemonic string
		if *restoreOpt == "" {
			mnemonic, _ = wallet.GenerateMnemonic(*entropyOpt)

			terminal.PrintLine()
			terminal.PrintInfoMsgf("🌱 Your wallet seed phrase:")
			terminal.PrintInfoMsgBoldf("   %s", mnemonic)
			terminal.PrintLine()
			terminal.PrintWarnMsgf("⚠️  CRITICAL: Write down this seed phrase and store it safely!")
			terminal.PrintWarnMsgf("   This is the ONLY way to recover your wallet if needed.")
			terminal.PrintWarnMsgf("   Never share it with anyone or store it electronically.")
			terminal.PrintLine()
			confirmed := prompt.PromptConfirm("Have you written down the seed phrase? Continue with initialization")
			if !confirmed {
				return
			}
		} else {
			mnemonic = *restoreOpt
			err := wallet.CheckMnemonic(*restoreOpt)
			terminal.FatalErrorCheck(err)
		}

		var password string
		if *passwordOpt == "" {
			terminal.PrintLine()
			terminal.PrintInfoMsgf("🔐 Set a password for your wallet")
			terminal.PrintInfoMsgf("   This password will be required to access your wallet")
			password = prompt.PromptPassword("Wallet Password", true)
		} else {
			password = *passwordOpt
		}

		var valNum int
		if *valNumOpt == 0 {
			terminal.PrintLine()
			terminal.PrintInfoMsgBoldf("🏛️  How many validators do you want to create?")
			terminal.PrintInfoMsgf("   • Each node can run up to 32 validators")
			terminal.PrintInfoMsgf("   • Each validator can stake up to 1,000 coins")
			terminal.PrintInfoMsgf("   • Choose based on your total stake amount")
			terminal.PrintLine()
			valNum = prompt.PromptInputWithRange("Number of Validators", 7, 1, 32)
		} else {
			if *valNumOpt < 1 || *valNumOpt > 32 {
				terminal.PrintErrorMsgf("%v is not in valid range of validator number, it should be between 1 and 32", *valNumOpt)

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

		wlt, rewardAddrs, err := cmd.CreateNode(context.Background(),
			valNum, chain, workingDir, mnemonic, password)
		terminal.FatalErrorCheck(err)

		// Recovering addresses
		if *restoreOpt != "" {
			cmd.RecoverWalletAddresses(wlt, password)
		}

		terminal.PrintLine()
		terminal.PrintInfoMsgBoldf("🏛️  Validator Addresses:")
		for i, addrInfo := range wlt.ListAddresses(wallet.OnlyValidatorAddresses()) {
			terminal.PrintInfoMsgf("   %d. %s", i+1, addrInfo.Address)
		}
		terminal.PrintLine()

		terminal.PrintInfoMsgBoldf("💰 Reward Address:")
		terminal.PrintInfoMsgf("   %s", rewardAddrs)
		terminal.PrintLine()

		terminal.PrintInfoMsgf("🌐 Network: %v", chain.String())
		terminal.PrintInfoMsgf("📁 Working Directory: %v", workingDir)
		terminal.PrintLine()
		terminal.PrintSuccessMsgf("✅ Pactus node successfully initialized!")
		terminal.PrintLine()
		terminal.PrintInfoMsgf("🚀 To start your node, run:")
		terminal.PrintInfoMsgBoldf("   %s start -w %s", cmd.PactusDaemonName(), workingDir)
	}
}
