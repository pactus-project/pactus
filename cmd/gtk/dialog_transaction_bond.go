//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_transaction_bond.ui
var uiTransactionBondDialog []byte

func broadcastTransactionBond(wlt *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionBondDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_bond")

	senderEntry := getComboBoxTextObj(builder, "id_combo_sender")
	senderHint := getLabelObj(builder, "id_hint_sender")
	receiverCombo := getComboBoxTextObj(builder, "id_combo_receiver")
	receiverHint := getLabelObj(builder, "id_hint_receiver")
	publicKeyEntry := getEntryObj(builder, "id_entry_public_key")
	amountEntry := getEntryObj(builder, "id_entry_amount")
	amountHint := getLabelObj(builder, "id_hint_amount")
	memoEntry := getEntryObj(builder, "id_entry_memo")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, ai := range wlt.AllAccountAddresses() {
		senderEntry.Append(ai.Address, ai.Address)
	}

	for _, ai := range wlt.AllValidatorAddresses() {
		receiverCombo.Append(ai.Address, ai.Address)
	}

	senderEntry.SetActive(0)

	onSenderChanged := func() {
		senderStr := senderEntry.GetActiveID()
		updateAccountHint(senderHint, senderStr, wlt)
	}

	onReceiverChanged := func() {
		receiverEntry, _ := receiverCombo.GetEntry()
		receiverStr, _ := receiverEntry.GetText()
		updateValidatorHint(receiverHint, receiverStr, wlt)
	}

	onAmountChanged := func() {
		amtStr, _ := amountEntry.GetText()
		updateFeeHint(amountHint, amtStr, wlt, payload.TypeBond)
	}

	onSend := func() {
		sender := senderEntry.GetActiveID()
		receiverEntry, _ := receiverCombo.GetEntry()
		receiver, _ := receiverEntry.GetText()
		publicKey, _ := publicKeyEntry.GetText()
		amountStr, _ := amountEntry.GetText()
		memo, _ := memoEntry.GetText()

		amt, err := amount.FromString(amountStr)
		if err != nil {
			errorCheck(err)

			return
		}

		opts := []wallet.TxOption{
			wallet.OptionMemo(memo),
		}

		trx, err := wlt.MakeBondTx(sender, receiver, publicKey, amt, opts...)
		if err != nil {
			errorCheck(err)

			return
		}
		msg := fmt.Sprintf(`
You are going to sign and broadcast this transaction:
<tt>
From:   %s
To:     %s
Amount: %s
Fee:    %s
Memo:   %s
</tt>
<b>THIS ACTION IS NOT REVERSIBLE. Do you want to continue?</b>`,
			sender, receiver, amt, trx.Fee(), trx.Memo())

		signAndBroadcastTransaction(dlg, msg, wlt, trx)

		dlg.Close()
	}

	onClose := func() {
		dlg.Close()
	}

	signals := map[string]interface{}{
		"on_sender_changed":   onSenderChanged,
		"on_receiver_changed": onReceiverChanged,
		"on_amount_changed":   onAmountChanged,
		"on_send":             onSend,
		"on_cancel":           onClose,
	}
	builder.ConnectSignals(signals)

	onSenderChanged()

	dlg.Run()
}
