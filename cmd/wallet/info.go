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

		cmd.PrintInfoMsgf("version: %d", wlt.Version())
		cmd.PrintInfoMsgf("created at: %s", wlt.CreationTime().Format(time.RFC3339))
		cmd.PrintInfoMsgf("is encrtypted: %t", wlt.IsEncrypted())
		cmd.PrintInfoMsgf("network: %s", wlt.Network().String())
	}
}
