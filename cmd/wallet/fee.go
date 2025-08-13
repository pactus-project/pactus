package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/spf13/cobra"
)

// buildFeeCmd builds sub-command to set the default fee for the wallet.
func buildFeeCmd(parentCmd *cobra.Command) {
	feeCmd := &cobra.Command{
		Use:   "fee [flags] <AMOUNT>",
		Short: "set the default fee for the wallet.",
		Args:  cobra.ExactArgs(1),
	}

	parentCmd.AddCommand(feeCmd)

	feeCmd.Run = func(_ *cobra.Command, args []string) {
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		fee, err := amount.FromString(args[0])
		cmd.FatalErrorCheck(err)

		wlt.SetDefaultFee(fee)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("Default fee is set to %s", fee.String())
	}
}
