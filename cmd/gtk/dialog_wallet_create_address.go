//go:build gtk

package main

import (
	_ "embed"

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

	onOk := func() {
		walletAddressLabel, err := addressLabel.GetText()
		fatalErrorCheck(err)

		walletAddressType := addressTypeCombo.GetActiveID()
		fatalErrorCheck(err)

		switch walletAddressType {
		case wallet.AddressTypeEd25519Account:
			passwordFetcher := func(wlt *wallet.Wallet) (string, bool) {
				if *passwordOpt != "" {
					return *passwordOpt, true
				}

				return getWalletPassword(wlt)
			}

			password, ok := passwordFetcher(ww.model.wallet)
			if !ok {
				return
			}

			_, err = ww.model.wallet.NewEd25519AccountAddress(walletAddressLabel, password)
		case wallet.AddressTypeBLSAccount:
			_, err = ww.model.wallet.NewBLSAccountAddress(walletAddressLabel)
		case wallet.AddressTypeValidator:
			_, err = ww.model.wallet.NewValidatorAddress(walletAddressLabel)
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
