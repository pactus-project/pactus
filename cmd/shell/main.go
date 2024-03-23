package main

import (
	"context"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	"github.com/NathanBaulch/protoc-gen-cobra/naming"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/www/grpc/basicauth"
	pb "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

const (
	defaultServerAddr     = "localhost:50051"
	defaultResponseFormat = "prettyjson"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "shell",
		Short: "Pactus Shell",
		Long:  `pactus-shell is a command line tool for interacting with the Pactus blockchain using gRPC`,
	}

	auth := &basicauth.BasicAuth{}

	client.RegisterFlagBinder(func(fs *pflag.FlagSet, namer naming.Namer) {
		fs.StringVar(&auth.Username, namer("auth-username"), "", "username for gRPC basic authentication")
		fs.StringVar(&auth.Password, namer("auth-password"), "", "password for gRPC basic authentication")
	})

	client.RegisterPreDialer(func(_ context.Context, opts *[]grpc.DialOption) error {
		*opts = append(*opts, grpc.WithPerRPCCredentials(auth))

		return nil
	})

	changeDefaultParameters := func(cmd *cobra.Command) *cobra.Command {
		_ = cmd.PersistentFlags().Lookup("server-addr").Value.Set(defaultServerAddr)
		cmd.PersistentFlags().Lookup("server-addr").DefValue = defaultServerAddr

		_ = cmd.PersistentFlags().Lookup("response-format").Value.Set(defaultResponseFormat)
		cmd.PersistentFlags().Lookup("response-format").DefValue = defaultResponseFormat

		return cmd
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
