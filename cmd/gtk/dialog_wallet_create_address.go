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
	addressTypeCombo.Append(wallet.AddressTypeEd25519Account, "ED25519 Account")
	addressTypeCombo.Append(wallet.AddressTypeBLSAccount, "BLS Account")
	addressTypeCombo.Append(wallet.AddressTypeValidator, "Validator")

	addressTypeCombo.SetActive(0)

	getLabelObj(builder, "id_label_account_password")
	passwordInput := getEntryObj(builder, "id_entry_account_password")

	addressTypeCombo.Connect("changed", func() {
		activeID := addressTypeCombo.GetActiveID()
		if activeID == wallet.AddressTypeEd25519Account {
			passwordInput.SetSensitive(true)
		} else {
			passwordInput.SetSensitive(false)
		}
	})

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	onOk := func() {
		walletAddressLabel, err := addressLabel.GetText()
		fatalErrorCheck(err)

		walletAddressType := addressTypeCombo.GetActiveID()
		fatalErrorCheck(err)

		password, err := passwordInput.GetText()
		fatalErrorCheck(err)

		if walletAddressType == wallet.AddressTypeEd25519Account && password == "" {
			passwordInput.SetName("entry_error")

			dialog := gtk.MessageDialogNew(dlg, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK,
				"Password is required for ED25519 Account.")
			dialog.Run()
			dialog.Destroy()
			passwordInput.SetSensitive(true)

			return
		}

		if walletAddressType == wallet.AddressTypeEd25519Account {
			_, err = ww.model.wallet.NewEd25519AccountAddress(walletAddressLabel, password)
		} else if walletAddressType == wallet.AddressTypeBLSAccount {
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
	signals := map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	dlg.SetModal(true)

	dlg.Run()
}
