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

func broadcastTransactionUnbond(model *walletModel) {
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
		receiverStr := getEntryText(receiverEntry)
		updateValidatorHint(validatorHint, receiverStr, wlt)
	}

	onSend := func() {
		validatorEntry, _ := validatorCombo.GetEntry()
		validator := getEntryText(validatorEntry)
		memo := getEntryText(memoEntry)

		opts := []wallet.TxOption{
			wallet.OptionMemo(memo),
		}

		trx, err := wlt.MakeUnbondTx(validator, opts...)
		if err != nil {
			showError(err)

			return
		}
		msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:     Unbond
Validator: %s
Fee:       %s
Memo:      %s
</tt>

You are going to sign and broadcast this transaction.
<b>⚠️ This action cannot be undone.</b>
Do you want to continue with this transaction?`, validator, trx.Fee(), trx.Memo())

		signAndBroadcastTransaction(dlg, msg, model, trx)
	}

	onClose := func() {
		dlg.Close()
	}

	signals := map[string]any{
		"on_validator_changed": onValidatorChanged,
		"on_send":              onSend,
		"on_cancel":            onClose,
	}
	builder.ConnectSignals(signals)

	runDialog(dlg)
}
