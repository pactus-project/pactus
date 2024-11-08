package main

import (
	"fmt"
	"strings"

	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/bls"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pactus-project/pactus/wallet/vault"
	"github.com/spf13/cobra"
)

// buildAllAddrCmd builds all sub-commands related to addresses.
func buildAllAddrCmd(parentCmd *cobra.Command) {
	addrCmd := &cobra.Command{
		Use:   "address",
		Short: "manage the address book",
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
		Short: "displays all stored addresses",
	}
	parentCmd.AddCommand(allAddressCmd)

	balanceOpt := allAddressCmd.Flags().Bool("balance",
		false, "displays the account balance for each address")

	stakeOpt := allAddressCmd.Flags().Bool("stake",
		false, "displays the validator stake for each address")

	allAddressCmd.Run = func(_ *cobra.Command, _ []string) {
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		for i, info := range wlt.AddressInfos() {
			line := fmt.Sprintf("%v- %s\t", i+1, info.Address)

			if *balanceOpt {
				balance, _ := wlt.Balance(info.Address)
				line += fmt.Sprintf("%s\t", balance.String())
			}

			if *stakeOpt {
				stake, _ := wlt.Stake(info.Address)
				line += fmt.Sprintf("%s\t", stake.String())
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
		Short: "creating a new address",
	}
	parentCmd.AddCommand(newAddressCmd)

	addressType := newAddressCmd.Flags().String("type",
		crypto.AddressTypeEd25519Account.String(), "the type of address: ed25519_account, bls_account and validator")

	newAddressCmd.Run = func(_ *cobra.Command, _ []string) {
		var addressInfo *vault.AddressInfo
		var err error

		label := cmd.PromptInput("Label")
		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		if *addressType == crypto.AddressTypeBLSAccount.String() {
			addressInfo, err = wlt.NewBLSAccountAddress(label)
		} else if *addressType == crypto.AddressTypeEd25519Account.String() {
			password := ""
			if wlt.IsEncrypted() {
				password = cmd.PromptPassword("Password", false)
			}
			addressInfo, err = wlt.NewEd25519AccountAddress(label, password)
		} else if *addressType == crypto.AddressTypeValidator.String() {
			addressInfo, err = wlt.NewValidatorAddress(label)
		} else {
			err = fmt.Errorf("invalid address type '%s'", *addressType)
		}
		cmd.FatalErrorCheck(err)

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintInfoMsgf("%s", addressInfo.Address)
	}
}

// buildBalanceCmd builds a command to display the balance of a given address.
func buildBalanceCmd(parentCmd *cobra.Command) {
	balanceCmd := &cobra.Command{
		Use:   "balance [flags] <ADDRESS>",
		Short: "displays the balance of an address",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(balanceCmd)

	balanceCmd.Run = func(_ *cobra.Command, args []string) {
		addr := args[0]

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()

		balance, _ := wlt.Balance(addr)
		stake, _ := wlt.Stake(addr)
		cmd.PrintInfoMsgf("balance: %s\tstake: %s",
			balance.String(), stake.String())
	}
}

// buildPrivateKeyCmd builds a command to show the private key of a given address.
func buildPrivateKeyCmd(parentCmd *cobra.Command) {
	privateKeyCmd := &cobra.Command{
		Use:   "priv [flags] <ADDRESS>",
		Short: "displays the private key for a specified address",
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
		Short: "displays the public key for a specified address",
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
		Short: "imports a private key into wallet",
	}
	parentCmd.AddCommand(importPrivateKeyCmd)

	passOpt := addPasswordOption(importPrivateKeyCmd)

	importPrivateKeyCmd.Run = func(_ *cobra.Command, _ []string) {
		prvStr := cmd.PromptInput("Private Key")

		wlt, err := openWallet()
		cmd.FatalErrorCheck(err)

		password := getPassword(wlt, *passOpt)

		maybeBLSPrivateKey := func(str string) bool {
			// BLS private keys start with "SECRET1P..." or "TSECRET1P...".
			return strings.Contains(strings.ToLower(str), "secret1p")
		}

		maybeEd25519PrivateKey := func(str string) bool {
			// Ed25519 private keys start with "SECRET1R..." or "TSECRET1R...".
			return strings.Contains(strings.ToLower(str), "secret1r")
		}

		switch {
		case maybeBLSPrivateKey(prvStr):
			blsPrv, err := bls.PrivateKeyFromString(prvStr)
			cmd.FatalErrorCheck(err)

			err = wlt.ImportBLSPrivateKey(password, blsPrv)
			cmd.FatalErrorCheck(err)

		case maybeEd25519PrivateKey(prvStr):
			ed25519Prv, err := ed25519.PrivateKeyFromString(prvStr)
			cmd.FatalErrorCheck(err)

			err = wlt.ImportEd25519PrivateKey(password, ed25519Prv)
			cmd.FatalErrorCheck(err)

		default:
			// The private key cannot be decoded as either BLS or Ed25519.
			cmd.PrintErrorMsgf("Invalid private key.")

			return
		}

		err = wlt.Save()
		cmd.FatalErrorCheck(err)

		cmd.PrintLine()
		cmd.PrintSuccessMsgf("Private Key imported successfully.")
	}
}

// buildSetLabelCmd build a command to set or update the label for an address.
func buildSetLabelCmd(parentCmd *cobra.Command) {
	setLabelCmd := &cobra.Command{
		Use:   "label [flags] <ADDRESS>",
		Short: "assigns or update a label for a specific address",
		Args:  cobra.ExactArgs(1),
	}
	parentCmd.AddCommand(setLabelCmd)

	setLabelCmd.Run = func(_ *cobra.Command, args []string) {
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
