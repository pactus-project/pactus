package tx

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
)

func BondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		stampOpt := c.String(cli.StringOpt{
			Name: "stamp",
			Desc: "Transaction stamp",
		})

		seqOpt := c.Int(cli.IntOpt{
			Name: "seq",
			Desc: "Transaction sequence number",
		})

		bonderOpt := c.String(cli.StringOpt{
			Name: "bonder",
			Desc: "Bonder address",
		})

		pubOpt := c.String(cli.StringOpt{
			Name: "pub",
			Desc: "Validator's public key",
		})

		stakeOpt := c.Int(cli.IntOpt{
			Name: "stake",
			Desc: "Stake amount",
		})

		feeOpt := c.Int(cli.IntOpt{
			Name: "fee",
			Desc: "Transaction fee",
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
			var bonder crypto.Address
			var pub crypto.PublicKey
			var seq int
			var stake int64
			var fee int64

			// ---
			if *seqOpt == 0 {
				cmd.PrintWarnMsg("Sequence number is not defined.")
				c.PrintHelp()
				return
			}
			seq = *seqOpt

			if *stakeOpt == 0 {
				cmd.PrintWarnMsg("Stake is not defined.")
				c.PrintHelp()
				return
			}
			stake = int64(*stakeOpt)

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

			if *bonderOpt == "" {
				cmd.PrintWarnMsg("bonder is not defined.")
				c.PrintHelp()
				return
			}
			bonder, err = crypto.AddressFromString(*bonderOpt)
			if err != nil {
				cmd.PrintErrorMsg("Bonder's address is wrong: %v", err)
				return
			}

			if *pubOpt == "" {
				cmd.PrintWarnMsg("Public key is not defined.")
				c.PrintHelp()
				return
			}
			pub, err = crypto.PublicKeyFromString(*pubOpt)
			if err != nil {
				cmd.PrintErrorMsg("Validator's public key is wrong: %v", err)
				return
			}

			trx := tx.NewBondTx(stamp, seq, bonder, pub, stake, fee, *memoOpt)
			bz, _ := trx.Encode()
			cmd.PrintInfoMsg("Unsigned transaction raw bytes:\n%x", bz)

		}
	}
}
