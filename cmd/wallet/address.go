package main

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
	"github.com/spf13/cobra"
)

// buildAllAddrCmd builds all sub-commands related to addresses.
func buildAllAddrCmd(parentCmd *cobra.Command) {
	addrCmd := &cobra.Command{
		Use:   "address",
		Short: "Manage the address book",
	}

	parentCmd.AddCommand(addrCmd)
	buildAllAddressesCmd(addrCmd)
	buildNewAddressCmd(addrCmd)
	buildBalanceCmd(addrCmd)
	buildPrivateKeyCmd(addrCmd)
	buildPublicKeyCmd(addrCmd)
	buildImportPrivateKeyCmd(addrCmd)
	buildSetLabelCmd(addrCmd)
}

// buildAllAddressesCmd builds a command to list all addresses from the wallet.
func buildAllAddressesCmd(parentCmd *cobra.Command) {
	allAddressCmd := &cobra.Command{
		Use:   "all",
		Short: "Display all stored addresses",
	}
	parentCmd.AddCommand(allAddressCmd)

	balanceOpt := allAddressCmd.Flags().Bool("balance",
		false, "Display the account balance for each address")

	stakeOpt := allAddressCmd.Flags().Bool("stake",
		false, "Display the validator stake for each address")

	allAddressCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		for i, info := range wlt.AddressInfos() {
			line := fmt.Sprintf("%v- %s\t", i+1, info.Address)

			if *balanceOpt {
				balance, _ := wlt.Balance(info.Address)
				line += fmt.Sprintf("%v\t", util.ChangeToCoin(balance))
			}

			if *stakeOpt {
				stake, _ := wlt.Stake(info.Address)
				line += fmt.Sprintf("%v\t", util.ChangeToCoin(stake))
			}

			line += info.Label
			cmd.PrintInfoMsgf(line)
		}
	}
}

// buildNewAddressCmd builds a command for creating a new wallet address.
func buildNewAddressCmd(parentCmd *cobra.Command) {
	newAddressCmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new address",
	}
	parentCmd.AddCommand(newAddressCmd)

	addressType := newAddressCmd.Flags().String("type",
		wallet.AddressTypeBLSAccount, "the type of address: bls_account or validator")

	newAddressCmd.Run = func(_ *cobra.Command, _ []string) {
		var addr string
		var err error

		label := cmd.PromptInput("Label")
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		if *addressType == wallet.AddressTypeBLSAccount {
			addr, err = wlt.NewBLSAccountAddress(label)
		} else if *addressType == wallet.AddressTypeValidator {
			addr, err = wlt.NewValidatorAddress(label)
		} else {
			err = fmt.Errorf("invalid address type '%s'", *addressType)
		}
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("%s", addr)
	}
}

// buildBalanceCmd builds a command to display the balance of a given address.
func buildBalanceCmd(parentCmd *cobra.Command) {
	balanceCmd := &cobra.Command{
		Use:   "balance [flags] <ADDRESS>",
		Short: "Display the balance of an address",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(balanceCmd)

	balanceCmd.Run = func(_ *cobra.Command, args []string) {
		addr := args[0]

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		balance, err := wlt.Balance(addr)
		cmd.FatalErrorCheck(err)

		stake, err := wlt.Stake(addr)
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("balance: %v\tstake: %v",
			util.ChangeToCoin(balance), util.ChangeToCoin(stake))
	}
}

// buildPrivateKeyCmd builds a command to show the private key of a given address.
func buildPrivateKeyCmd(parentCmd *cobra.Command) {
	privateKeyCmd := &cobra.Command{
		Use:   "priv [flags] <ADDRESS>",
		Short: "Display the private key for a specified address",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(privateKeyCmd)

	passOpt := addPasswordOption(privateKeyCmd)

	privateKeyCmd.Run = func(_ *cobra.Command, args []string) {
		addr := args[0]

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		password := getPassword(wlt, *passOpt)
		prv, err := wlt.PrivateKey(password, addr)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintWarnMsgf("Private Key: %v", prv)
	}
}

// buildPublicKeyCmd builds a command to show the public key of a given address.
func buildPublicKeyCmd(parentCmd *cobra.Command) {
	publicKeyCmd := &cobra.Command{
		Use:   "pub [flags] <ADDRESS>",
		Short: "Display the public key for a specified address",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(publicKeyCmd)

	publicKeyCmd.Run = func(_ *cobra.Command, args []string) {
		addr := args[0]

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		info := wlt.AddressInfo(addr)
		if info == nil {
			cmd.PrintErrorMsgf("Address not found")

			return
		}

		cmd.PrintLine()
		cmd.PrintInfoMsgf("Public Key: %v", info.PublicKey)
		if info.Path != "" {
			cmd.PrintInfoMsgf("Path: %v", info.Path)
		}
	}
}

// buildImportPrivateKeyCmd build a command to import a private key into the wallet.
func buildImportPrivateKeyCmd(parentCmd *cobra.Command) {
	importPrivateKeyCmd := &cobra.Command{
		Use:   "import",
		Short: "Import a private key into wallet",
	}
	parentCmd.AddCommand(importPrivateKeyCmd)

	passOpt := addPasswordOption(importPrivateKeyCmd)

	importPrivateKeyCmd.Run = func(_ *cobra.Command, _ []string) {
		prvStr := cmd.PromptInput("Private Key")

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		prv, err := bls.PrivateKeyFromString(prvStr)
		cmd.FatalErrorCheck(err)

		password := getPassword(wlt, *passOpt)
		err = wlt.ImportPrivateKey(password, prv)
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsgf("Private Key imported.") // TODO: display imported addresses
	}
}

// buildSetLabelCmd build a command to set or update the label for an address.
func buildSetLabelCmd(parentCmd *cobra.Command) {
	setLabelCmd := &cobra.Command{
		Use:   "label [flags] <ADDRESS>",
		Short: "Assign or update a label for a specific address",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(setLabelCmd)

	setLabelCmd.Run = func(c *cobra.Command, args []string) {
		addr := args[0]

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		oldLabel := wlt.Label(addr)
		newLabel := cmd.PromptInputWithSuggestion("Label", oldLabel)

		err = wlt.SetLabel(addr, newLabel)
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsgf("Label set successfully")
	}
}
