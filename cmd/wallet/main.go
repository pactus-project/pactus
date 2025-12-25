package main

import (
	"context"
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util/terminal"
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
	opts := make([]wallet.OpenWalletOption, 0)

	if *serverAddrsOpt != nil {
		opts = append(opts, wallet.WithCustomServers(*serverAddrsOpt))
	}

	if *timeoutOpt > 0 {
		opts = append(opts, wallet.WithTimeout(time.Duration(*timeoutOpt)*time.Second))
	}

	if *offlineOpt {
		opts = append(opts, wallet.WithOfflineMode())
	}

	wlt, err := wallet.Open(context.Background(), *pathOpt, opts...)
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
		"specifies the timeout duration for the connection in seconds")

	buildCreateCmd(rootCmd)
	buildRecoverCmd(rootCmd)
	buildGetSeedCmd(rootCmd)
	buildFeeCmd(rootCmd)
	buildPasswordCmd(rootCmd)
	buildSendCmd(rootCmd)
	buildAddressCmd(rootCmd)
	buildTransactionsCmd(rootCmd)
	buildInfoCmd(rootCmd)
	buildNeuterCmd(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		terminal.PrintErrorMsgf(err.Error())
	}
}
