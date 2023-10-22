//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"

	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_wallet_create_address.ui
var uiWalletCreateAddressDialog []byte

func createAddress(ww *widgetWallet) {
	builder, err := gtk.BuilderNewFromString(string(uiWalletCreateAddressDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_wallet_create_address")
	addressLabel := getEntryObj(builder, "id_entry_account_label")

	addressTypeCombo := getComboBoxTextObj(builder, "id_combo_address_type")
	addressTypeCombo.Append(wallet.AddressTypeBLSAccount, "Account")
	addressTypeCombo.Append(wallet.AddressTypeValidator, "Validator")

	addressTypeCombo.SetActive(0)

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	onOk := func() {
		walletAddressLabel, err := addressLabel.GetText()
		fatalErrorCheck(err)

		walletAddressType := addressTypeCombo.GetActiveID()
		fatalErrorCheck(err)

		if walletAddressType == wallet.AddressTypeBLSAccount {
			_, err = ww.model.wallet.NewBLSAccountAddress(walletAddressLabel)
		} else if walletAddressType == wallet.AddressTypeValidator {
			_, err = ww.model.wallet.NewValidatorAddress(walletAddressLabel)
		} else {
			err = fmt.Errorf("invalid address type '%s'", walletAddressType)
		}

		errorCheck(err)

		err = ww.model.wallet.Save()
		errorCheck(err)

		ww.model.rebuildModel()

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
