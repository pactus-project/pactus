package main

import (
	"github.com/pactus-project/pactus/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:               "pactus-daemon",
		Short:             "Pactus daemon",
		CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	}

	// Hide the sub-command
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	buildVersionCmd(rootCmd)
	buildInitCmd(rootCmd)
	buildStartCmd(rootCmd)

	err := rootCmd.Execute()
	if err != nil {
		cmd.PrintErrorMsgf("%s", err)
	}
}
