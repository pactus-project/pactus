package main

import (
	"fmt"

	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

var pathArg *string
var offlineOpt *bool
var serverAddrOpt *string

func addPasswordOption(c *cobra.Command) *string {
	return c.Flags().StringP("password", "p",
		"", "the wallet password")
}

func openWallet() (*wallet.Wallet, error) {
	if !*offlineOpt && *serverAddrOpt != "" {
		wallet, err := wallet.Open(*pathArg, true)
		if err != nil {
			return nil, err
		}

		err = wallet.Connect(*serverAddrOpt)
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return wallet, err
	}
	wallet, err := wallet.Open(*pathArg, *offlineOpt)
	if err != nil {
		return nil, err
	}
	return wallet, nil
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "pactus-wallet",
		Short: "Pactus wallet",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("use --help")
		},
	}

	pathArg = rootCmd.Flags().String("PATH", "", "the path to the wallet file")
	offlineOpt = rootCmd.Flags().Bool("offline", false, "offline mode")
	serverAddrOpt = rootCmd.Flags().String("server", "", "server gRPC address")

	buildGenerateCmd(rootCmd)
	buildRecoverCmd(rootCmd)
	buildGetSeedCmd(rootCmd)
	buildChangePasswordCmd(rootCmd)

	app.Command("address", "Manage address book", func(k *cli.Cmd) {
		k.Command("new", "Creating a new address", NewAddress())
		k.Command("all", "Show all addresses", AllAddresses())
		k.Command("label", "Set label for the an address", SetLabel())
		k.Command("balance", "Show the balance of an address", Balance())
		k.Command("pub", "Show the public key of an address", PublicKey())
		k.Command("priv", "Show the private key of an address", PrivateKey())
		k.Command("import", "Import a private key into wallet", ImportPrivateKey())
	})

	app.Command("tx", "Create, sign and publish a transaction", func(k *cli.Cmd) {
		k.Command("bond", "Create, sign and publish a Bond transaction", BondTx())
		k.Command("transfer", "Create, sign and publish a Transfer transaction", TransferTx())
		k.Command("unbond", "Create, sign and publish an Unbond transaction", UnbondTx())
		k.Command("withdraw", "Create, sign and publish a Withdraw transaction", WithdrawTx())
	})
	
	app.Command("history", "Check the wallet history", func(k *cli.Cmd) {
		k.Command("add", "Add a transaction to the wallet history", AddToHistory())
		k.Command("get", "Show the transaction history of any address", ShowHistory())
	})

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
