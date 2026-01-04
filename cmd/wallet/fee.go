package main

import (
	"context"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/spf13/cobra"
)

// buildFeeCmd builds sub-command to set the default fee for the wallet.
func buildFeeCmd(parentCmd *cobra.Command) {
	feeCmd := &cobra.Command{
		Use:   "fee [flags] <AMOUNT>",
		Short: "set the default fee for the wallet",
		Args:  cobra.ExactArgs(1),
	}

	parentCmd.AddCommand(feeCmd)

	feeCmd.Run = func(_ *cobra.Command, args []string) {
		wlt, err := openWallet(context.Background())
		terminal.FatalErrorCheck(err)

		fee, err := amount.FromString(args[0])
		terminal.FatalErrorCheck(err)

		err = wlt.SetDefaultFee(fee)
		terminal.FatalErrorCheck(err)

		terminal.PrintInfoMsgf("Default fee is set to %s", fee.String())
	}
}
