package main

import (
	"context"
	"time"

	"github.com/pactus-project/pactus/util/terminal"
	"github.com/spf13/cobra"
)

// buildInfoCmd builds all sub-commands related to the wallet information.
func buildInfoCmd(parentCmd *cobra.Command) {
	infoCmd := &cobra.Command{
		Use:   "info",
		Short: "retrieving the wallet information",
	}

	parentCmd.AddCommand(infoCmd)

	infoCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet(context.Background())
		terminal.FatalErrorCheck(err)

		info := wlt.Info()

		terminal.PrintInfoMsgf("Version: %d", info.Version)
		terminal.PrintInfoMsgf("UUID: %s", info.UUID)
		terminal.PrintInfoMsgf("Driver: %s", info.Driver)
		terminal.PrintInfoMsgf("Created At: %s", info.CreatedAt.Format(time.RFC1123))
		terminal.PrintInfoMsgf("Default Fee: %s", info.DefaultFee.String())
		terminal.PrintInfoMsgf("Encrypted: %t", info.Encrypted)
		terminal.PrintInfoMsgf("Neutered: %t", info.Neutered)
		terminal.PrintInfoMsgf("Network: %s", info.Network)
	}
}
