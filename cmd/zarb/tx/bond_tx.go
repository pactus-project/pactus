package tx

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/tx"
	"github.com/zarbchain/zarb-go/util"
)

func BondTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		stampOpt := c.String(cli.StringOpt{
			Name: "stamp",
			Desc: "Transaction stamp if not specified will query from RPC server",
		})

		seqOpt := c.Int(cli.IntOpt{
			Name: "seq",
			Desc: "Transaction sequence number if not specified will query from RPC server",
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

		authOpt := c.String(cli.StringOpt{
			Name: "a auth",
			Desc: "Passphrase of the key file",
		})
		keyFileOpt := c.String(cli.StringOpt{
			Name: "k keyfile",
			Desc: "Path to the encrypted key file",
		})

		grpcOpt := c.String(cli.StringOpt{
			Name: "e endpoint",
			Desc: "gRPC server address",
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
			var auth string

			// ---
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

			if *bonderOpt == "" {
				cmd.PrintWarnMsg("Bonder address is not defined.")
				c.PrintHelp()
				return
			}
			bonder, err = crypto.AddressFromString(*bonderOpt)
			if err != nil {
				cmd.PrintErrorMsg("Bonder address is not valid: %v", err)
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

			//sign transaction
			if *keyFileOpt == "" {
				cmd.PrintWarnMsg("Please specify a key file to sign.")
				c.PrintHelp()
				return
			}
			if *authOpt == "" {
				auth = cmd.PromptPassphrase("Passphrase: ", false)
			} else {
				auth = *authOpt
			}

			//RPC
			seq = *seqOpt
			if *seqOpt == 0 {
				seq, err = util.GetSequence(*grpcOpt, bonder)
				if err != nil {
					cmd.PrintErrorMsg("Couldn't retrieve Sequence number from RPC Server: %v", err)
					return
				}
			}

			stamp, err = crypto.HashFromString(*stampOpt)
			if err != nil {
				stamp, err = util.GetStamp(*grpcOpt)
				if err != nil {
					cmd.PrintErrorMsg("Couldn't retrieve stamp from RPC Server: %v", err)
					return
				}
			}
			//fulfill transaction payload
			trx := tx.NewBondTx(stamp, seq, bonder, pub, stake, fee, *memoOpt)

			signAndPublish(trx, *keyFileOpt, auth, *grpcOpt)
		}
	}
}
