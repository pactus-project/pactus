package main

import (
	"context"
	"time"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pactus-project/pactus/wallet/provider/remote"
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

func openWallet(ctx context.Context) (*wallet.Wallet, error) {
	wlt, err := wallet.Open(ctx, *pathOpt)
	if err != nil {
		return nil, err
	}

	if !*offlineOpt {
		var opts []remote.RemoteProviderOption
		if *serverAddrsOpt != nil {
			opts = append(opts, remote.WithCustomServers(*serverAddrsOpt))
		}

		if *timeoutOpt > 0 {
			opts = append(opts, remote.WithTimeout(time.Duration(*timeoutOpt)*time.Second))
		}

		provider, err := remote.NewRemoteBlockchainProvider(ctx, wlt.Info().Network, opts...)
		if err != nil {
			return nil, err
		}

		wlt.SetProvider(provider)
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
	serverAddrsOpt = rootCmd.PersistentFlags().StringSlice("servers", nil, "servers gRPC address")
	timeoutOpt = rootCmd.PersistentFlags().Int("timeout", 1,
		"specifies the timeout duration for the connection in seconds")

	buildCreateCmd(rootCmd)
	buildRecoverCmd(rootCmd)
	buildGetSeedCmd(rootCmd)
	buildFeeCmd(rootCmd)
	buildPasswordCmd(rootCmd)
	buildSendCmd(rootCmd)
	buildAddressCmd(rootCmd)
	buildTransactionCmd(rootCmd)
	buildInfoCmd(rootCmd)
	buildNeuterCmd(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		terminal.PrintErrorMsgf(err.Error())
	}
}
