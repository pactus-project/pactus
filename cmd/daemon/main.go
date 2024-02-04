package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/version"
	"github.com/spf13/cobra"
)

func init() {
	version.AppType = "daemon"
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

	err := rootCmd.Execute()
	if err != nil {
		cmd.PrintErrorMsgf("%s", err)
	}
}
