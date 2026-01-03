package main

import (
	"context"

	"github.com/pactus-project/pactus/util/prompt"
	"github.com/pactus-project/pactus/util/terminal"
	"github.com/spf13/cobra"
)

// buildPasswordCmd builds a command to update the wallet's password.
func buildPasswordCmd(parentCmd *cobra.Command) {
	passwordCmd := &cobra.Command{
		Use:   "password",
		Short: "changes the wallet's password",
	}
	parentCmd.AddCommand(passwordCmd)
	passOpt := addPasswordOption(passwordCmd)

	passwordCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet(context.Background())
		terminal.FatalErrorCheck(err)

		oldPassword := getPassword(wlt, *passOpt)
		newPassword := prompt.PromptPassword("New Password", true)

		err = wlt.UpdatePassword(oldPassword, newPassword)
		terminal.FatalErrorCheck(err)

		terminal.PrintLine()
		terminal.PrintWarnMsgf("Your wallet password successfully updated.")
	}
}
