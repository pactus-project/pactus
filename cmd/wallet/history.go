package main

import (
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/spf13/cobra"
)

// buildAllHistoryCmd builds all sub-commands related to the wallet history.
func buildAllHistoryCmd(parentCmd *cobra.Command) {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "retrieving the transaction history of the wallet.",
	}

	parentCmd.AddCommand(historyCmd)
	buildAddToHistoryCmd(historyCmd)
	buildShowHistoryCmd(historyCmd)
}

// buildAddToHistoryCmd builds a command for adding a transaction to the wallet's history.
func buildAddToHistoryCmd(parentCmd *cobra.Command) {
	addToHistoryCmd := &cobra.Command{
		Use:   "add [flags] <ID>",
		Short: "adds a transaction to the wallet's history.",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(addToHistoryCmd)

	addToHistoryCmd.Run = func(_ *cobra.Command, args []string) {
		txID := args[0]

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		id, err := hash.FromString(txID)
		cmd.FatalErrorCheck(err)

		err = wlt.AddTransaction(id)
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("Transaction successfully added to the wallet.")
	}
}

// buildShowHistoryCmd builds a command for displaying the transaction history of a specific address.
func buildShowHistoryCmd(parentCmd *cobra.Command) {
	showHistoryCmd := &cobra.Command{
		Use:   "get [flags] <ADDRESS>",
		Short: "displays the transaction history for a given address.",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(showHistoryCmd)

	showHistoryCmd.Run = func(_ *cobra.Command, args []string) {
		addr := args[0]

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		history := wlt.GetHistory(addr)
		for i, h := range history {
			if h.Time != nil {
				cmd.PrintInfoMsgf("%d %v %v %v %s\t%v",
					i+1, h.Time.Format(time.RFC822), h.TxID, h.PayloadType, h.Desc, h.Amount)
			} else {
				cmd.PrintInfoMsgf("%d              %v  %s\t%v",
					i+1, h.TxID, h.Desc, h.Amount)
			}
		}
	}
}
