package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zarbchain/zarb-go/validator"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/genesis"
	"github.com/zarbchain/zarb-go/keystore/key"
	"github.com/zarbchain/zarb-go/node"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
)

//Start starts the zarb node
func Start() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {

		workingDir := c.String(cli.StringOpt{
			Name:  "w working-dir",
			Desc:  "Working directory of the configuration and genesis files",
			Value: ".",
		})
		privateKey := c.String(cli.StringOpt{
			Name: "p privatekey",
			Desc: "Private key of the node's validator",
		})
		keyFile := c.String(cli.StringOpt{
			Name: "k keyfile",
			Desc: "Path to the encrypted key file contains validator's private key",
		})
		keyFileAuth := c.String(cli.StringOpt{
			Name: "a auth",
			Desc: "Key file's passphrase",
		})

		c.Spec = "[-w=<working directory>] [-p=<validator's private key>] | [-k=<path to the key file>] [-a=<key file's password>]"
		c.LongDesc = "Starting the node"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			path, _ := filepath.Abs(*workingDir)
			var keyObj *key.Key
			switch {
			case *keyFile == "" && *privateKey == "":
				f := path + "/validator_key.json"
				if util.PathExists(f) {
					kj, err := key.DecryptKeyFile(f, "")
					if err != nil {
						cmd.PrintErrorMsg("Aborted! %v", err)
						return
					}
					keyObj = kj
				} else {
					// Creating KeyObject from Private Key
					kj, err := cmd.PromptPrivateKey("Please enter the privateKey for the validator: ")
					if err != nil {
						cmd.PrintErrorMsg("Aborted! %v", err)
						return
					}
					keyObj = kj
				}
			case *keyFile != "" && *keyFileAuth != "":
				//Creating KeyObject from keystore
				passphrase := *keyFileAuth
				kj, err := key.DecryptKeyFile(*keyFile, passphrase)
				if err != nil {
					cmd.PrintErrorMsg("Aborted! %v", err)
					return
				}
				keyObj = kj
			case *keyFile != "" && *keyFileAuth == "":
				//Creating KeyObject from keystore
				passphrase := cmd.PromptPassphrase("Passphrase: ", false)
				kj, err := key.DecryptKeyFile(*keyFile, passphrase)
				if err != nil {
					cmd.PrintErrorMsg("Aborted! %v", err)
					return
				}
				keyObj = kj
			case *privateKey != "":
				// Creating KeyObject from Private Key
				pv, err := crypto.PrivateKeyFromString(*privateKey)
				if err != nil {
					cmd.PrintErrorMsg("Aborted! %v", err)
					return
				}
				keyObj, _ = key.NewKey(pv.PublicKey().Address(), pv)
			}

			cmd.PrintInfoMsg("Validator address: %v", keyObj.Address())

			// change working directory
			if err := os.Chdir(path); err != nil {
				cmd.PrintErrorMsg("Unable to changes working directory. %v", err)
				return
			}
			configFile := "./config.toml"
			genesisFile := "./genesis.json"

			gen, err := genesis.LoadFromFile(genesisFile)
			if err != nil {
				cmd.PrintErrorMsg("Could not obtain genesis. %v", err)
				return
			}

			conf, err := config.LoadFromFile(configFile)
			if err != nil {
				cmd.PrintErrorMsg("Could not obtain config. %v", err)
				return
			}

			err = conf.Check()
			if err != nil {
				cmd.PrintErrorMsg("Config is invalid - %v", err)
				return
			}

			cmd.PrintInfoMsg("You are running a zarb block chain node version: %v. Welcome! ", version.NodeVersion.String())

			privVal := validator.NewPrivValidator(keyObj.PrivateKey())
			node, err := node.NewNode(gen, conf, privVal)
			if err != nil {
				cmd.PrintErrorMsg("Could not create node. %v", err)
				return
			}

			if err := node.Start(); err != nil {
				cmd.PrintErrorMsg("Could not start node. %v", err)
				return
			}

			cmd.TrapSignal(func() {
				node.Stop()
				cmd.PrintInfoMsg("Exiting ...")
			})

			// run forever (the node will not be returned)
			select {}
		}
	}
}
