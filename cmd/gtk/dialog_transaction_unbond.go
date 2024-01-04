//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_transaction_unbond.ui
var uiTransactionUnBondDialog []byte

func broadcastTransactionUnBond(wlt *wallet.Wallet, valAddrs []crypto.Address) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionUnBondDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_unbond")

	validatorCombo := getComboBoxTextObj(builder, "id_combo_validator")
	validatorHint := getLabelObj(builder, "id_hint_validator")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, addr := range valAddrs {
		validatorCombo.Append(addr.String(), addr.String())
	}

	onValidatorChanged := func() {
		receiverEntry, _ := validatorCombo.GetEntry()
		receiverStr, _ := receiverEntry.GetText()
		updateValidatorHint(validatorHint, receiverStr, wlt)
	}

	onSend := func() {
		validatorEntry, _ := validatorCombo.GetEntry()
		validator, _ := validatorEntry.GetText()

		trx, err := wlt.MakeUnbondTx(validator)
		if err != nil {
			errorCheck(err)
			return
		}
		msg := fmt.Sprintf(`
You are going to sign and broadcast this transaction:

Validator: %v

THIS ACTION IS NOT REVERSIBLE. Do you want to continue?`, validator)

		signAndBroadcastTransaction(dlg, msg, wlt, trx)

		dlg.Close()
	}

	onClose := func() {
		dlg.Close()
	}

	signals := map[string]interface{}{
		"on_validator_changed": onValidatorChanged,
		"on_send":              onSend,
		"on_cancel":            onClose,
	}
	builder.ConnectSignals(signals)

	dlg.Run()
}
