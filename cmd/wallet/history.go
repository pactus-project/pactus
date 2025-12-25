package main

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/terminal"
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
		terminal.FatalErrorCheck(err)

		id, err := hash.FromString(txID)
		terminal.FatalErrorCheck(err)

		err = wlt.AddTransaction(id)
		terminal.FatalErrorCheck(err)

		terminal.PrintInfoMsgf("Transaction successfully added to the wallet.")
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
		terminal.FatalErrorCheck(err)

		transactions := wlt.ListTransactions(addr)
		for i, trx := range transactions {
			terminal.PrintInfoMsgf("%d %v %v %v %v\t%v",
				i+1, trx.CreatedAt.Format("02 Jan 06 15:04"), trx.ID[:12], trx.PayloadType, trx.Status, trx.Amount)
		}
	}
}
