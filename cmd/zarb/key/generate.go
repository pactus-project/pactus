package key

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/utils"
)

// Generate creates a new account and stores the keyfile in the disk
func Generate() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			keyObj := key.GenKey()
			passphrase := cmd.PromptPassphrase("Passphrase: ", true)
			label := cmd.PromptInput("Label: ")
			keyjson, err := key.EncryptKey(keyObj, passphrase, label)
			if err != nil {
				cmd.PrintErrorMsg("Failed to encrypt: %v", err)
				return
			}
			keyfilepath := cmd.ZarbKeystoreDir() + keyObj.Address().String() + ".json"

			// Store the file to disk.
			if err := utils.WriteFile(keyfilepath, keyjson); err != nil {
				cmd.PrintErrorMsg("Failed to write the key file: %v", err)
				return
			}

			fmt.Println()
			cmd.PrintInfoMsg("Key path: %v", keyfilepath)
			cmd.PrintInfoMsg("Address: %v", keyObj.Address())
			cmd.PrintInfoMsg("Public key: %v", keyObj.PublicKey())
		}
	}
}
