package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

func SendTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "sender address",
		})

		toArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "receiver address",
		})

		amtArg := c.Float64(cli.Float64Arg{
			Name: "AMOUNT",
			Desc: "the amount to be transferred",
		})
		stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			w, err := openWallet()
			cmd.FatalErrorCheck(err)

			opts := []wallet.TxOption{
				wallet.OptionStamp(*stampOpt),
				wallet.OptionFee(util.CoinToChange(*feeOpt)),
				wallet.OptionSequence(int32(*seqOpt)),
				wallet.OptionMemo(*memoOpt),
			}

			trx, err := w.MakeSendTx(*fromArg, *toArg, util.CoinToChange(*amtArg),
				opts...)
			cmd.FatalErrorCheck(err)

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign this \033[1mSend\033[0m transition:")
			cmd.PrintInfoMsg("From  : %s", *fromArg)
			cmd.PrintInfoMsg("To    : %s", *toArg)
			cmd.PrintInfoMsg("Amount: %.9f", *amtArg)
			cmd.PrintInfoMsg("Fee   : %.9f", util.ChangeToCoin(trx.Fee()))

			signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
		}
	}
}

func BondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "sender account address",
		})

		toArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "receiver validator address",
		})

		amtArg := c.Float64(cli.Float64Arg{
			Name: "STAKE",
			Desc: "stake amount",
		})

		pubKeyOpt := c.String(cli.StringOpt{
			Name: "pub",
			Desc: "validator public key",
		})

		stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
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
			cmd.PrintInfoMsg("You are going to sign this \033[1mBond\033[0m transition:")
			cmd.PrintInfoMsg("Account  : %s", *fromArg)
			cmd.PrintInfoMsg("Validator: %s", *toArg)
			cmd.PrintInfoMsg("Amount   : %.9f", *amtArg)
			cmd.PrintInfoMsg("Fee      : %.9f", util.ChangeToCoin(trx.Fee()))

			signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
		}
	}
}

func UnbondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "validator's address",
		})
		stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
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
			cmd.PrintInfoMsg("You are going to sign this \033[1mUnbond\033[0m transition:")
			cmd.PrintInfoMsg("Validator: %s", *fromArg)
			cmd.PrintInfoMsg("Fee      : %.9f", util.ChangeToCoin(trx.Fee()))

			signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
		}
	}
}

func WithdrawTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		fromArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "withdraw from Validator address",
		})

		toArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "deposit to account address",
		})

		amtArg := c.Float64(cli.Float64Arg{
			Name: "AMOUNT",
			Desc: "the amount to be transferred",
		})
		stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
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
			cmd.PrintInfoMsg("You are going to sign this \033[1mWithdraw\033[0m transition:")
			cmd.PrintInfoMsg("Validator: %s", *fromArg)
			cmd.PrintInfoMsg("Account  : %s", *toArg)
			cmd.PrintInfoMsg("Amount   : %.9f", *amtArg)
			cmd.PrintInfoMsg("Fee      : %.9f", util.ChangeToCoin(trx.Fee()))

			signAndPublishTx(w, trx, *noConfirmOpt, *passOpt)
		}
	}
}

func addCommonTxOptions(c *cli.Cmd) (*string, *int, *float64, *string, *bool) {
	stampOpt := c.String(cli.StringOpt{
		Name:  "stamp",
		Desc:  "transaction stamp, if not specified will query from gRPC server",
		Value: "",
	})
	seqOpt := c.Int(cli.IntOpt{
		Name:  "seq",
		Desc:  "transaction sequence, if not specified will query from gRPC server",
		Value: 0,
	})
	feeOpt := c.Float64(cli.Float64Opt{
		Name:  "fee",
		Desc:  "transaction fee, if not specified will calculate automatically",
		Value: 0,
	})
	memoOpt := c.String(cli.StringOpt{
		Name:  "memo",
		Desc:  "transaction memo, maximum should be 64 character (optional)",
		Value: "",
	})
	noConfirmOpt := c.Bool(cli.BoolOpt{
		Name:  "no-confirm",
		Desc:  "no confirmation question (optional)",
		Value: false,
	})

	return stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt
}

func signAndPublishTx(w *wallet.Wallet, trx *tx.Tx, noConfirm bool, pass string) {
	cmd.PrintLine()
	password := getPassword(w, pass)
	err := w.SignTransaction(password, trx)
	cmd.FatalErrorCheck(err)

	bs, _ := trx.Bytes()
	cmd.PrintInfoMsg("Signed transaction data: %x", bs)
	cmd.PrintLine()

	if !w.IsOffline() {
		if !noConfirm {
			cmd.PrintInfoMsg("You are going to broadcast the signed transition:")
			cmd.PrintWarnMsg("THIS ACTION IS NOT REVERSIBLE")
			confirmed := cmd.PromptConfirm("Do you want to continue")
			if !confirmed {
				return
			}
		}
		res, err := w.BroadcastTransaction(trx)
		cmd.FatalErrorCheck(err)

		err = w.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsg("Transaction hash: %s", res)
	}
}

func getPassword(wallet *wallet.Wallet, passOpt string) string {
	password := passOpt
	if wallet.IsEncrypted() && password == "" {
		password = cmd.PromptPassword("Wallet password", false)
	}
	return password
}
