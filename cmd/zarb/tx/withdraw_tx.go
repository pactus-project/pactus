package tx

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func WithdrawTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		stampOpt := c.String(cli.StringOpt{
			Name: "stamp",
			Desc: "Transaction stamp",
		})

		seqOpt := c.Int(cli.IntOpt{
			Name: "seq",
			Desc: "Transaction sequence number",
		})

		fromOpt := c.String(cli.StringOpt{
			Name: "from",
			Desc: "withdraw from Validator address",
		})

		toOpt := c.String(cli.StringOpt{
			Name: "to",
			Desc: "Deposit to address",
		})

		amountOpt := c.Int(cli.IntOpt{
			Name: "amount",
			Desc: "The amount to be transferred",
		})

		feeOpt := c.Int(cli.IntOpt{
			Name: "fee",
			Desc: "Transaction fee",
		})

		memoOpt := c.String(cli.StringOpt{
			Name:  "memo",
			Desc:  "Transaction memo (Optional)",
			Value: "",
		})
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			var err error
			var stamp crypto.Hash
			var from crypto.Address
			var to crypto.Address
			var seq int
			var amount int64
			var fee int64

			// ---
			if *seqOpt == 0 {
				cmd.PrintWarnMsg("Sequence number is not defined.")
				c.PrintHelp()
				return
			}
			seq = *seqOpt

			if *amountOpt == 0 {
				cmd.PrintWarnMsg("Amount is not defined.")
				c.PrintHelp()
				return
			}
			amount = int64(*amountOpt)

			if *feeOpt == 0 {
				cmd.PrintWarnMsg("Fee is not defined.")
				c.PrintHelp()
				return
			}
			fee = int64(*feeOpt)

			if *stampOpt == "" {
				cmd.PrintWarnMsg("stamp is not defined.")
				c.PrintHelp()
				return
			}
			stamp, err = crypto.HashFromString(*stampOpt)
			if err != nil {
				cmd.PrintErrorMsg("Stamp is wrong: %v", err)
				return
			}

			if *fromOpt == "" {
				cmd.PrintWarnMsg("Validator address is not defined.")
				c.PrintHelp()
				return
			}
			from, err = crypto.AddressFromString(*fromOpt)
			if err != nil {
				cmd.PrintErrorMsg("Validator address is not valid: %v", err)
				return
			}

			if *toOpt == "" {
				cmd.PrintWarnMsg("Deposit to address is not defined.")
				c.PrintHelp()
				return
			}
			to, err = crypto.AddressFromString(*toOpt)
			if err != nil {
				cmd.PrintErrorMsg("Deposit to address is not valid: %v", err)
				return
			}

			trx := tx.NewWithdrawTx(stamp, seq, from, to, amount, fee, *memoOpt)
			bz, _ := trx.Encode()
			cmd.PrintInfoMsg("Unsigned transaction raw bytes:\n%x", bz)

		}
	}
}
