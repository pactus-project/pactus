package cmd

import pb "github.com/pactus-project/pactus/www/grpc/gen/cobra"

func init() {
	rootCmd.AddCommand(pb.BlockchainClientCommand())
	rootCmd.AddCommand(pb.NetworkClientCommand())
	rootCmd.AddCommand(pb.TransactionClientCommand())
	rootCmd.AddCommand(pb.WalletClientCommand())
}
