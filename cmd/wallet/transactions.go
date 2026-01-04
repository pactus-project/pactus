package main

import (
	"context"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
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

		wlt, err := openWallet(context.Background())
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

	directionOpt := listCmd.Flags().String("direction", "any",
		"filter transactions by direction: any, incoming, outgoing")
	countOpt := listCmd.Flags().Int("count", 10, "number of transactions to list")
	skipOpt := listCmd.Flags().Int("skip", 0, "number of transactions to skip")

	listCmd.Run = func(_ *cobra.Command, args []string) {
		var direction wallet.TxDirection
		switch *directionOpt {
		case "any":
			direction = wallet.TxDirectionAny
		case "incoming":
			direction = wallet.TxDirectionIncoming
		case "outgoing":
			direction = wallet.TxDirectionOutgoing
		default:
			terminal.PrintErrorMsgf("invalid direction: %s", *directionOpt)

			return
		}

		opts := []wallet.ListTransactionsOption{
			wallet.WithAddress(args[0]),
			wallet.WithDirection(direction),
			wallet.WithCount(*countOpt),
			wallet.WithSkip(*skipOpt),
		}

		wlt, err := openWallet(context.Background())
		terminal.FatalErrorCheck(err)

		transactions := wlt.ListTransactions(opts...)
		for i, trx := range transactions {
			terminal.PrintInfoMsgf("%d %v %v %v %v\t%v",
				i+1, trx.CreatedAt.Format("02 Jan 06 15:04"), trx.ID[:12], trx.PayloadType, trx.Status, trx.Amount)
		}
	}
}
