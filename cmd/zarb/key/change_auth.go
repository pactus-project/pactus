package key

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
)

// ChangeAuth changes the passphrase of the key file
func ChangeAuth() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		keyFileArg := c.String(cli.StringArg{
			Name: "KEYFILE",
			Desc: "Path to the encrypted key file",
		})
		authOpt := c.String(cli.StringOpt{
			Name: "a auth",
			Desc: "Passphrase of the key file",
		})

		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			if *keyFileArg == "" {
				cmd.PrintWarnMsg("Key file is not specified.")
				c.PrintHelp()
				return
			}
			path := *keyFileArg
			// Decrypt key with passphrase.
			var oldAuth string
			if *authOpt == "" {
				oldAuth = cmd.PromptPassphrase("Passphrase: ", false)
			} else {
				oldAuth = *authOpt
			}
			ek, err := key.NewEncryptedKey(path)
			if err != nil {
				cmd.PrintErrorMsg("Failed to read the key: %v", err)
				return
			}
			keyObj, err := ek.Decrypt(oldAuth)
			if err != nil {
				cmd.PrintErrorMsg("Failed to decrypt: %v", err)
				return
			}
			//Prompt for the new passphrase
			newAuth := cmd.PromptPassphrase("New passphrase: ", true)
			//Prompt for the label
			label := cmd.PromptInputWithSuggestion("New label: ", ek.Label)
			// Encrypt key with passphrase.
			err = key.EncryptKeyToFile(keyObj, path, newAuth, label)
			if err != nil {
				cmd.PrintErrorMsg("Failed to encrypt: %v", err)
				return
			}

			fmt.Println()
			cmd.PrintSuccessMsg("Password changed successfully")
		}
	}
}
