package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"path/filepath"
	"strings"

	cli "github.com/jawher/mow.cli"
	"github.com/zarbchain/zarb-go/cmd"
	"github.com/zarbchain/zarb-go/node"
	"github.com/zarbchain/zarb-go/node/config"
	"github.com/zarbchain/zarb-go/types/crypto"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/genesis"
	"github.com/zarbchain/zarb-go/util"
	"github.com/zarbchain/zarb-go/version"
	"github.com/zarbchain/zarb-go/wallet"
)

//Start starts the zarb node
func Start() func(c *cli.Cmd) {
	return func(c *cli.Cmd) {
		workingDirOpt := c.String(cli.StringOpt{
			Name:  "w working-dir",
			Desc:  "Working directory to read node configuration and genesis files",
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
			workingDir, _ := filepath.Abs(*workingDirOpt)
			// change working directory
			if err := os.Chdir(workingDir); err != nil {
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

			gen, err := genesis.LoadFromFile(cmd.ZarbGenesisPath(workingDir))
			if err != nil {
				cmd.PrintErrorMsg("Aborted! Could not obtain genesis. %v", err)
				return
			}

			if gen.Params().IsTestnet() {
				crypto.DefaultHRP = "tzc"
			}

			conf, err := config.LoadFromFile(cmd.ZarbConfigPath(workingDir))
			if err != nil {
				cmd.PrintErrorMsg("Aborted! Could not obtain config. %v", err)
				return
			}

			if err = conf.SanityCheck(); err != nil {
				cmd.PrintErrorMsg("Aborted! Config is invalid. %v", err)
				return
			}

			signer, err := makeSigner(workingDir, keyFileOpt, privateKeyOpt)
			if err != nil {
				cmd.PrintErrorMsg("Aborted! %v", err)
				return
			}

			validatorAddr := signer.Address()
			rewardAddr := conf.State.RewardAddress
			if rewardAddr == "" {
				rewardAddr = validatorAddr.String()
			}
			cmd.PrintInfoMsg("You are running a zarb block chain version: %s. Welcome! ", version.Version())
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

// makeSigner makes a signer object from the validator private key.
// Private key obtains in this order:
// 1- From key file option (--key-file <path>)
// 2- From private key option (--private-key <secret>)
// 3- From 'validator_key' file insied the working directory
// 4- From the first address of the default_wallet
func makeSigner(workingDir string, keyFileOpt, privateKeyOpt *string) (crypto.Signer, error) {
	prvStr := ""
	switch {
	case *keyFileOpt == "" && *privateKeyOpt == "":
		keyPath := workingDir + "/validator_key"
		walletPath := cmd.ZarbDefaultWalletPath(workingDir)
		if util.PathExists(keyPath) {
			data, err := util.ReadFile(keyPath)
			if err != nil {
				return nil, err
			}
			prvStr = strings.TrimSpace(string(data))
		} else if util.PathExists(walletPath) {
			wallet, err := wallet.OpenWallet(walletPath)
			if err != nil {
				return nil, err
			}
			addrInfos := wallet.AddressInfos()
			if len(addrInfos) == 0 {
				return nil, fmt.Errorf("validator address is not defined")
			}
			password := ""
			if wallet.IsEncrypted() {
				password = cmd.PromptPassword("Wallet password", false)
			}
			valPrvKeyStr, err := wallet.PrivateKey(password, addrInfos[0].Address)
			if err != nil {
				return nil, err
			}
			prvStr = valPrvKeyStr
		} else {
			// Creating KeyObject from Private Key
			prvStr = cmd.PromptInput("Please enter the validator private key")
		}

	case *keyFileOpt != "":
		// Creating KeyObject from keystore
		data, err := util.ReadFile(*keyFileOpt)
		if err != nil {
			return nil, err
		}
		prvStr = strings.TrimSpace(string(data))

	case *privateKeyOpt != "":
		prvStr = *privateKeyOpt
	}

	prv, err := bls.PrivateKeyFromString(prvStr)
	if err != nil {
		return nil, err
	}

	return crypto.NewSigner(prv), nil
}
