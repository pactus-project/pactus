package main

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
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
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		for i, info := range wallet.AddressLabels() {
			line := fmt.Sprintf("%v- %s\t", i+1, info.Address)

			if *balanceOpt {
				balance, _ := wallet.Balance(info.Address)
				line += fmt.Sprintf("%v\t", util.ChangeToCoin(balance))
			}

			if *stakeOpt {
				stake, _ := wallet.Stake(info.Address)
				line += fmt.Sprintf("%v\t", util.ChangeToCoin(stake))
			}

			line += info.Label
			if info.Imported {
				line += " (Imported)"
			}

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

	newAddressCmd.Run = func(_ *cobra.Command, _ []string) {
		label := cmd.PromptInput("Label")
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		addr, err := wallet.DeriveNewAddress(label)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
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

		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		balance, err := wallet.Balance(addr)
		cmd.FatalErrorCheck(err)

		stake, err := wallet.Stake(addr)
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

		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		password := getPassword(wallet, *passOpt)
		prv, err := wallet.PrivateKey(password, addr)
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

		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		info := wallet.AddressInfo(addr)
		if info == nil {
			cmd.PrintErrorMsgf("Address not found")
			return
		}

		cmd.PrintLine()
		cmd.PrintInfoMsgf("Public Key: %v", info.Pub.String())
		if !info.Imported {
			cmd.PrintInfoMsgf("Path: %v", info.Path.String())
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

		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		prv, err := bls.PrivateKeyFromString(prvStr)
		cmd.FatalErrorCheck(err)

		password := getPassword(wallet, *passOpt)
		err = wallet.ImportPrivateKey(password, prv)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsgf("Private Key imported. Address: %v",
			prv.PublicKey().Address())
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

		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		oldLabel := wallet.Label(addr)
		newLabel := cmd.PromptInputWithSuggestion("Label", oldLabel)

		err = wallet.SetLabel(addr, newLabel)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsgf("Label set successfully")
	}
}
