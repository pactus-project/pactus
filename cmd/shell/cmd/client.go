package cmd

import (
	"context"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	"github.com/NathanBaulch/protoc-gen-cobra/naming"
	"github.com/pactus-project/pactus/www/grpc/basicauth"
	pb "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

var Auth = &basicauth.BasicAuth{}

const (
	defaultServerAddr     = "localhost:50051"
	defaultResponseFormat = "prettyjson"
)

func init() {
	client.RegisterFlagBinder(func(fs *pflag.FlagSet, namer naming.Namer) {
		fs.StringVar(&Auth.Username, namer("username"), Auth.Username, "username for gRPC basic authentication")
		fs.StringVar(&Auth.Password, namer("password"), Auth.Password, "password for gRPC basic authentication")
	})

	client.RegisterPreDialer(func(_ context.Context, opts *[]grpc.DialOption) error {
		cred := Auth
		basicAuthCreds := basicauth.Credentials{Token: basicauth.MakeCredentials(cred.Username, cred.Password)}

		*opts = append(*opts, grpc.WithPerRPCCredentials(basicAuthCreds))

		return nil
	})

	rootCmd.AddCommand(changeDefaultParameters(pb.BlockchainClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.NetworkClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.TransactionClientCommand()))
	rootCmd.AddCommand(changeDefaultParameters(pb.WalletClientCommand()))
}

func changeDefaultParameters(cmd *cobra.Command) *cobra.Command {
	_ = cmd.PersistentFlags().Lookup("server-addr").Value.Set(defaultServerAddr)
	cmd.PersistentFlags().Lookup("server-addr").DefValue = defaultServerAddr

	_ = cmd.PersistentFlags().Lookup("response-format").Value.Set(defaultResponseFormat)
	cmd.PersistentFlags().Lookup("response-format").DefValue = defaultResponseFormat

	return cmd
}
