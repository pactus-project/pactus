package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

var (
	pathOpt       *string
	offlineOpt    *bool
	serverAddrOpt *string
)

func addPasswordOption(c *cobra.Command) *string {
	return c.Flags().StringP("password", "p",
		"", "the wallet password")
}

func openWallet() (*wallet.Wallet, error) {
	if !*offlineOpt && *serverAddrOpt != "" {
		wlt, err := wallet.Open(*pathOpt, true)
		if err != nil {
			return nil, err
		}

		err = wlt.Connect(*serverAddrOpt)
		if err != nil {
			return nil, err
		}

		return wlt, err
	}
	wlt, err := wallet.Open(*pathOpt, *offlineOpt)
	if err != nil {
		return nil, err
	}

	return wlt, nil
}

func main() {
	rootCmd := &cobra.Command{
		Use:               "pactus-wallet",
		Short:             "Pactus wallet",
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	// Hide the "help" sub-command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	pathOpt = rootCmd.PersistentFlags().String("path",
		cmd.PactusDefaultWalletPath(cmd.PactusDefaultHomeDir()), "the path to the wallet file")
	offlineOpt = rootCmd.PersistentFlags().Bool("offline", false, "offline mode")
	serverAddrOpt = rootCmd.PersistentFlags().String("server", "", "server gRPC address")

	buildCreateCmd(rootCmd)
	buildRecoverCmd(rootCmd)
	buildGetSeedCmd(rootCmd)
	buildChangePasswordCmd(rootCmd)
	buildAllTransactionCmd(rootCmd)
	buildAllAddrCmd(rootCmd)
	buildAllHistoryCmd(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		cmd.PrintErrorMsgf("%s", err)
	}
}
