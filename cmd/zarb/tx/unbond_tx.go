package tx

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func UnbondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		stampOpt := c.String(cli.StringOpt{
			Name: "stamp",
			Desc: "Transaction stamp",
		})

		seqOpt := c.Int(cli.IntOpt{
			Name: "seq",
			Desc: "Transaction sequence number",
		})

		valOpt := c.String(cli.StringOpt{
			Name: "val",
			Desc: "Validator's address",
		})

		memoOpt := c.String(cli.StringOpt{
			Name:  "memo",
			Desc:  "Transaction memo",
			Value: "",
		})
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			var err error
			var stamp crypto.Hash
			var validator crypto.Address
			var seq int

			// ---
			if *seqOpt == 0 {
				cmd.PrintWarnMsg("Sequence number is not defined.")
				c.PrintHelp()
				return
			}
			seq = *seqOpt

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

			if *valOpt == "" {
				cmd.PrintWarnMsg("Validator address is not defined.")
				c.PrintHelp()
				return
			}
			validator, err = crypto.AddressFromString(*valOpt)
			if err != nil {
				cmd.PrintErrorMsg("Validator address is not valid: %v", err)
				return
			}

			trx := tx.NewUnbondTx(stamp, seq, validator, *memoOpt)
			bz, _ := trx.Encode()
			cmd.PrintInfoMsg("Unsigned transaction raw bytes:\n%x", bz)

		}
	}
}
