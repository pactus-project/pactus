package cmd

import (
	"context"

	"github.com/NathanBaulch/protoc-gen-cobra/client"
	"github.com/NathanBaulch/protoc-gen-cobra/naming"
	"github.com/pactus-project/pactus/util"
	pb "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"google.golang.org/grpc"
)

const (
	defaultServerAddr     = "localhost:50051"
	defaultResponseFormat = "prettyjson"
)

func init() {
	client.RegisterFlagBinder(func(fs *pflag.FlagSet, namer naming.Namer) {
		fs.StringVar(&Auth.Username, namer("auth-username"), Auth.Username, "username for authentication")
		fs.StringVar(&Auth.Password, namer("auth-password"), Auth.Password, "password for authentication")
	})

	client.RegisterPreDialer(func(_ context.Context, opts *[]grpc.DialOption) error {
		cred := Auth
		basicAuthCreds := basicAuthCredentials{Token: util.BasicAuth(cred.Username, cred.Password)}

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
