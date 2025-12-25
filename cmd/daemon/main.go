package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/pactus-project/pactus/version"
	"github.com/spf13/cobra"
)

func init() {
	version.NodeAgent.AppType = "daemon"
}

func main() {
	rootCmd := &cobra.Command{
		Use:               "pactus-daemon",
		Short:             "Pactus daemon",
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	// Hide the "help" sub-command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	buildVersionCmd(rootCmd)
	buildInitCmd(rootCmd)
	buildStartCmd(rootCmd)
	buildPruneCmd(rootCmd)
	buildImportCmd(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		terminal.PrintErrorMsgf(err.Error())
	}
}

func addWorkingDirOption(c *cobra.Command) *string {
	return c.Flags().StringP("working-dir", "w", cmd.PactusDefaultHomeDir(),
		"the path to the working directory that keeps the wallets and node files")
}
