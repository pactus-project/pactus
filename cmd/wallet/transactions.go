package main

import (
	"context"
	"fmt"

	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/types"
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
		defer wlt.Close()

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
		var direction types.TxDirection
		switch *directionOpt {
		case "any":
			direction = types.TxDirectionAny
		case "incoming":
			direction = types.TxDirectionIncoming
		case "outgoing":
			direction = types.TxDirectionOutgoing
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
		defer wlt.Close()

		transactions := wlt.ListTransactions(opts...)
		const (
			noWidth      = 4
			timeWidth    = 16
			idWidth      = 65
			addressWidth = 16
			amountWidth  = 12
			typeWidth    = 10
			statusWidth  = 10
		)

		headerFmt := fmt.Sprintf("%%-%dv %%-%dv %%-%dv %%-%dv %%-%dv %%-%dv %%-%dv %%-%dv",
			noWidth, timeWidth, idWidth, addressWidth, addressWidth, amountWidth, typeWidth, statusWidth)
		rowFmt := headerFmt

		terminal.PrintInfoMsgBoldf(headerFmt,
			"No", "Time", "ID", "Sender", "Receiver", "Amount", "Type", "Status")

		fit := func(s string, width int) string {
			if len(s) > width {
				if width <= 3 {
					return s[:width]
				}

				return fmt.Sprintf("%-*s", width, s[:width-3]+"...")
			}

			return fmt.Sprintf("%-*s", width, s)
		}

		shortAddr := func(addr string) string {
			const keep = 6
			if len(addr) <= (keep*2)+3 {
				return addr
			}
			return fmt.Sprintf("%s...%s", addr[:keep], addr[len(addr)-keep:])
		}

		for i, trx := range transactions {
			terminal.PrintInfoMsgf(rowFmt,
				i+1,
				trx.CreatedAt.Format("2/1/2006 15:04"),
				fit(trx.ID, idWidth),
				fit(shortAddr(trx.Sender), addressWidth),
				fit(shortAddr(trx.Receiver), addressWidth),
				fit(trx.Amount.String(), amountWidth),
				fit(trx.PayloadType.String(), typeWidth),
				fit(trx.Status.String(), statusWidth),
			)
		}
	}
}
