package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/spf13/cobra"
)

func buildNeuterCmd(parentCmd *cobra.Command) {
	neuterCmd := &cobra.Command{
		Use:   "neuter",
		Short: "convert full wallet to read-only wallet and can only be used to retrieve balances or stakes",
	}
	parentCmd.AddCommand(neuterCmd)

	neuterCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		neuteredWallet := wlt.Neuter()

		err = neuteredWallet.Save()
		cmd.FatalErrorCheck(err)
	}
}
