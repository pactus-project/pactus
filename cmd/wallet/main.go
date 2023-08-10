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

	// transaction commands
	buildTransactionCmd(rootCmd)
	buildTransferTxCmd(txCmd)
	buildBondTxCmd(txCmd)
	buildUnbondTxCmd(txCmd)
	buildWithdrawTxCmd(txCmd)

	// address commands
	buildAddrCmd(rootCmd)
	buildAllAddressesCmd(addrCmd)
	buildNewAddressCmd(addrCmd)
	buildBalanceCmd(addrCmd)
	buildPrivateKeyCmd(addrCmd)
	buildPublicKeyCmd(addrCmd)
	buildImportPrivateKeyCmd(addrCmd)
	buildSetLabelCmd(addrCmd)

	// history commands
	buildHistoryCmd(rootCmd)
	buildAddToHistoryCmd(historyCmd)
	buildShowHistoryCmd(historyCmd)

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
