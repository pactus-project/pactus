package main

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/util"
	"github.com/spf13/cobra"
)

func buildAllAddrCmd(parentCmd *cobra.Command) {
	addrCmd := &cobra.Command{
		Use:   "address",
		Short: "Manage address book",
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

// AllAddresses lists all the wallet addresses.
func buildAllAddressesCmd(parentCmd *cobra.Command) {
	allAddressCmd := &cobra.Command{
		Use:   "all",
		Short: "Show all addresses",
	}
	parentCmd.AddCommand(allAddressCmd)

	balanceOpt := allAddressCmd.Flags().Bool("balance",
		false, "show account balance")

	stakeOpt := allAddressCmd.Flags().Bool("stake",
		false, "show validator stake")

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

// NewAddress creates a new address.
func buildNewAddressCmd(parentCmd *cobra.Command) {
	newAddressCmd := &cobra.Command{
		Use:   "new",
		Short: "Creating a new address",
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

// Balance shows the balance of an address.
func buildBalanceCmd(parentCmd *cobra.Command) {
	balanceCmd := &cobra.Command{
		Use:   "balance",
		Short: "Show the balance of an address",
	}
	parentCmd.AddCommand(balanceCmd)

	addrArg := addAddressArg(parentCmd)

	balanceCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		balance, err := wallet.Balance(*addrArg)
		cmd.FatalErrorCheck(err)

		stake, err := wallet.Stake(*addrArg)
		cmd.FatalErrorCheck(err)

		cmd.PrintInfoMsgf("balance: %v\tstake: %v",
			util.ChangeToCoin(balance), util.ChangeToCoin(stake))
	}
}

// PrivateKey returns the private key of an address.
func buildPrivateKeyCmd(parentCmd *cobra.Command) {
	privateKeyCmd := &cobra.Command{
		Use:   "priv",
		Short: "Show the private key of an address",
	}
	parentCmd.AddCommand(privateKeyCmd)

	addrArg := addAddressArg(parentCmd)
	passOpt := addPasswordOption(parentCmd)

	privateKeyCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		password := getPassword(wallet, *passOpt)
		prv, err := wallet.PrivateKey(password, *addrArg)
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintWarnMsgf("Private Key: %v", prv)
	}
}

// PublicKey returns the public key of an address.
func buildPublicKeyCmd(parentCmd *cobra.Command) {
	publicKeyCmd := &cobra.Command{
		Use:   "pub",
		Short: "Show the public key of an address",
	}
	parentCmd.AddCommand(publicKeyCmd)

	addrArg := addAddressArg(parentCmd)

	publicKeyCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		info := wallet.AddressInfo(*addrArg)
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

// ImportPrivateKey imports a private key into the wallet.
func buildImportPrivateKeyCmd(parentCmd *cobra.Command) {
	importPrivateKeyCmd := &cobra.Command{
		Use:   "import",
		Short: "Import a private key into wallet",
	}
	parentCmd.AddCommand(importPrivateKeyCmd)

	passOpt := addPasswordOption(parentCmd)

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

// SetLabel set label for the address.
func buildSetLabelCmd(parentCmd *cobra.Command) {
	setLabelCmd := &cobra.Command{
		Use:   "label",
		Short: "Set label for the an address",
	}
	parentCmd.AddCommand(setLabelCmd)

	addrArg := addAddressArg(parentCmd)

	setLabelCmd.Run = func(_ *cobra.Command, _ []string) {
		wallet, err := openWallet()
		cmd.FatalErrorCheck(err)

		oldLabel := wallet.Label(*addrArg)
		newLabel := cmd.PromptInputWithSuggestion("Label", oldLabel)

		err = wallet.SetLabel(*addrArg, newLabel)
		cmd.FatalErrorCheck(err)

		err = wallet.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsgf("Label set successfully")
	}
}

func addAddressArg(c *cobra.Command) *string {
	return c.Flags().String("ADDRESS",
		"", "address string")
}
