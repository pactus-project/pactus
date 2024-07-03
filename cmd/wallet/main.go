package main

import (
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

var (
	pathOpt       *string
	offlineOpt    *bool
	serverAddrOpt *string
	timeoutOpt    *int
)

func addPasswordOption(c *cobra.Command) *string {
	return c.Flags().StringP("password", "p",
		"", "the wallet password")
}

func openWallet() (*wallet.Wallet, error) {
	wlt, err := wallet.Open(*pathOpt, *offlineOpt)
	if err != nil {
		return nil, err
	}

	if *serverAddrOpt != "" {
		wlt.SetServerAddr(*serverAddrOpt)
	}

	wlt.SetClientTimeout(time.Duration(*timeoutOpt) * time.Second)

	return wlt, err
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
	timeoutOpt = rootCmd.PersistentFlags().Int("timeout", 1,
		"specifies the timeout duration for the client connection in seconds. "+
			"this timeout determines how long the client will attempt to establish a connection with the server "+
			"before failing with an error. Adjust this value based on network conditions and server response times "+
			"to ensure a balance between responsiveness and reliability.")

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
