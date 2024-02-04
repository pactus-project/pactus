package cmd

import (
	pb "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/spf13/cobra"
)

const (
	defaultServerAddr     = "localhost:50051"
	defaultResponseFormat = "prettyjson"
)

func init() {
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
