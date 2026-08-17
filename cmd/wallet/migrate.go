package main

import (
	"context"

	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildMigrateCmd builds the migrate command to convert a legacy JSON wallet
// to the SQLite format in place.
func buildMigrateCmd(parentCmd *cobra.Command) {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrate the legacy JSON wallet to the SQLite format",
	}
	parentCmd.AddCommand(migrateCmd)

	migrateCmd.Run = func(_ *cobra.Command, _ []string) {
		srcPath := *pathOpt

		err := wallet.Migrate(context.Background(), srcPath)
		terminal.FatalErrorCheck(err)

		terminal.PrintSuccessMsgf("wallet migrated to SQLite format at %s", srcPath)
	}
}
