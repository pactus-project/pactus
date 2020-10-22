package key

import (
	"fmt"
	"io/ioutil"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/util"
)

// ChangeAuth changes the passphrase of the key file
func ChangeAuth() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		keyFile := c.String(cli.StringArg{
			Name: "KEYFILE",
			Desc: "Path to the encrypted key file",
		})

		c.Spec = "KEYFILE"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			if *keyFile == "" {
				cmd.PrintWarnMsg("Key file is not specified.")
				c.PrintHelp()
				return
			}
			//Read the key from the keyfile
			keyjson, err := ioutil.ReadFile(*keyFile)
			if err != nil {
				cmd.PrintErrorMsg("Failed to read the keyfile: %v", err)
				return
			}
			// Decrypt key with passphrase.
			passphrase := cmd.PromptPassphrase("Old passphrase: ", false)
			keyObj, err := key.DecryptKey(keyjson, passphrase)
			if err != nil {
				cmd.PrintErrorMsg("Failed to decrypt: %v", err)
				return
			}
			//Prompt for the new passphrase
			passphrase = cmd.PromptPassphrase("New passphrase: ", true)
			//Prompt for the label
			label := cmd.PromptInput("New label: ")
			// Encrypt key with passphrase.
			keyjson, err = key.EncryptKey(keyObj, passphrase, label)
			if err != nil {
				cmd.PrintErrorMsg("Failed to encrypt: %v", err)
				return
			}
			// Store the file to disk.
			if err := util.WriteFile(*keyFile, keyjson); err != nil {
				cmd.PrintErrorMsg("Failed to write the key file: %v", err)
				return
			}

			fmt.Println()
			cmd.PrintSuccessMsg("Password changed successfully")
		}
	}
}
