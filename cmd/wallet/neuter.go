package main

import (
	"fmt"

	"github.com/ipfs/boxo/util"
	"github.com/pactus-project/pactus/cmd"
	"github.com/spf13/cobra"
)

func buildNeuterCmd(parentCmd *cobra.Command) {
	neuterCmd := &cobra.Command{
		Use:   "neuter",
		Short: "convert full wallet to neutered wallet and can only be used to retrieve balances or stakes",
	}
	parentCmd.AddCommand(neuterCmd)

	neuterCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		path := wlt.Path() + ".neutered"

		if util.FileExists(path) {
			cmd.FatalErrorCheck(fmt.Errorf("neutered wallet already exists, at %s", path))
		}

		neuteredWallet := wlt.Neuter(path)

		err = neuteredWallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintSuccessMsgf("neutered wallet created at %s", path)
	}
}
