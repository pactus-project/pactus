package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildAllTransactionCmd builds all sub-commands related to the transactions.
func buildAllTransactionCmd(parentCmd *cobra.Command) {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "creating, signing and publishing a transaction",
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
		cmd.FatalErrorCheck(err)

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFeeFromString(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTimeOpt)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeTransferTx(sender, receiver, amt, opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mTransfer\033[0m transition:")
		cmd.PrintInfoMsgf("From  : %s", sender)
		cmd.PrintInfoMsgf("To    : %s", receiver)
		cmd.PrintInfoMsgf("Amount: %s", amt)
		cmd.PrintInfoMsgf("Fee   : %s", trx.Fee())
		cmd.PrintInfoMsgf("Memo  : %s", trx.Memo())

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
		cmd.FatalErrorCheck(err)

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFeeFromString(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTime)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeBondTx(sender, receiver, *pubKeyOpt, amt, opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mBond\033[0m transition:")
		cmd.PrintInfoMsgf("Account  : %s", sender)
		cmd.PrintInfoMsgf("Validator: %s", receiver)
		cmd.PrintInfoMsgf("Stake    : %s", amt)
		cmd.PrintInfoMsgf("Fee      : %s", trx.Fee())
		cmd.PrintInfoMsgf("Memo     : %s", trx.Memo())

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
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFeeFromString(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTime)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeUnbondTx(from, opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mUnbond\033[0m transition:")
		cmd.PrintInfoMsgf("Validator: %s", from)
		cmd.PrintInfoMsgf("Fee      : %s", trx.Fee())
		cmd.PrintInfoMsgf("Memo     : %s", trx.Memo())

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
		cmd.FatalErrorCheck(err)

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionFeeFromString(*feeOpt),
			wallet.OptionLockTime(uint32(*lockTime)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := wlt.MakeWithdrawTx(sender, receiver, amt, opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mWithdraw\033[0m transition:")
		cmd.PrintInfoMsgf("Validator: %s", sender)
		cmd.PrintInfoMsgf("Account  : %s", receiver)
		cmd.PrintInfoMsgf("Amount   : %s", amt)
		cmd.PrintInfoMsgf("Fee      : %s", trx.Fee())
		cmd.PrintInfoMsgf("Memo     : %s", trx.Memo())

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
	cmd.PrintLine()
	password := getPassword(wlt, pass)
	err := wlt.SignTransaction(password, trx)
	cmd.FatalErrorCheck(err)

	bs, _ := trx.Bytes()
	cmd.PrintInfoMsgf("Signed transaction data: %x", bs)
	cmd.PrintLine()

	if !wlt.IsOffline() {
		if !noConfirm {
			cmd.PrintInfoMsgf("You are going to broadcast the signed transition:")
			cmd.PrintWarnMsgf("THIS ACTION IS NOT REVERSIBLE")
			confirmed := cmd.PromptConfirm("Do you want to continue")
			if !confirmed {
				return
			}
		}
		res, err := wlt.BroadcastTransaction(trx)
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("Transaction hash: %s", res)
	}
}

func getPassword(wlt *wallet.Wallet, passOpt string) string {
	password := passOpt
	if wlt.IsEncrypted() && password == "" {
		password = cmd.PromptPassword("Wallet password", false)
	}

	return password
}
