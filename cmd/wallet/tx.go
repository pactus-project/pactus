package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

func buildAllTransactionCmd(parentCmd *cobra.Command) {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Create, sign and publish a transaction",
	}

	parentCmd.AddCommand(txCmd)
	buildTransferTxCmd(txCmd)
	buildBondTxCmd(txCmd)
	buildUnbondTxCmd(txCmd)
	buildWithdrawTxCmd(txCmd)
}

func buildTransferTxCmd(parentCmd *cobra.Command) {
	transferCmd := &cobra.Command{
		Use:   "transfer",
		Short: "Create, sign and publish a Transfer transaction",
	}
	parentCmd.AddCommand(transferCmd)

	fromArg := transferCmd.Flags().String("FROM", "", "sender address")
	toArg := transferCmd.Flags().String("TO", "", "receiver address")
	amtArg := transferCmd.Flags().Float64("AMOUNT", 0, "the amount to be transferred")

	stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(parentCmd)
	passOpt := addPasswordOption(parentCmd)

	transferCmd.Run = func(_ *cobra.Command, _ []string) {
		w, err := openWallet()
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionStamp(*stampOpt),
			wallet.OptionFee(util.CoinToChange(*feeOpt)),
			wallet.OptionSequence(int32(*seqOpt)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := w.MakeTransferTx(*fromArg, *toArg, util.CoinToChange(*amtArg),
			opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mSend\033[0m transition:")
		cmd.PrintInfoMsgf("From  : %s", *fromArg)
		cmd.PrintInfoMsgf("To    : %s", *toArg)
		cmd.PrintInfoMsgf("Amount: %.9f", *amtArg)
		cmd.PrintInfoMsgf("Fee   : %.9f", util.ChangeToCoin(trx.Fee()))

		signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
	}
}

func buildBondTxCmd(parentCmd *cobra.Command) {
	bondTransactionCmd := &cobra.Command{
		Use:   "bond",
		Short: "Create, sign and publish a Bond transaction",
	}
	parentCmd.AddCommand(bondTransactionCmd)

	fromArg := bondTransactionCmd.Flags().String("FROM", "", "sender address")
	toArg := bondTransactionCmd.Flags().String("TO", "", "receiver validator address")
	amtArg := bondTransactionCmd.Flags().Float64("STAKE", 0, "stake amount")
	pubKeyOpt := bondTransactionCmd.Flags().String("pub", "", "validator public key")

	stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(parentCmd)
	passOpt := addPasswordOption(parentCmd)

	bondTransactionCmd.Run = func(_ *cobra.Command, _ []string) {
		w, err := openWallet()
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionStamp(*stampOpt),
			wallet.OptionFee(util.CoinToChange(*feeOpt)),
			wallet.OptionSequence(int32(*seqOpt)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := w.MakeBondTx(*fromArg, *toArg, *pubKeyOpt,
			util.CoinToChange(*amtArg), opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mBond\033[0m transition:")
		cmd.PrintInfoMsgf("Account  : %s", *fromArg)
		cmd.PrintInfoMsgf("Validator: %s", *toArg)
		cmd.PrintInfoMsgf("Amount   : %.9f", *amtArg)
		cmd.PrintInfoMsgf("Fee      : %.9f", util.ChangeToCoin(trx.Fee()))

		signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
	}
}

func buildUnbondTxCmd(parentCmd *cobra.Command) {
	unbondTransactionCmd := &cobra.Command{
		Use:   "unbond",
		Short: "Create, sign and publish an Unbond transaction",
	}
	parentCmd.AddCommand(unbondTransactionCmd)

	fromArg := unbondTransactionCmd.Flags().String("ADDR", "",
		"validator's address")

	stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(parentCmd)
	passOpt := addPasswordOption(parentCmd)

	unbondTransactionCmd.Run = func(_ *cobra.Command, _ []string) {
		w, err := openWallet()
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionStamp(*stampOpt),
			wallet.OptionFee(util.CoinToChange(*feeOpt)),
			wallet.OptionSequence(int32(*seqOpt)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := w.MakeUnbondTx(*fromArg, opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mUnbond\033[0m transition:")
		cmd.PrintInfoMsgf("Validator: %s", *fromArg)
		cmd.PrintInfoMsgf("Fee      : %.9f", util.ChangeToCoin(trx.Fee()))

		signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
	}
}

func buildWithdrawTxCmd(parentCmd *cobra.Command) {
	withdrawTransactionCmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Create, sign and publish a Withdraw transaction",
	}
	parentCmd.AddCommand(withdrawTransactionCmd)

	fromArg := withdrawTransactionCmd.Flags().String("FROM",
		"", "withdraw from Validator address")

	toArg := withdrawTransactionCmd.Flags().String("TO",
		"", "deposit to account address")

	amtArg := withdrawTransactionCmd.Flags().Float64("AMOUNT",
		0, "the amount to be transferred")

	stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(parentCmd)
	passOpt := addPasswordOption(parentCmd)

	withdrawTransactionCmd.Run = func(_ *cobra.Command, _ []string) {
		w, err := openWallet()
		cmd.FatalErrorCheck(err)

		opts := []wallet.TxOption{
			wallet.OptionStamp(*stampOpt),
			wallet.OptionFee(util.CoinToChange(*feeOpt)),
			wallet.OptionSequence(int32(*seqOpt)),
			wallet.OptionMemo(*memoOpt),
		}

		trx, err := w.MakeWithdrawTx(*fromArg, *toArg,
			util.CoinToChange(*amtArg), opts...)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("You are going to sign this \033[1mWithdraw\033[0m transition:")
		cmd.PrintInfoMsgf("Validator: %s", *fromArg)
		cmd.PrintInfoMsgf("Account  : %s", *toArg)
		cmd.PrintInfoMsgf("Amount   : %.9f", *amtArg)
		cmd.PrintInfoMsgf("Fee      : %.9f", util.ChangeToCoin(trx.Fee()))

		signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
	}
}

func addCommonTxOptions(c *cobra.Command) (*string, *int, *float64, *string, *bool) {
	stampOpt := c.Flags().String("stamp", "",
		"transaction stamp, if not specified will query from gRPC server")

	seqOpt := c.Flags().Int("seq", 0,
		"transaction sequence, if not specified will query from gRPC server")

	feeOpt := c.Flags().Float64("fee", 0,
		"transaction fee, if not specified will calculate automatically")

	memoOpt := c.Flags().String("memo", "",
		"transaction memo, maximum should be 64 character (optional)")

	noConfirmOpt := c.Flags().Bool("no-confirm", false,
		"no confirmation question (optional)")

	return stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt
}

func signAndPublishTx(w *wallet.Wallet, trx *tx.Tx, noConfirm bool, pass string) {
	cmd.PrintLine()
	password := getPassword(w, pass)
	err := w.SignTransaction(password, trx)
	cmd.FatalErrorCheck(err)

	bs, _ := trx.Bytes()
	cmd.PrintInfoMsgf("Signed transaction data: %x", bs)
	cmd.PrintLine()

	if !w.IsOffline() {
		if !noConfirm {
			cmd.PrintInfoMsgf("You are going to broadcast the signed transition:")
			cmd.PrintWarnMsgf("THIS ACTION IS NOT REVERSIBLE")
			confirmed := cmd.PromptConfirm("Do you want to continue")
			if !confirmed {
				return
			}
		}
		res, err := w.BroadcastTransaction(trx)
		cmd.FatalErrorCheck(err)

		err = w.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("Transaction hash: %s", res)
	}
}

func getPassword(wallet *wallet.Wallet, passOpt string) string {
	password := passOpt
	if wallet.IsEncrypted() && password == "" {
		password = cmd.PromptPassword("Wallet password", false)
	}
	return password
}
