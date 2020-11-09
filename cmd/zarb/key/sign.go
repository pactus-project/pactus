package key

import (
	"fmt"
	"io/ioutil"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/keystore/key"
)

// Sign the message with the private key and returns the signature hash
func Sign() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		messageFile := c.String(cli.StringOpt{
			Name: "f file",
			Desc: "Message file path to read the file and sign the message inside",
		})
		message := c.String(cli.StringOpt{
			Name: "m message",
			Desc: "Text message to sign",
		})
		privateKey := c.String(cli.StringOpt{
			Name: "p private-Key",
			Desc: "Private key to sign the message",
		})
		keyFile := c.String(cli.StringOpt{
			Name: "k keyfile",
			Desc: "Path to the encrypted key file",
		})
		keyFileAuth := c.String(cli.StringOpt{
			Name: "a auth",
			Desc: "Key file's passphrase",
		})

		c.Spec = "[-f=<message file>] | [-m=<message to sign>]" +
			" [-p=<private key>] | [-k=<path to the key file>] [-a=<key file's passphrase>]"
		c.LongDesc = "Signing a message "
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			var msg []byte
			var err error
			//extract the message to be signed
			if *message != "" {
				msg = []byte(*message)
			} else if *messageFile != "" {
				msg, err = ioutil.ReadFile(*messageFile)
				if err != nil {
					cmd.PrintErrorMsg("Failed to read the file: %v", err)
					return
				}
			} else {
				cmd.PrintWarnMsg("Please enter a message to sign.")
				c.PrintHelp()
				return
			}

			var signature *crypto.Signature
			var pv crypto.PrivateKey
			//Sign the message with the private key
			if *privateKey != "" {
				pv, err = crypto.PrivateKeyFromString(*privateKey)
				if err != nil {
					cmd.PrintErrorMsg("%v", err)
					return
				}
				signature = pv.Sign(msg)
			} else if *keyFile != "" {
				var passphrase string
				if *keyFileAuth == "" {
					passphrase = cmd.PromptPassphrase("Passphrase: ", false)
				} else {
					passphrase = *keyFileAuth
				}

				kj, err := key.DecryptKeyFile(*keyFile, passphrase)
				if err != nil {
					cmd.PrintErrorMsg("Failed to decrypt: %v", err)
					return
				}
				pv = kj.PrivateKey()
				signature = pv.Sign(msg)
			} else {
				cmd.PrintWarnMsg("Please specify a key file to sign.")
				c.PrintHelp()
				return
			}

			//display the signature
			fmt.Println()
			cmd.PrintInfoMsg("Signature: %v", signature)
		}
	}
}
