package main

import (
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildAllTransactionCmd builds all sub-commands related to the transactions.
func buildAllTransactionCmd(parentCmd *cobra.Command) {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "create, sign and publish a transaction",
	}

	parentCmd.AddCommand(txCmd)
	buildTransferTxCmd(txCmd)
	buildBondTxCmd(txCmd)
	buildUnbondTxCmd(txCmd)
	buildWithdrawTxCmd(txCmd)
}

// buildTransferTxCmd builds a command for create, sign and publish a `Transfer` transaction.
func buildTransferTxCmd(parentCmd *cobra.Command) {
	transferCmd := &cobra.Command{
		Use:   "transfer [flags] <FROM> <TO> <AMOUNT>",
		Short: "create, sign and publish a `Transfer` transaction",
		Args:  cobra.ExactArgs(3),
	}
	parentCmd.AddCommand(transferCmd)

	lockTimeOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(transferCmd)
	passOpt := addPasswordOption(transferCmd)

	transferCmd.Run = func(_ *cobra.Command, args []string) {
		sender := args[0]
		receiver := args[1]
		amt, err := amount.FromString(args[2])
		terminal.FatalErrorCheck(err)

		wlt, err := openWallet()
		terminal.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFee(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTimeOpt)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeTransferTx(sender, receiver, amt, opts...)
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintInfoMsgf("📝 Transaction Details:")
		terminal.PrintInfoMsgf("   Type   : Transfer")
		terminal.PrintInfoMsgf("   From   : %s", sender)
		terminal.PrintInfoMsgf("   To     : %s", receiver)
		terminal.PrintInfoMsgf("   Amount : %s", amt)
		terminal.PrintInfoMsgf("   Fee    : %s", trx.Fee())
		terminal.PrintInfoMsgf("   Memo   : %s", trx.Memo())

		signAndPublishTx(wlt, trx, *noConfirmOpt, *passOpt)
	}
}

// buildBondTxCmd builds a command for create, sign and publish a `Bond` transaction.
func buildBondTxCmd(parentCmd *cobra.Command) {
	bondCmd := &cobra.Command{
		Use:   "bond [flags] <ACCOUNT> <VALIDATOR> <STAKE>",
		Short: "create, sign and publish a `Bond` transaction",
		Args:  cobra.ExactArgs(3),
	}
	parentCmd.AddCommand(bondCmd)

	pubKeyOpt := bondCmd.Flags().String("pub", "", "validator's public key")
	lockTime, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(bondCmd)
	passOpt := addPasswordOption(bondCmd)

	bondCmd.Run = func(_ *cobra.Command, args []string) {
		sender := args[0]
		receiver := args[1]
		amt, err := amount.FromString(args[2])
		terminal.FatalErrorCheck(err)

		wlt, err := openWallet()
		terminal.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFee(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTime)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeBondTx(sender, receiver, *pubKeyOpt, amt, opts...)
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintInfoMsgf("📝 Transaction Details:")
		terminal.PrintInfoMsgf("   Type     : Bond")
		terminal.PrintInfoMsgf("   Account  : %s", sender)
		terminal.PrintInfoMsgf("   Validator: %s", receiver)
		terminal.PrintInfoMsgf("   Stake    : %s", amt)
		terminal.PrintInfoMsgf("   Fee      : %s", trx.Fee())
		terminal.PrintInfoMsgf("   Memo     : %s", trx.Memo())

		signAndPublishTx(wlt, trx, *noConfirmOpt, *passOpt)
	}
}

// buildUnbondTxCmd builds a command for create, sign and publish a `Unbond` transaction.
func buildUnbondTxCmd(parentCmd *cobra.Command) {
	unbondCmd := &cobra.Command{
		Use:   "unbond [flags] <ADDRESS>",
		Short: "create, sign and publish an `Unbond` transaction",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(unbondCmd)

	lockTime, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(unbondCmd)
	passOpt := addPasswordOption(unbondCmd)

	unbondCmd.Run = func(_ *cobra.Command, args []string) {
		from := args[0]

		wlt, err := openWallet()
		terminal.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFee(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTime)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeUnbondTx(from, opts...)
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintInfoMsgf("📝 Transaction Details:")
		terminal.PrintInfoMsgf("   Type     : Unbond")
		terminal.PrintInfoMsgf("   Validator: %s", from)
		terminal.PrintInfoMsgf("   Fee      : %s", trx.Fee())
		terminal.PrintInfoMsgf("   Memo     : %s", trx.Memo())

		signAndPublishTx(wlt, trx, *noConfirmOpt, *passOpt)
	}
}

// buildWithdrawTxCmd builds a command for create, sign and publish a `Withdraw` transaction.
func buildWithdrawTxCmd(parentCmd *cobra.Command) {
	withdrawCmd := &cobra.Command{
		Use:   "withdraw [flags] <VALIDATOR> <ACCOUNT> <STAKE>",
		Short: "create, sign and publish a `Withdraw` transaction",
		Args:  cobra.ExactArgs(3),
	}
	parentCmd.AddCommand(withdrawCmd)

	lockTime, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(withdrawCmd)
	passOpt := addPasswordOption(withdrawCmd)

	withdrawCmd.Run = func(_ *cobra.Command, args []string) {
		sender := args[0]
		receiver := args[1]
		amt, err := amount.FromString(args[2])
		terminal.FatalErrorCheck(err)

		wlt, err := openWallet()
		terminal.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFee(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTime)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeWithdrawTx(sender, receiver, amt, opts...)
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintInfoMsgf("📝 Transaction Details:")
		terminal.PrintInfoMsgf("   Type     : Withdraw")
		terminal.PrintInfoMsgf("   Validator: %s", sender)
		terminal.PrintInfoMsgf("   Account  : %s", receiver)
		terminal.PrintInfoMsgf("   Amount   : %s", amt)
		terminal.PrintInfoMsgf("   Fee      : %s", trx.Fee())
		terminal.PrintInfoMsgf("   Memo     : %s", trx.Memo())

		signAndPublishTx(wlt, trx, *noConfirmOpt, *passOpt)
	}
}

func addCommonTxOptions(cobra *cobra.Command) (lockTimeOpt *int, feeOpt, memoOpt *string, noConfirmOpt *bool) {
	lockTimeOpt = cobra.Flags().Int("lock-time", 0,
		"transaction lock-time, if not specified will be the latest height")

	feeOpt = cobra.Flags().String("fee", "",
		"transaction fee in PAC, if not specified will set to the estimated value")

	memoOpt = cobra.Flags().String("memo", "",
		"transaction memo, maximum should be 64 character")

	noConfirmOpt = cobra.Flags().Bool("no-confirm", false,
		"no confirmation question")

	return lockTimeOpt, feeOpt, memoOpt, noConfirmOpt
}

func signAndPublishTx(wlt *wallet.Wallet, trx *tx.Tx, noConfirm bool, pass string) {
	terminal.PrintLine()
	password := getPassword(wlt, pass)
	err := wlt.SignTransaction(password, trx)
	terminal.FatalErrorCheck(err)

	bs, _ := trx.Bytes()
	terminal.PrintLine()
	terminal.PrintSuccessMsgf("✅ Transaction signed successfully")
	terminal.PrintInfoMsgf("   Signed transaction data: %x", bs)
	terminal.PrintLine()

	if !wlt.IsOffline() {
		if !noConfirm {
			terminal.PrintWarnMsgf("⚠️  You are going to broadcast the signed transaction")
			terminal.PrintWarnMsgf("   This action cannot be undone")
			terminal.PrintLine()
			confirmed := prompt.PromptConfirm("Do you want to continue")
			if !confirmed {
				return
			}
		}
		res, err := wlt.BroadcastTransaction(trx)
		terminal.FatalErrorCheck(err)

		err = wlt.Save()
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintSuccessMsgf("✅ Transaction broadcast successfully!")
		terminal.PrintInfoMsgf("   Transaction ID: %s", res)
	}
}

func getPassword(wlt *wallet.Wallet, passOpt string) string {
	password := passOpt
	if wlt.IsEncrypted() && password == "" {
		password = prompt.PromptPassword("Wallet password", false)
	}

	return password
}
