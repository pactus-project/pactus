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

// buildTransactionCmd builds all sub-commands related to the wallet transactions.
func buildTransactionCmd(parentCmd *cobra.Command) {
	transactionCmd := &cobra.Command{
		Use:   "transaction",
		Short: "Manage and view wallet transactions",
		Long:  "The 'transactions' command allows you to list existing transactions with optional filtering",
	}

	parentCmd.AddCommand(transactionCmd)
	buildTransactionAddCmd(transactionCmd)
	buildTransactionListCmd(transactionCmd)
}

// buildTransactionsAddCmd builds the command for adding a transaction to the wallet.
func buildTransactionAddCmd(parentCmd *cobra.Command) {
	addCmd := &cobra.Command{
		Use:   "add [flags] <ID>",
		Short: "Add a transaction to the wallet",
		Long:  "Add a specific transaction to the wallet using its transaction ID",
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

// buildTransactionListCmd builds the command for listing wallet transactions.
func buildTransactionListCmd(parentCmd *cobra.Command) {
	listCmd := &cobra.Command{
		Use:   "list [flags]",
		Short: "List transactions in the wallet",
		Long: `List transactions stored in the wallet.

Examples:
  # List last 10 transactions (default)
  transaction list

  # List last 20 outgoing transactions
  transaction list --direction outgoing --count 20

  # List transactions for a specific address
  transaction list --address pc1xyz...`,
	}
	parentCmd.AddCommand(listCmd)

	directionOpt := listCmd.Flags().String("direction", "any",
		"Filter transactions by direction: 'any', 'incoming', or 'outgoing'. Default is 'any'.")
	addressOpt := listCmd.Flags().String("address", "*",
		"Filter transactions by wallet address. Use '*' to include all addresses.")
	countOpt := listCmd.Flags().Int("count", 10,
		"Maximum number of transactions to display. Default is 10.")
	skipOpt := listCmd.Flags().Int("skip", 0,
		"Number of transactions to skip for pagination. Default is 0.")

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
			wallet.WithAddress(*addressOpt),
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

		fit := func(text string, width int) string {
			if len(text) > width {
				if width <= 3 {
					return text[:width]
				}

				return fmt.Sprintf("%-*s", width, text[:width-3]+"...")
			}

			return fmt.Sprintf("%-*s", width, text)
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
