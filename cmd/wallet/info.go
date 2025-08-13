package main

import (
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/spf13/cobra"
)

// buildInfoCmd builds all sub-commands related to the wallet information.
func buildInfoCmd(parentCmd *cobra.Command) {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "retrieving the wallet information.",
	}

	parentCmd.AddCommand(infoCmd)

	infoCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		info := wlt.Info()

		cmd.PrintInfoMsgf("Version: %d", info.Version)
		cmd.PrintInfoMsgf("UUID: %s", info.UUID)
		cmd.PrintInfoMsgf("Default fee: %s", info.DefaultFee.String())
		cmd.PrintInfoMsgf("Created at: %s", info.CreatedAt.Format(time.RFC3339))
		cmd.PrintInfoMsgf("Is encrtypted: %t", info.Encrypted)
		cmd.PrintInfoMsgf("Network: %s", info.Network)
	}
}
