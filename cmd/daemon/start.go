package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"path/filepath"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/config"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/bls"
	"github.com/zarbchain/zarb-go/genesis"
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
			Desc: "Path to the key file contains validator's private key",
		})
		pprofOpt := c.String(cli.StringOpt{
			Name: "pprof",
			Desc: "debug pprof server address(not recommended to expose to internet)",
		})

		c.LongDesc = "Starting the node from working directory"
		c.Before = func() { fmt.Println(cmd.ZARB) }
		c.Action = func() {
			configFile := "./config.toml"
			genesisFile := "./genesis.json"
			var err error

			workspace, _ := filepath.Abs(*workingDirOpt)

			signer, err := makeSigner(workspace, keyFileOpt, privateKeyOpt)
			if err != nil {
				cmd.PrintErrorMsg("Aborted! %v", err)
				return
			}

			// change working directory
			if err := os.Chdir(workspace); err != nil {
				cmd.PrintErrorMsg("Aborted! Unable to changes working directory. %v", err)
				return
			}

			// separate pprof handlers from DefaultServeMux.
			pprofMux := http.DefaultServeMux
			http.DefaultServeMux = http.NewServeMux()
			if *pprofOpt != "" {
				cmd.PrintWarnMsg("Starting Debug pprof server on: %v", *pprofOpt)
				go func() {
					err := http.ListenAndServe(*pprofOpt, pprofMux)
					if err != nil {
						cmd.PrintErrorMsg("Could not initialize pprof server. %v", err)
					}
				}()
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

			validatorAddr := signer.Address()
			rewardAddr := conf.State.RewardAddress
			if rewardAddr == "" {
				rewardAddr = validatorAddr.String()
			}
			cmd.PrintInfoMsg("You are running a zarb block chain agent: %v. Welcome! ", version.Version())
			cmd.PrintInfoMsg("Validator address: %v", validatorAddr)
			cmd.PrintInfoMsg("Reward address : %v", rewardAddr)
			cmd.PrintLine()

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

func makeSigner(workspace string, keyFileOpt, privateKeyOpt *string) (crypto.Signer, error) {
	prvHex := ""
	switch {
	case *keyFileOpt == "" && *privateKeyOpt == "":
		path := workspace + "/validator_key"
		if util.PathExists(path) {
			data, err := util.ReadFile(path)
			if err != nil {
				return nil, err
			}
			prvHex = string(data)
		} else {
			// Creating KeyObject from Private Key
			prvHex = cmd.PromptInput("Please enter the private key in hex format: ")
		}

	case *keyFileOpt != "":
		// Creating KeyObject from keystore
		data, err := util.ReadFile(*keyFileOpt)
		if err != nil {
			return nil, err
		}
		prvHex = string(data)
	case *privateKeyOpt != "":
		prvHex = *privateKeyOpt
	}

	prv, err := bls.PrivateKeyFromString(prvHex)
	if err != nil {
		return nil, err
	}

	return crypto.NewSigner(prv), nil
}
