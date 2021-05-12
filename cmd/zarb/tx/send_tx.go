package tx

import (
	"encoding/json"
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func SendTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		stampOpt := c.String(cli.StringOpt{
			Name: "stamp",
			Desc: "Transaction stamp",
		})

		seqOpt := c.Int(cli.IntOpt{
			Name:      "seq",
			Desc:      "Transaction sequence number",
			HideValue: true,
		})

		senderOpt := c.String(cli.StringOpt{
			Name: "sender",
			Desc: "Sender address",
		})

		receiverOpt := c.String(cli.StringOpt{
			Name: "receiver",
			Desc: "Receiver address",
		})

		amountOpt := c.Int(cli.IntOpt{
			Name:      "amount",
			Desc:      "The amount to be transferred",
			HideValue: true,
		})

		feeOpt := c.Int(cli.IntOpt{
			Name:      "fee",
			Desc:      "Transaction fee",
			HideValue: true,
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
			var sender crypto.Address
			var receiver crypto.Address
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
				cmd.PrintWarnMsg("Stake is not defined.")
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

			if *senderOpt == "" {
				cmd.PrintWarnMsg("Sender address is not defined.")
				c.PrintHelp()
				return
			}
			sender, err = crypto.AddressFromString(*senderOpt)
			if err != nil {
				cmd.PrintErrorMsg("Sender address is not valid: %v", err)
				return
			}

			if *receiverOpt == "" {
				cmd.PrintWarnMsg("Receiver address is not defined.")
				c.PrintHelp()
				return
			}
			receiver, err = crypto.AddressFromString(*receiverOpt)
			if err != nil {
				cmd.PrintErrorMsg("Receiver address is not valid: %v", err)
				return
			}

			trx := tx.NewSendTx(stamp, seq, sender, receiver, amount, fee, *memoOpt)
			bz, _ := trx.Encode()
			js, _ := json.MarshalIndent(trx, " ", " ")
			cmd.PrintInfoMsg("Transaction format:\n%s", js)
			cmd.PrintLine()
			cmd.PrintInfoMsg("Transaction raw bytes:\n%x", bz)

		}
	}
}
