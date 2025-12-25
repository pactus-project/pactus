package main

import (
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/spf13/cobra"
)

// buildTransactionsCmd builds all sub-commands related to the wallet transactions.
func buildTransactionsCmd(parentCmd *cobra.Command) {
	transactionsCmd := &cobra.Command{
		Use:   "transactions",
		Short: "retrieving the list of transaction for the wallet",
	}

	parentCmd.AddCommand(transactionsCmd)
	buildTransactionsAddCmd(transactionsCmd)
	buildTransactionsListCmd(transactionsCmd)
}

// buildTransactionsAddCmd builds a command for adding a transaction to the wallet.
func buildTransactionsAddCmd(parentCmd *cobra.Command) {
	addCmd := &cobra.Command{
		Use:   "add [flags] <ID>",
		Short: "Add a transaction to the wallet",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(addCmd)

	addCmd.Run = func(_ *cobra.Command, args []string) {
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

// buildTransactionsListCmd builds a command for listing transactions of a specific address.
func buildTransactionsListCmd(parentCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list [flags] <ADDRESS>",
		Short: "List transactions for a given address",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(listCmd)

	listCmd.Run = func(_ *cobra.Command, args []string) {
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
