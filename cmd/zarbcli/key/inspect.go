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
		keyFileArg := c.String(cli.StringArg{
			Name: "KEYFILE",
			Desc: "Path to the encrypted key file",
		})
		showPrivateOpt := c.Bool(cli.BoolOpt{
			Name: "e expose-private-key",
			Desc: "expose the private key in the output",
		})
		authOpt := c.String(cli.StringOpt{
			Name: "a auth",
			Desc: "Passphrase of the key file",
		})
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			ek, err := key.NewEncryptedKey(*keyFileArg)
			if err != nil {
				cmd.PrintErrorMsg("Failed to read the key file: %v", err)
				return
			}
			var auth string
			if *authOpt == "" {
				auth = cmd.PromptPassphrase("Passphrase: ", false)
			} else {
				auth = *authOpt
			}
			// Decrypt key with passphrase.
			keyObj, err := ek.Decrypt(auth)
			if err != nil {
				cmd.PrintErrorMsg("Failed to decrypt: %v", err)
				return
			}

			cmd.PrintLine()
			cmd.PrintInfoMsg("Label: %v", ek.Label)
			cmd.PrintInfoMsg("Address: %v", keyObj.Address())
			cmd.PrintInfoMsg("Public key: %v", keyObj.PublicKey())
			if *showPrivateOpt {
				cmd.PrintDangerMsg("Private key: %v", keyObj.PrivateKey())
			}
		}
	}
}
