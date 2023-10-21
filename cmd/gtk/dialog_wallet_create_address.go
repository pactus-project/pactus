//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	w "github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_wallet_create_address.ui
var uiWalletCreateAddressDialog []byte

func createAddress(wallet *w.Wallet, addCreatedAddess func(address string, label string)) {
	builder, err := gtk.BuilderNewFromString(string(uiWalletCreateAddressDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_wallet_create_address")
	addressLabel := getEntryObj(builder, "id_entry_account_label")
	oldPasswordLabel := getLabelObj(builder, "id_label_account_label")

	accountTypeCombo := getComboBoxTextObj(builder, "id_combo_account_type")
	accountTypeCombo.Append("bls_account", "Account")
	accountTypeCombo.Append("validator", "Validator")

	accountTypeCombo.SetActive(0)

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	if !wallet.IsEncrypted() {
		addressLabel.SetVisible(false)
		oldPasswordLabel.SetVisible(false)
	}

	onOk := func() {
		walletAddressLabel, err := addressLabel.GetText()
		fatalErrorCheck(err)

		walletAccountType := accountTypeCombo.GetActiveID()
		fatalErrorCheck(err)

		var address string
		if walletAccountType == w.AddressTypeBLSAccount {
			address, err = wallet.NewBLSAccountAddress(walletAddressLabel)
		} else if walletAccountType == w.AddressTypeValidator {
			address, err = wallet.NewValidatorAddress(walletAddressLabel)
		} else {
			errorMsg := fmt.Sprintf("Invalid address type '%s'. Supported address types are '%s' and '%s'", walletAccountType, w.AddressTypeBLSAccount, w.AddressTypeValidator)
			showWarningDialog(dlg, errorMsg)
			return
		}

		errorCheck(err)
		addCreatedAddess(address, walletAddressLabel)

		err = wallet.Save()
		errorCheck(err)

		dlg.Close()
	}

	onCancel := func() {
		dlg.Close()
	}

	// Map the handlers to callback functions, and connect the signals
	// to the Builder.
	signals := map[string]interface{}{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	dlg.Run()
}
