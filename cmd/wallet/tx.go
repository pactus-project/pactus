package main

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/tx/payload"
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

		amountArg := c.String(cli.StringArg{
			Name: "AMOUNT",
			Desc: "the amount to be transferred",
		})
		stampOpt, seqOpt, memoOpt, feeOpt := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			trx, err := wallet.MakeSendTx(*stampOpt, *seqOpt, *fromArg, *toArg, *amountArg, *feeOpt, *memoOpt)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast a Send transition to the network:")
			cmd.PrintInfoMsg("From: %s", *fromArg)
			cmd.PrintInfoMsg("To: %s", *toArg)
			cmd.PrintInfoMsg("Amount: %s", *amountArg)

			signAndPublishTx(wallet, trx)
		}
	}
}

func BondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		senderArg := c.String(cli.StringArg{
			Name: "FROM",
			Desc: "sender account address",
		})

		pubArg := c.String(cli.StringArg{
			Name: "TO",
			Desc: "validator public key",
		})

		stakeArg := c.String(cli.StringArg{
			Name: "STAKE",
			Desc: "stake amount",
		})
		stampOpt, seqOpt, memoOpt, feeOpt := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			trx, err := wallet.MakeBondTx(*stampOpt, *seqOpt, *senderArg, *pubArg, *stakeArg, *feeOpt, *memoOpt)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast a bond transition to the network.")
			cmd.PrintInfoMsg("Account: %s", *senderArg)
			cmd.PrintInfoMsg("Validator: %s", trx.Payload().(*payload.BondPayload).PublicKey.Address())
			cmd.PrintInfoMsg("Stake: %s", *stakeArg)

			signAndPublishTx(wallet, trx)
		}
	}
}

func UnbondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		valArg := c.String(cli.StringArg{
			Name: "ADDR",
			Desc: "validator's address",
		})
		stampOpt, seqOpt, memoOpt, _ := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			trx, err := wallet.MakeUnbondTx(*stampOpt, *seqOpt, *valArg, *memoOpt)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast an Unbond transition to the network:")
			cmd.PrintInfoMsg("Validator: %s", *valArg)

			signAndPublishTx(wallet, trx)

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

		amountArg := c.String(cli.StringArg{
			Name: "AMOUNT",
			Desc: "the amount to be transferred",
		})
		stampOpt, seqOpt, memoOpt, feeOpt := addCommonTxOptions(c)

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			wallet, err := wallet.OpenWallet(*path)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			trx, err := wallet.MakeWithdrawTx(*stampOpt, *seqOpt, *fromArg, *toArg, *amountArg, *feeOpt, *memoOpt)
			if err != nil {
				cmd.PrintDangerMsg(err.Error())
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("You are going to sign and broadcast a Withdraw transition to the network.")
			cmd.PrintInfoMsg("Validator: %s", *fromArg)
			cmd.PrintInfoMsg("Account: %s", *toArg)
			cmd.PrintInfoMsg("Amount: %s", *amountArg)

			signAndPublishTx(wallet, trx)
		}
	}
}

func addCommonTxOptions(c *cli.Cmd) (*string, *string, *string, *string) {
	stampOpt := c.String(cli.StringOpt{
		Name: "stamp",
		Desc: "transaction stamp, if not specified will query from gRPC server",
	})

	seqOpt := c.String(cli.StringOpt{
		Name: "seq",
		Desc: "transaction sequence, if not specified will query from gRPC server",
	})
	memoOpt := c.String(cli.StringOpt{
		Name:  "memo",
		Desc:  "transaction memo, maximum should be 64 character (optional)",
		Value: "",
	})
	feeOpt := c.String(cli.StringOpt{
		Name:  "fee",
		Desc:  "transaction fee, if not specified will calculate automatically",
		Value: "",
	})

	return stampOpt, seqOpt, memoOpt, feeOpt
}

func signAndPublishTx(wallet *wallet.Wallet, trx *tx.Tx) {
	cmd.PrintWarnMsg("THIS ACTION IS NOT REVERSIBLE")
	confirmed := cmd.PromptConfirm("Do you want to continue? ")
	if !confirmed {
		return
	}

	passphrase := getPassphrase(wallet)
	res, err := wallet.SignAndBroadcast(passphrase, trx)
	if err != nil {
		cmd.PrintDangerMsg(err.Error())
		return
	}
	cmd.PrintInfoMsg(res)
}

func getPassphrase(wallet *wallet.Wallet) string {
	passphrase := ""
	if wallet.IsEncrypted() {
		passphrase = cmd.PromptPassphrase("Wallet password: ", false)
	}
	return passphrase
}
