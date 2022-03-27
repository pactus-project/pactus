package tx

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	"github.com/zarbchain/zarb-go/tx"
	grpcclient "github.com/zarbchain/zarb-go/www/grpc/client"
)

func SendTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		stampOpt := c.String(cli.StringOpt{
			Name: "stamp",
			Desc: "Transaction stamp if not specified will query from RPC server",
		})

		seqOpt := c.Int(cli.IntOpt{
			Name: "seq",
			Desc: "Transaction sequence number if not specified will query from RPC server",
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
			Desc: "gRPC server address if not specified will just print raw signed transaction",
		})
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			var err error
			var stamp hash.Stamp
			var sender crypto.Address
			var receiver crypto.Address
			var seq int32
			var amount int64
			var fee int64
			var auth string

			// ---
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
			if seqOpt != nil {
				seq = int32(*seqOpt)
			} else {
				seq, err = grpcclient.GetSequence(promptRPCEndpoint(grpcOpt), sender)
				if err != nil {
					cmd.PrintErrorMsg("Couldn't retrieve sequence number from RPC Server: %v", err)
					return
				}
			}

			if stampOpt == nil || *stampOpt == "" {
				stamp, err = grpcclient.GetStamp(promptRPCEndpoint(grpcOpt))
				if err != nil {
					cmd.PrintErrorMsg("Couldn't retrieve stamp from RPC Server: %v", err)
					return
				}
			} else {
				stamp, err = hash.StampFromString(*stampOpt)
				if err != nil {
					cmd.PrintErrorMsg("Couldn't decode stamp from input: %v", err)
					return
				}
			}

			//fulfill transaction payload
			trx := tx.NewSendTx(stamp, seq, sender, receiver, amount, fee, *memoOpt)

			//sign transaction
			signAndPublish(trx, *keyFileOpt, auth, grpcOpt)
		}
	}
}
