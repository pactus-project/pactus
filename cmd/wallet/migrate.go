package main

import (
	"context"
	"strings"

	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildMigrateCmd builds the migrate command to convert a legacy JSON wallet
// to the SQLite format.
func buildMigrateCmd(parentCmd *cobra.Command) {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "migrate the legacy JSON wallet to the SQLite format",
	}
	parentCmd.AddCommand(migrateCmd)

	toOpt := migrateCmd.Flags().String("to", "",
		"path to save the migrated SQLite wallet (default: the source path without the .json extension)")

	migrateCmd.Run = func(_ *cobra.Command, _ []string) {
		srcPath := *pathOpt
		dstPath := *toOpt
		if dstPath == "" {
			dstPath = strings.TrimSuffix(srcPath, ".json")
			if dstPath == srcPath {
				dstPath = srcPath + ".sqlite"
			}
		}

		wlt, err := wallet.Migrate(context.Background(), srcPath, dstPath, wallet.WithOfflineProvider())
		terminal.FatalErrorCheck(err)
		defer wlt.Close()

		terminal.PrintSuccessMsgf("wallet migrated to %s", dstPath)
		terminal.PrintInfoMsgf("the legacy JSON wallet is kept at %s; remove it manually if it is no longer needed", srcPath)
	}
}
