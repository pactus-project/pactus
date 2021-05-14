package tx

import (
	"fmt"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/tx"
)

func SendTx() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		stampOpt := c.String(cli.StringOpt{
			Name: "stamp",
			Desc: "Transaction stamp if not specified will get from RPC server",
		})

		seqOpt := c.Int(cli.IntOpt{
			Name: "seq",
			Desc: "Transaction sequence number if not specified will get from RPC server",
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
			Desc: "gRPC server address",
		})
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			var err error
			var stamp crypto.Hash
			var sender crypto.Address
			var receiver crypto.Address
			var rpc string
			var seq int
			var amount int64
			var fee int64
			var auth string

			// ---
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
			if *grpcOpt == "" {
				rpc = cmd.PromptInput("gRPC server address: ")
			} else {
				rpc = *grpcOpt
			}

			grpcClient, err := cmd.GetRPCClient(rpc)
			if err != nil {
				cmd.PrintErrorMsg("Couldn't connect to RPC Server: %v", err)
				return
			}

			seq = *seqOpt
			if *seqOpt == 0 {
				seq, err = cmd.GetSequence(grpcClient, sender)
				if err != nil {
					cmd.PrintErrorMsg("Couldn't retrieve Sequence number from RPC Server: %v", err)
					return
				}
			}

			stamp, err = crypto.HashFromString(*stampOpt)
			if err != nil {
				stamp, err = cmd.GetStamp(grpcClient)
				if err != nil {
					cmd.PrintErrorMsg("Couldn't retrieve stamp from RPC Server: %v", err)
					return
				}
			}

			//fulfill transaction payload
			trx := tx.NewSendTx(stamp, seq, sender, receiver, amount, fee, *memoOpt)

			//sign transaction
			k, err := key.DecryptKeyFile(*keyFileOpt, auth)
			if err != nil {
				cmd.PrintErrorMsg("Couldn't retrieve Key: %v", err)
				return
			}
			k.ToSigner().SignMsg(trx)

			cmd.PrintWarnMsg("you are about to publish:")
			cmd.PrintJSONObject(trx)
			confirm := cmd.PromptInput("press y/yes to continue:")
			if !strings.HasSuffix(strings.ToLower(confirm), "y") {
				cmd.PrintWarnMsg("Opration aborted!")
				return
			}

			// publish
			signedTrx, _ := trx.Encode()
			if id, err := cmd.SendTx(grpcClient, signedTrx); err != nil {
				cmd.PrintErrorMsg("Couldn't publish transaction: %v", err)
				return
			} else {
				cmd.PrintSuccessMsg("transaction sent with Id: %v", id)
			}
		}
	}
}
