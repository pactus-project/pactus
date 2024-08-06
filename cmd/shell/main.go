package main

import (
	"encoding/base64"
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/inancgumus/screen"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util/shell"
	pb "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

const (
	defaultServerAddr     = "localhost:50051"
	defaultResponseFormat = "prettyjson"
)

var _prefix string

func main() {
	var (
		serverAddr string
		username   string
		password   string
	)

	rootCmd := &cobra.Command{
		Use:   "shell",
		Short: "Pactus Shell",
		Long:  `pactus-shell is a command line tool for interacting with the Pactus blockchain using gRPC`,
	}

	sh := shell.New(rootCmd, nil,
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.Green),
		prompt.OptionDescriptionBGColor(prompt.Black),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionLivePrefix(livePrefix),
	)

	sh.Flags().StringVar(&serverAddr, "server-addr", defaultServerAddr, "gRPC server address")
	sh.Flags().StringVar(&username, "auth-username", "",
		"username for gRPC basic authentication")

	sh.Flags().StringVar(&password, "auth-password", "",
		"username for gRPC basic authentication")

	sh.PreRun = func(_ *cobra.Command, _ []string) {
		cls()
		cmd.PrintInfoMsgf("Welcome to PactusBlockchain shell\n\n- Home: https//pactus.org\n- " +
			"Docs: https://docs.pactus.org")
		cmd.PrintLine()
		_prefix = fmt.Sprintf("pactus@%s > ", serverAddr)
	}

	sh.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		if username != "" && password != "" {
			auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
			md := metadata.Pairs("authorization", "Basic "+auth)
			ctx := metadata.NewOutgoingContext(cmd.Context(), md)
			cmd.SetContext(ctx)
		}

		return nil
	}

	changeDefaultParameters := func(c *cobra.Command) *cobra.Command {
		_ = c.PersistentFlags().Lookup("server-addr").Value.Set(defaultServerAddr)
		c.PersistentFlags().Lookup("server-addr").DefValue = defaultServerAddr

		_ = c.PersistentFlags().Lookup("response-format").Value.Set(defaultResponseFormat)
		c.PersistentFlags().Lookup("response-format").DefValue = defaultResponseFormat

		return c
	}

	rootCmd.AddCommand(changeDefaultParameters(pb.BlockchainClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.NetworkClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.TransactionClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.WalletClientCommand()))
	rootCmd.AddCommand(clearScreen())

	err := sh.Execute()
	if err != nil {
		cmd.PrintErrorMsgf("%s", err)
	}
}

func livePrefix() (string, bool) {
	return _prefix, true
}

func clearScreen() *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "clear screen",
		Run: func(_ *cobra.Command, _ []string) {
			cls()
		},
	}
}

func cls() {
	screen.MoveTopLeft()
	screen.Clear()
}
