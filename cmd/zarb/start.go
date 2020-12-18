package main

import (
	"fmt"
	"os"
	"path/filepath"

	cli "github.com/jawher/mow.cli"
	"github.com/sasha-s/go-deadlock"
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

		workingDirOpt := c.String(cli.StringOpt{
			Name:  "w working-dir",
			Desc:  "Working directory of the configuration and genesis files",
			Value: ".",
		})
		privateKeyOpt := c.String(cli.StringOpt{
			Name: "p private-key",
			Desc: "Validator's private key",
		})
		keyFileOpt := c.String(cli.StringOpt{
			Name: "k key-file",
			Desc: "Path to the encrypted key file contains validator's private key",
		})
		authOpt := c.String(cli.StringOpt{
			Name: "a auth",
			Desc: "Passphrase of the key file",
		})
		wizardOpt := c.Bool(cli.BoolOpt{
			Name:  "wizard",
			Desc:  "Start new node in wizard mode",
			Value: false,
		})

		deadlockdOpt := c.Bool(cli.BoolOpt{
			Name:  "deadlock",
			Desc:  "Enable deadlock detection mode",
			Value: false,
		})

		c.Spec = "[-w=<path>] [-p=<private_key>] | ([-k=<path>] [-a=<passphrase>]) | [--wizard] | [--deadlock] "
		c.LongDesc = "Starting the node"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {

			if !*deadlockdOpt {
				// Disable dead-lock detection, Should we define a flag for this?
				deadlock.Opts.Disable = false
			}

			configFile := "./config.toml"
			genesisFile := "./genesis.json"
			var err error
			var keyObj *key.Key
			var workspace string

			if *wizardOpt {
				defaultWorkspace := os.Getenv("HOME") + "/.zarb"

				workspace = cmd.PromptInput(fmt.Sprintf("Enter the workspace path (%v): ", defaultWorkspace))
				if workspace == "" {
					workspace = defaultWorkspace
				}

				workspace, err = filepath.Abs(workspace)
				if err != nil {
					cmd.PrintErrorMsg("Aborted! %v", err)
					return
				}

				if !util.IsDirNotExistsOrEmpty(workspace) {
					cmd.PrintErrorMsg("Workspace is not empty. %v", workspace)
					return
				}

				gen := genesis.Testnet()
				conf := makeConfigfile()

				conf.Network.Name = "zarb-testnet"
				conf.Network.Bootstrap.Addresses = []string{"/ip4/47.254.199.97/tcp/35470/ipfs/12D3KooWJy3oZ1mZh4TbLZKLBzAGJnGwyrbo2mn8oyb4zu121uQD"}
				conf.Network.Bootstrap.MinPeerThreshold = 1

				// save genesis file to file system
				genFile := workspace + "/genesis.json"
				if err := gen.SaveToFile(genFile); err != nil {
					cmd.PrintErrorMsg("Failed to write genesis file: %v", err)
					return
				}

				// save config file to file system
				confFile := workspace + "/config.toml"
				if err := conf.SaveToFile(confFile); err != nil {
					cmd.PrintErrorMsg("Failed to write config file: %v", err)
					return
				}

				keyObj = key.GenKey()
				if err := key.EncryptKeyFile(keyObj, workspace+"/validator_key.json", "", ""); err != nil {
					cmd.PrintErrorMsg("Failed to write key file: %v", err)
					return
				}

			} else {

				workspace = *workingDirOpt
				if workspace == "." {
					if !util.PathExists(genesisFile) {
						c.PrintHelp()
						return
					}
				}

				workspace, err = filepath.Abs(workspace)
				if err != nil {
					cmd.PrintErrorMsg("Aborted! %v", err)
					return
				}
				switch {
				case *keyFileOpt == "" && *privateKeyOpt == "":
					f := workspace + "/validator_key.json"
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
				case *keyFileOpt != "" && *authOpt != "":
					//Creating KeyObject from keystore
					passphrase := *authOpt
					kj, err := key.DecryptKeyFile(*keyFileOpt, passphrase)
					if err != nil {
						cmd.PrintErrorMsg("Aborted! %v", err)
						return
					}
					keyObj = kj
				case *keyFileOpt != "" && *authOpt == "":
					//Creating KeyObject from keystore
					passphrase := cmd.PromptPassphrase("Passphrase: ", false)
					kj, err := key.DecryptKeyFile(*keyFileOpt, passphrase)
					if err != nil {
						cmd.PrintErrorMsg("Aborted! %v", err)
						return
					}
					keyObj = kj
				case *privateKeyOpt != "":
					// Creating KeyObject from Private Key
					pv, err := crypto.PrivateKeyFromString(*privateKeyOpt)
					if err != nil {
						cmd.PrintErrorMsg("Aborted! %v", err)
						return
					}
					keyObj, _ = key.NewKey(pv.PublicKey().Address(), pv)
				}
			}

			// change working directory
			if err := os.Chdir(workspace); err != nil {
				cmd.PrintErrorMsg("Unable to changes working directory. %v", err)
				return
			}

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

			validatorAddr := keyObj.Address()
			mintbaseAddr := conf.State.MintbaseAddress
			if mintbaseAddr == nil {
				mintbaseAddr = &validatorAddr
			}
			cmd.PrintInfoMsg("You are running a zarb block chain node version: %v. Welcome! ", version.NodeVersion.String())
			cmd.PrintInfoMsg("Validator address: %v", validatorAddr)
			cmd.PrintInfoMsg("Mintbase address : %v", mintbaseAddr)
			cmd.PrintInfoMsg("")

			signer := keyObj.ToSigner()
			node, err := node.NewNode(gen, conf, signer)
			if err != nil {
				cmd.PrintErrorMsg("Could not initialize node. %v", err)
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
