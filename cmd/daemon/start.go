package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof" // #nosec
	"os"
	"path/filepath"
	"time"

	cli "github.com/jawher/mow.cli"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/node"
	"github.com/pactus-project/pactus/node/config"
	"github.com/pactus-project/pactus/types/genesis"
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
				crypto.AddressHRP = "tpc"
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

			walletPath := cmd.PactusDefaultWalletPath(workingDir)
			wallet, err := wallet.OpenWallet(walletPath, true)
			if err != nil {
				cmd.PrintErrorMsg("Aborted! %v", err)
				return
			}
			addrInfos := wallet.AddressLabels()
			if len(addrInfos) == 0 {
				cmd.PrintErrorMsg("Aborted! %v", err)
				return
			}
			password := ""
			if wallet.IsEncrypted() {
				password = cmd.PromptPassword("Wallet password", false)
			}

			signers := make([]crypto.Signer, conf.NumValidators)
			rewardAddrs := make([]crypto.Address, conf.NumValidators)
			for i := 0; i < conf.NumValidators; i++ {
				prvKey, err := wallet.PrivateKey(password, addrInfos[i*2].Address)
				if err != nil {
					cmd.PrintErrorMsg("Aborted! %v", err)
					return
				}

				addr, err := crypto.AddressFromString(addrInfos[(i*2)+1].Address)
				if err != nil {
					cmd.PrintErrorMsg("Aborted! %v", err)
					return
				}

				signers[i] = crypto.NewSigner(prvKey)
				rewardAddrs[i] = addr

				cmd.PrintInfoMsg("Validator address %v: %s", i+1, addrInfos[i*2].Address)
				cmd.PrintInfoMsg("Reward    address %v: %s", i+1, addrInfos[(i*2)+1].Address)
			}

			cmd.PrintLine()

			node, err := node.NewNode(gen, conf, signers, rewardAddrs)
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
