package main

import (
	"encoding/base64"
	"fmt"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	"github.com/NathanBaulch/protoc-gen-cobra/naming"
	"github.com/pactus-project/pactus/cmd"
	pb "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc/metadata"
)

const (
	defaultServerAddr     = "localhost:50051"
	defaultResponseFormat = "prettyjson"
)

func main() {
	var (
		username string
		password string
	)

	rootCmd := &cobra.Command{
		Use:   "shell",
		Short: "Pactus Shell",
		Long:  `pactus-shell is a command line tool for interacting with the Pactus blockchain using gRPC`,
	}

	client.RegisterFlagBinder(func(fs *pflag.FlagSet, namer naming.Namer) {
		fs.StringVar(&username, namer("auth-username"), "", "username for gRPC basic authentication")
		fs.StringVar(&password, namer("auth-password"), "", "password for gRPC basic authentication")
	})

	changeDefaultParameters := func(c *cobra.Command) *cobra.Command {
		_ = c.PersistentFlags().Lookup("server-addr").Value.Set(defaultServerAddr)
		c.PersistentFlags().Lookup("server-addr").DefValue = defaultServerAddr

		_ = c.PersistentFlags().Lookup("response-format").Value.Set(defaultResponseFormat)
		c.PersistentFlags().Lookup("response-format").DefValue = defaultResponseFormat

		return c
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		if username != "" && password != "" {
			auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", username, password)))
			md := metadata.Pairs("authorization", "Basic "+auth)
			ctx := metadata.NewOutgoingContext(cmd.Context(), md)
			cmd.SetContext(ctx)
		}

		return nil
	}

	rootCmd.AddCommand(changeDefaultParameters(pb.BlockchainClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.NetworkClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.TransactionClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.WalletClientCommand()))

	err := rootCmd.Execute()
	if err != nil {
		cmd.PrintErrorMsgf("%s", err)
	}
}
