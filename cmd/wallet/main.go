package main

import (
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

var (
	pathOpt        *string
	offlineOpt     *bool
	serverAddrsOpt *[]string
	timeoutOpt     *int
)

func addPasswordOption(c *cobra.Command) *string {
	return c.Flags().StringP("password", "p",
		"", "the wallet password")
}

func openWallet() (*wallet.Wallet, error) {
	opts := make([]wallet.Option, 0)

	if *serverAddrsOpt != nil {
		opts = append(opts, wallet.WithCustomServers(*serverAddrsOpt))
	}

	if *timeoutOpt > 0 {
		opts = append(opts, wallet.WithTimeout(time.Duration(*timeoutOpt)*time.Second))
	}

	wlt, err := wallet.Open(*pathOpt, *offlineOpt, opts...)
	if err != nil {
		return nil, err
	}

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
	serverAddrsOpt = rootCmd.PersistentFlags().StringSlice("servers", []string{}, "servers gRPC address")
	timeoutOpt = rootCmd.PersistentFlags().Int("timeout", 1,
		"specifies the timeout duration for the client connection in seconds")

	buildCreateCmd(rootCmd)
	buildRecoverCmd(rootCmd)
	buildGetSeedCmd(rootCmd)
	buildFeeCmd(rootCmd)
	buildChangePasswordCmd(rootCmd)
	buildAllTransactionCmd(rootCmd)
	buildAllAddrCmd(rootCmd)
	buildAllHistoryCmd(rootCmd)
	buildInfoCmd(rootCmd)
	buildNeuterCmd(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		cmd.PrintErrorMsgf("%s", err)
	}
}
