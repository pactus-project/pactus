//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/crypto"
)

//go:embed assets/ui/dialog_wallet_create_address.ui
var uiWalletCreateAddressDialog []byte

func createAddress(wdgWallet *widgetWallet) {
	builder, err := gtk.BuilderNewFromString(string(uiWalletCreateAddressDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_wallet_create_address")
	addressLabel := getEntryObj(builder, "id_entry_account_label")

	addressTypeCombo := getComboBoxTextObj(builder, "id_combo_address_type")
	addressTypeCombo.Append(crypto.AddressTypeEd25519Account.String(), "ED25519 Account")
	addressTypeCombo.Append(crypto.AddressTypeBLSAccount.String(), "BLS Account")
	addressTypeCombo.Append(crypto.AddressTypeValidator.String(), "Validator")

	addressTypeCombo.SetActive(0)

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	onOk := func() {
		walletAddressLabel, err := addressLabel.GetText()
		fatalErrorCheck(err)

		walletAddressType := addressTypeCombo.GetActiveID()
		fatalErrorCheck(err)

		switch walletAddressType {
		case crypto.AddressTypeEd25519Account.String():
			password, ok := getWalletPassword(wdgWallet.model.wallet)
			if !ok {
				return
			}

			_, err = wdgWallet.model.wallet.NewEd25519AccountAddress(walletAddressLabel, password)
		case crypto.AddressTypeBLSAccount.String():
			_, err = wdgWallet.model.wallet.NewBLSAccountAddress(walletAddressLabel)
		case crypto.AddressTypeValidator.String():
			_, err = wdgWallet.model.wallet.NewValidatorAddress(walletAddressLabel)
		}

		if err != nil {
			showError(err)

			return
		}

		err = wdgWallet.model.wallet.Save()
		fatalErrorCheck(err)

		wdgWallet.model.rebuildModel()

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
