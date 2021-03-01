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
		deadlockOpt := c.Bool(cli.BoolOpt{
			Name:  "disable-deadlock",
			Desc:  "Disable deadlock detection mode",
			Value: false,
		})

		c.LongDesc = "Starting the node from working directory"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			configFile := "./config.toml"
			genesisFile := "./genesis.json"
			var err error
			var keyObj *key.Key
			var workspace string

			if *deadlockOpt {
				// Disable dead-lock detection
				deadlock.Opts.Disable = true
			}

			workspace = *workingDirOpt
			if workspace == "." {
				if !util.PathExists(genesisFile) {
					cmd.PrintErrorMsg("Aborted! No genesis file")
					return
				}
			}

			workspace, err = filepath.Abs(workspace)
			if err != nil {
				cmd.PrintErrorMsg("Aborted! %v", err)
				return
			}

			// change working directory
			if err := os.Chdir(workspace); err != nil {
				cmd.PrintErrorMsg("Aborted! Unable to changes working directory. %v", err)
				return
			}

			gen, err := genesis.LoadFromFile(genesisFile)
			if err != nil {
				cmd.PrintErrorMsg("Aborted! Could not obtain genesis. %v", err)
				return
			}

			conf, err := config.LoadFromFile(configFile)
			if err != nil {
				cmd.PrintErrorMsg("Aborted! Could not obtain config. %v", err)
				return
			}

			if err = conf.SanityCheck(); err != nil {
				cmd.PrintErrorMsg("Aborted! Config is invalid. %v", err)
				return
			}

			keyObj, err = retrievePrivateKey(workspace, keyFileOpt, authOpt, privateKeyOpt)
			if err != nil {
				cmd.PrintErrorMsg("Aborted! %v", err)
				return
			}

			validatorAddr := keyObj.Address()
			mintbaseAddr := conf.State.MintbaseAddress
			if mintbaseAddr == "" {
				mintbaseAddr = validatorAddr.String()
			}
			cmd.PrintInfoMsg("You are running a zarb block chain node version: %v. Welcome! ", version.NodeVersion.String())
			cmd.PrintInfoMsg("Validator address: %v", validatorAddr)
			cmd.PrintInfoMsg("Mintbase address : %v", mintbaseAddr)
			cmd.PrintLine()

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

func retrievePrivateKey(workspace string, keyFileOpt, authOpt, privateKeyOpt *string) (*key.Key, error) {

	switch {
	case *keyFileOpt == "" && *privateKeyOpt == "":
		f := workspace + "/validator_key.json"
		if util.PathExists(f) {
			kj, err := key.DecryptKeyFile(f, "")
			if err != nil {
				return nil, err
			}
			return kj, nil
		}
		// Creating KeyObject from Private Key
		kj, err := cmd.PromptPrivateKey("Please enter the privateKey for the validator: ")
		if err != nil {
			return nil, err
		}
		return kj, nil

	case *keyFileOpt != "" && *authOpt != "":
		// Creating KeyObject from keystore
		auth := *authOpt
		kj, err := key.DecryptKeyFile(*keyFileOpt, auth)
		if err != nil {
			return nil, err
		}
		return kj, nil
	case *keyFileOpt != "" && *authOpt == "":
		// Creating KeyObject from keystore
		auth := cmd.PromptPassphrase("Passphrase: ", false)
		kj, err := key.DecryptKeyFile(*keyFileOpt, auth)
		if err != nil {
			return nil, err
		}
		return kj, nil
	case *privateKeyOpt != "":
		// Creating KeyObject from Private Key
		pv, err := crypto.PrivateKeyFromString(*privateKeyOpt)
		if err != nil {
			return nil, err
		}
		return key.NewKey(pv.PublicKey().Address(), pv)
	}

	return nil, fmt.Errorf("Invalid input")
}
