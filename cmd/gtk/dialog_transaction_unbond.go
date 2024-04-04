//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_transaction_unbond.ui
var uiTransactionUnbondDialog []byte

func broadcastTransactionUnbond(wlt *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionUnbondDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_unbond")

	validatorCombo := getComboBoxTextObj(builder, "id_combo_validator")
	validatorHint := getLabelObj(builder, "id_hint_validator")
	memoEntry := getEntryObj(builder, "id_entry_memo")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, ai := range wlt.AllValidatorAddresses() {
		validatorCombo.Append(ai.Address, ai.Address)
	}

	onValidatorChanged := func() {
		receiverEntry, _ := validatorCombo.GetEntry()
		receiverStr, _ := receiverEntry.GetText()
		updateValidatorHint(validatorHint, receiverStr, wlt)
	}

	onSend := func() {
		validatorEntry, _ := validatorCombo.GetEntry()
		validator, _ := validatorEntry.GetText()
		memo, _ := memoEntry.GetText()

		opts := []wallet.TxOption{
			wallet.OptionMemo(memo),
		}

		trx, err := wlt.MakeUnbondTx(validator, opts...)
		if err != nil {
			errorCheck(err)

			return
		}
		msg := fmt.Sprintf(`
You are going to sign and broadcast this transaction:
<tt>
Validator: %s
Memo:      %s
</tt>
<b>THIS ACTION IS NOT REVERSIBLE. Do you want to continue?</b>`, validator, trx.Memo())

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
