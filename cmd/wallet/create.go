package main

import (
	"context"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildCreateCmd builds a command to create a new wallet.
func buildCreateCmd(parentCmd *cobra.Command) {
	createCmd := &cobra.Command{
		Use:   "create",
		Short: "create a new wallet",
	}
	parentCmd.AddCommand(createCmd)

	testnetOpt := createCmd.Flags().Bool("testnet", false,
		"create a wallet for the testnet environment")
	entropyOpt := createCmd.Flags().Int("entropy", 128,
		"specify the entropy bit length")

	createCmd.Run = func(_ *cobra.Command, _ []string) {
		password := prompt.PromptPassword("Password", true)
		mnemonic, err := wallet.GenerateMnemonic(*entropyOpt)
		terminal.FatalErrorCheck(err)

		network := genesis.Mainnet
		if *testnetOpt {
			network = genesis.Testnet
		}
		wlt, err := wallet.Create(context.Background(), *pathOpt, mnemonic, password, network)
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintSuccessMsgf("✅ Wallet successfully created at: %s", wlt.Path())
		terminal.PrintLine()
		terminal.PrintInfoMsgf("🌱 Your wallet seed phrase:")
		terminal.PrintInfoMsgBoldf("   %v", mnemonic)
		terminal.PrintLine()
		terminal.PrintWarnMsgf("⚠️  CRITICAL: Write down this seed phrase and store it safely!")
		terminal.PrintWarnMsgf("   This is the ONLY way to recover your wallet if needed.")
		terminal.PrintWarnMsgf("   Never share it with anyone or store it electronically.")
	}
}
