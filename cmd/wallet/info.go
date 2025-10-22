package main

import (
	"time"

	"github.com/pactus-project/pactus/util/terminal"
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
		terminal.FatalErrorCheck(err)

		info := wlt.Info()

		terminal.PrintInfoMsgf("Version: %d", info.Version)
		terminal.PrintInfoMsgf("UUID: %s", info.UUID)
		terminal.PrintInfoMsgf("Default fee: %s", info.DefaultFee.String())
		terminal.PrintInfoMsgf("Created at: %s", info.CreatedAt.Format(time.RFC3339))
		terminal.PrintInfoMsgf("Is encrtypted: %t", info.Encrypted)
		terminal.PrintInfoMsgf("Network: %s", info.Network)
	}
}
