package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"path/filepath"
	"strings"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/node/config"
	"github.com/pactus-project/pactus/types/genesis"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/version"
	"github.com/pactus-project/pactus/wallet"
)

// Start starts the pactus node.
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
		c.Before = func() { fmt.Println(cmd.Pactus) }
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
				server := &http.Server{
					Addr:              *pprofOpt,
					ReadHeaderTimeout: 3 * time.Second,
					Handler:           pprofMux,
				}
				go func() {
					err := server.ListenAndServe()
					if err != nil {
						cmd.PrintErrorMsg("Could not initialize pprof server. %v", err)
					}
				}()
			}

			gen, err := genesis.LoadFromFile(cmd.PactusGenesisPath(workingDir))
			if err != nil {
				cmd.PrintErrorMsg("Aborted! Could not obtain genesis. %v", err)
				return
			}

			if gen.Params().IsTestnet() {
				crypto.AddressHRP = "tzc"
				crypto.PublicKeyHRP = "tpublic"
				crypto.PrivateKeyHRP = "tsecret"
				crypto.XPublicKeyHRP = "txpublic"
				crypto.XPrivateKeyHRP = "txsecret"
			}

			conf, err := config.LoadFromFile(cmd.PactusConfigPath(workingDir))
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
			cmd.PrintInfoMsg("You are running a pactus block chain version: %s. Welcome! ", version.Version())
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
// If none of them, it asks user to enter the private key.
func makeSigner(workingDir string, keyFileOpt, privateKeyOpt *string) (crypto.Signer, error) {
	prvStr := ""
	switch {
	case *keyFileOpt == "" && *privateKeyOpt == "":
		keyPath := workingDir + "/validator_key"
		walletPath := cmd.PactusDefaultWalletPath(workingDir)
		if util.PathExists(keyPath) {
			data, err := util.ReadFile(keyPath)
			if err != nil {
				return nil, err
			}
			prvStr = strings.TrimSpace(string(data))
		} else if util.PathExists(walletPath) {
			valKeyWallet, err := getValidatorKeyFromWallet(walletPath)
			if err != nil {
				return nil, err
			}
			prvStr = valKeyWallet
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

func getValidatorKeyFromWallet(walletPath string) (string, error) {
	wallet, err := wallet.OpenWallet(walletPath, true)
	if err != nil {
		return "", err
	}
	addrInfos := wallet.AddressLabels()
	if len(addrInfos) == 0 {
		return "", fmt.Errorf("validator address is not defined")
	}
	password := ""
	if wallet.IsEncrypted() {
		password = cmd.PromptPassword("Wallet password", false)
	}
	prvKey, err := wallet.PrivateKey(password, addrInfos[0].Address)
	if err != nil {
		return "", err
	}

	return prvKey.String(), nil
}
