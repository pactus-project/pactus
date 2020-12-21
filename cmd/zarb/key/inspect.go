package key

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
)

//Inspect displays various information of the keyfile
func Inspect() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		keyFile := c.String(cli.StringArg{
			Name: "KEYFILE",
			Desc: "Path to the encrypted key file",
		})
		showPrivate := c.Bool(cli.BoolOpt{
			Name: "e expose-private-key",
			Desc: "expose the private key in the output",
		})
		c.Spec = "KEYFILE [-e]"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			ek, err := key.NewEncryptedKey(*keyFile)
			if err != nil {
				cmd.PrintErrorMsg("Failed to read the key file: %v", err)
				return
			}
			// Decrypt key with passphrase.
			passphrase := cmd.PromptPassphrase("Passphrase: ", false)
			keyObj, err := ek.Decrypt(passphrase)
			if err != nil {
				cmd.PrintErrorMsg("Failed to decrypt: %v", err)
				return
			}

			fmt.Println()
			cmd.PrintInfoMsg("Label: %v", ek.Label)
			cmd.PrintInfoMsg("Address: %v", keyObj.Address())
			cmd.PrintInfoMsg("Public key: %v", keyObj.PublicKey())
			if *showPrivate {
				cmd.PrintInfoMsg("Private key: %v", keyObj.PrivateKey())
			}
		}
	}
}
