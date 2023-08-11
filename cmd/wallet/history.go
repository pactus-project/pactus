package main

import (
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/spf13/cobra"
)

func buildAllHistoryCmd(parentCmd *cobra.Command) {
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Check the wallet history",
	}

	parentCmd.AddCommand(historyCmd)
	buildAddToHistoryCmd(historyCmd)
	buildShowHistoryCmd(historyCmd)
}

func buildAddToHistoryCmd(parentCmd *cobra.Command) {
	var addToHistoryCmd = &cobra.Command{
		Use:   "add",
		Short: "Add a transaction to the wallet history",
	}
	parentCmd.AddCommand(addToHistoryCmd)

	txID := addToHistoryCmd.Flags().String("ID", "", "transaction id")

	addToHistoryCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		id, err := hash.FromString(*txID)
		cmd.FatalErrorCheck(err)

		err = wallet.AddTransaction(id)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsg("Transaction added to wallet")
	}
}

func buildShowHistoryCmd(parentCmd *cobra.Command) {
	showHistoryCmd := &cobra.Command{
		Use:   "get",
		Short: "Show the transaction history of any address",
	}
	parentCmd.AddCommand(showHistoryCmd)
	addrArg := addAddressArg(parentCmd)

	showHistoryCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		history := wallet.GetHistory(*addrArg)
		for i, h := range history {
			if h.Time != nil {
				cmd.PrintInfoMsg("%d %v %v %v %s\t%v",
					i+1, h.Time.Format(time.RFC822), h.TxID, h.PayloadType, h.Desc, h.Amount)
			} else {
				cmd.PrintInfoMsg("%d              %v  %s\t%v",
					i+1, h.TxID, h.Desc, h.Amount)
			}
		}
	}
}
