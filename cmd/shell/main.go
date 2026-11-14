package main

import (
	"encoding/base64"
	"fmt"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	"github.com/NathanBaulch/protoc-gen-cobra/naming"
	"github.com/c-bata/go-prompt"
	"github.com/inancgumus/screen"
	"github.com/pactus-project/pactus/util/shell"
	"github.com/pactus-project/pactus/util/terminal"
	pb "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/metadata"
)

const (
	defaultServerAddr     = "localhost:50051"
	defaultResponseFormat = "prettyjson"
)

var _prefix string

// createRootCommand creates and configures the root command with all subcommands
// This function contains all the logic from main() but is testable.
func createRootCommand() *cobra.Command {
	var (
		serverAddr string
		username   string
		password   string
	)

	rootCmd := &cobra.Command{
		Use:          "interactive",
		Short:        "Pactus Shell",
		SilenceUsage: true,
		Long:         "pactus-shell is a command line tool for interacting with the Pactus blockchain using gRPC",
	}

	interactive := shell.New(
		rootCmd, nil,
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.Green),
		prompt.OptionDescriptionBGColor(prompt.Black),
		prompt.OptionDescriptionTextColor(prompt.White),
		prompt.OptionLivePrefix(livePrefix),
	)

	client.RegisterFlagBinder(func(fs *pflag.FlagSet, namer naming.Namer) {
		fs.StringVar(&username, namer("auth-username"), "", "username for gRPC basic authentication")
		fs.StringVar(&password, namer("auth-password"), "", "password for gRPC basic authentication")
	})

	interactive.Flags().StringVar(&serverAddr, "server-addr", defaultServerAddr, "gRPC server address")
	interactive.Flags().StringVar(&username, "auth-username", "",
		"username for gRPC basic authentication")

	interactive.Flags().StringVar(&password, "auth-password", "",
		"password for gRPC basic authentication")

	interactive.PreRun = func(_ *cobra.Command, _ []string) {
		// The flags of a command are reset to their defaults before each run.
		// The credentials of the session become these defaults,
		// so they are not lost after the first command.
		setFlagDefault(rootCmd, "auth-username", username)
		setFlagDefault(rootCmd, "auth-password", password)

		cls()
		terminal.PrintInfoMsgf("Welcome to PactusBlockchain interactive mode\n\n- Home: https://pactus.org\n- " +
			"Docs: https://docs.pactus.org")
		terminal.PrintLine()
		_prefix = fmt.Sprintf("pactus@%s > ", serverAddr)
	}

	interactive.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		setAuthContext(cmd, username, password)
	}

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, _ []string) {
		setAuthContext(cmd, username, password)
	}

	changeDefaultParameters := func(cobra *cobra.Command) *cobra.Command {
		_ = cobra.PersistentFlags().Lookup("server-addr").Value.Set(defaultServerAddr)
		cobra.PersistentFlags().Lookup("server-addr").DefValue = defaultServerAddr
		_ = cobra.PersistentFlags().Lookup("response-format").Value.Set(defaultResponseFormat)
		cobra.PersistentFlags().Lookup("response-format").DefValue = defaultResponseFormat

		return cobra
	}

	rootCmd.AddCommand(changeDefaultParameters(pb.BlockchainClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.NetworkClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.TransactionClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.WalletClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.UtilsClientCommand()))
	rootCmd.AddCommand(clearScreen())

	interactive.Use = "interactive"
	interactive.Short = "Start pactus-shell in interactive mode"

	rootCmd.AddCommand(interactive)

	return rootCmd
}

// setFlagDefault sets the default value, and the value, of the named flag
// in the given command and in all its sub-commands.
func setFlagDefault(cmd *cobra.Command, name, value string) {
	for _, flags := range []*pflag.FlagSet{cmd.Flags(), cmd.PersistentFlags()} {
		if flag := flags.Lookup(name); flag != nil {
			flag.DefValue = value
			_ = flag.Value.Set(value)
		}
	}

	for _, sub := range cmd.Commands() {
		setFlagDefault(sub, name, value)
	}
}

func main() {
	rootCmd := createRootCommand()

	err := rootCmd.Execute()
	if err != nil {
		terminal.PrintErrorMsgf(err.Error())
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

func setAuthContext(c *cobra.Command, username, password string) {
	if username != "" && password != "" {
		auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
		md := metadata.Pairs("authorization", "Basic "+auth)
		ctx := metadata.NewOutgoingContext(c.Context(), md)
		c.SetContext(ctx)
	}
}
