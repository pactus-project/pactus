//go:build gtk

package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/wallet"
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

	okBtn := getButtonObj(builder, "id_button_ok")
	okBtn.SetImage(OkIcon())

	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	dlg.Connect("key-press-event", func(_ *gtk.Dialog, event *gdk.Event) {
		keyEvent := &gdk.EventKey{Event: event}
		// Check if the Enter key was pressed
		if keyEvent.KeyVal() == gdk.KEY_Return || keyEvent.KeyVal() == gdk.KEY_KP_Enter {
			// Prevent Enter from triggering default behavior twice
			_, _ = okBtn.Emit("clicked")
		}
	})

	onOk := func() {
		okBtn.SetSensitive(false)

		walletAddressLabel := getEntryText(addressLabel)

		walletAddressType := addressTypeCombo.GetActiveID()
		fatalErrorCheck(err)

		switch walletAddressType {
		case crypto.AddressTypeEd25519Account.String():
			password, ok := getWalletPassword(wdgWallet.model)
			if !ok {
				return
			}

			_, err = wdgWallet.model.NewAddress(crypto.AddressTypeEd25519Account, walletAddressLabel,
				wallet.WithPassword(password))
		case crypto.AddressTypeBLSAccount.String():
			_, err = wdgWallet.model.NewAddress(crypto.AddressTypeBLSAccount, walletAddressLabel)
		case crypto.AddressTypeValidator.String():
			_, err = wdgWallet.model.NewAddress(crypto.AddressTypeValidator, walletAddressLabel)
		}

		if err != nil {
			okBtn.SetSensitive(true)
			showError(err)

			return
		}

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

	runDialog(dlg)
}
