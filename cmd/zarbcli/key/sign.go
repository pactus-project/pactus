package key

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/tx"
)

// Sign the message with the private key and returns the signature
func Sign() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		messageFileOpt := c.String(cli.StringOpt{
			Name: "f file",
			Desc: "A file path to sign its content",
		})
		messageOpt := c.String(cli.StringOpt{
			Name: "m message",
			Desc: "Text message to sign",
		})
		transactionOpt := c.String(cli.StringOpt{
			Name: "t tx",
			Desc: "Raw transaction to sign",
		})
		keyFileOpt := c.String(cli.StringOpt{
			Name: "k keyfile",
			Desc: "Path to the encrypted key file",
		})
		authOpt := c.String(cli.StringOpt{
			Name: "a auth",
			Desc: "Passphrase of the key file",
		})

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			var trx *tx.Tx
			var msg []byte
			var err error
			//extract the message to be signed
			if *messageOpt != "" {
				msg = []byte(*messageOpt)
			} else if *messageFileOpt != "" {
				msg, err = ioutil.ReadFile(*messageFileOpt)
				if err != nil {
					cmd.PrintErrorMsg("Failed to read the file: %v", err)
					return
				}
			} else if *transactionOpt != "" {
				bz, err := hex.DecodeString(*transactionOpt)
				if err != nil {
					cmd.PrintErrorMsg("Invalid input: %v", err)
					return
				}
				trx = new(tx.Tx)
				err = trx.Decode(bz)
				if err != nil {
					cmd.PrintErrorMsg("Invalid transaction: %v", err)
					return
				}
			} else {
				cmd.PrintWarnMsg("Please specify a message or transaction to sign.")
				c.PrintHelp()
				return
			}

			//Sign the message with the private key
			if *keyFileOpt == "" {
				cmd.PrintWarnMsg("Please specify a key file to sign.")
				c.PrintHelp()
				return
			}
			var auth string
			if *authOpt == "" {
				auth = cmd.PromptPassphrase("Passphrase: ", false)
			} else {
				auth = *authOpt
			}

			key, err := key.DecryptKeyFile(*keyFileOpt, auth)
			if err != nil {
				cmd.PrintErrorMsg("Failed to decrypt: %v", err)
				return
			}

			if trx != nil {
				key.ToSigner().SignMsg(trx)
				bz, _ := trx.Encode()

				fmt.Println()
				cmd.PrintInfoMsg("Signed raw transaction:\n%x", bz)
			} else {
				signature := key.ToSigner().SignData(msg)

				fmt.Println()
				cmd.PrintInfoMsg("Signature: %s", signature)
			}
		}
	}
}
