package key

import (
	"fmt"

	cli "github.com/jawher/mow.cli"
	"github.com/tyler-smith/go-bip39"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/keystore/key"
)

// Recover tries to recover a key from the seed
func Recover() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			mnemonic := cmd.PromptInput("Seed: ")
			passphrase := cmd.PromptInput("Passphrase: ")

			seed, err := bip39.NewSeedWithErrorChecking(mnemonic, passphrase)
			if err != nil {
				cmd.PrintErrorMsg("Seed is not correct: %v", err)
				return
			}
			keyObj, err := key.FromSeed(seed)
			if err != nil {
				cmd.PrintErrorMsg("Failed to create key from the seed: %v", err)
				return
			}

			keyfilepath := cmd.ZarbKeystoreDir() + keyObj.Address().String() + "_recovered.json"
			err = key.EncryptKeyToFile(keyObj, keyfilepath, passphrase, "recovered")
			if err != nil {
				cmd.PrintErrorMsg("Failed to encrypt: %v", err)
				return
			}

			fmt.Println()
			cmd.PrintInfoMsg("Key path: %v", keyfilepath)
			cmd.PrintInfoMsg("Address: %v", keyObj.Address())
		}
	}
}
