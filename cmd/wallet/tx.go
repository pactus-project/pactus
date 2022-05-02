package main

import (
	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/types/tx"
	"github.com/zarbchain/zarb-go/wallet"
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
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			opts := []wallet.TxOption{
				wallet.OptionStamp(*stampOpt),
				wallet.OptionFee(coinToChange(*feeOpt)),
				wallet.OptionSequence(int32(*seqOpt)),
				wallet.OptionMemo(*memoOpt),
			}

			trx, err := w.MakeSendTx(*fromArg, *toArg, coinToChange(*amtArg),
				opts...)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast a Send transition to the network:")
			cmd.PrintInfoMsg("From: %s", *fromArg)
			cmd.PrintInfoMsg("To: %s", *toArg)
			cmd.PrintInfoMsg("Amount: %s (%s)", *amtArg, coinToChange(*amtArg))
			cmd.PrintInfoMsg("Fee: %s (%s)", changeToCoin(trx.Fee()), trx.Fee())

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
			Name: "pubkey",
			Desc: "validator public key",
		})

		stampOpt, seqOpt, feeOpt, memoOpt, noConfirmOpt := addCommonTxOptions(c)
		passOpt := addPasswordOption(c)

		c.Before = func() {}
		c.Action = func() {
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			opts := []wallet.TxOption{
				wallet.OptionStamp(*stampOpt),
				wallet.OptionFee(coinToChange(*feeOpt)),
				wallet.OptionSequence(int32(*seqOpt)),
				wallet.OptionMemo(*memoOpt),
			}

			trx, err := w.MakeBondTx(*fromArg, *toArg, *pubKeyOpt,
				coinToChange(*amtArg), opts...)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast a bond transition to the network.")
			cmd.PrintInfoMsg("Account: %s", *fromArg)
			cmd.PrintInfoMsg("Validator: %s", *toArg)
			cmd.PrintInfoMsg("Amount: %v (%v)", *amtArg, coinToChange(*amtArg))
			cmd.PrintInfoMsg("Fee: %v (%v)", changeToCoin(trx.Fee()), trx.Fee())

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
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			opts := []wallet.TxOption{
				wallet.OptionStamp(*stampOpt),
				wallet.OptionFee(coinToChange(*feeOpt)),
				wallet.OptionSequence(int32(*seqOpt)),
				wallet.OptionMemo(*memoOpt),
			}

			trx, err := w.MakeUnbondTx(*fromArg, opts...)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast an Unbond transition to the network:")
			cmd.PrintInfoMsg("Validator: %s", *fromArg)

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
			w, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			opts := []wallet.TxOption{
				wallet.OptionStamp(*stampOpt),
				wallet.OptionFee(coinToChange(*feeOpt)),
				wallet.OptionSequence(int32(*seqOpt)),
				wallet.OptionMemo(*memoOpt),
			}

			trx, err := w.MakeWithdrawTx(*fromArg, *toArg,
				coinToChange(*amtArg), opts...)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast a Withdraw transition to the network.")
			cmd.PrintInfoMsg("Validator: %s", *fromArg)
			cmd.PrintInfoMsg("Account: %s", *toArg)
			cmd.PrintInfoMsg("Amount: %s", *amtArg)

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

func signAndPublishTx(wallet *wallet.Wallet, trx *tx.Tx, noConfirm bool,
	pass string) {
	if !noConfirm {
		cmd.PrintWarnMsg("THIS ACTION IS NOT REVERSIBLE")
		confirmed := cmd.PromptConfirm("Do you want to continue")
		if !confirmed {
			return
		}
	}

	password := getPassword(wallet, pass)
	res, err := wallet.SignAndBroadcast(password, trx)
	if err != nil {
		cmd.PrintDangerMsg(err.Error())
		return
	}
	cmd.PrintInfoMsg(res)
}

func getPassword(wallet *wallet.Wallet, passOpt string) string {
	password := passOpt
	if wallet.IsEncrypted() && password == "" {
		password = cmd.PromptPassword("Wallet password", false)
	}
	return password
}
