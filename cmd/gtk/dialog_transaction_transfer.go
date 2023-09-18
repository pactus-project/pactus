//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_transaction_transfer.ui
var uiTransactionTransferDialog []byte

func broadcastTransactionSend(wallet *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionTransferDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_transfer")

	senderEntry := getComboBoxTextObj(builder, "id_combo_sender")
	senderHint := getLabelObj(builder, "id_hint_sender")
	receiverEntry := getEntryObj(builder, "id_entry_receiver")
	receiverHint := getLabelObj(builder, "id_hint_receiver")
	amountEntry := getEntryObj(builder, "id_entry_amount")
	amountHint := getLabelObj(builder, "id_hint_amount")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, i := range wallet.AddressLabels() {
		senderEntry.Append(i.Address, i.Address)
	}
	senderEntry.SetActive(0)

	onSenderChanged := func() {
		senderStr := senderEntry.GetActiveID()
		updateAccountHint(senderHint, senderStr, wallet)
	}

	onReceiverChanged := func() {
		receiverStr, _ := receiverEntry.GetText()
		updateAccountHint(receiverHint, receiverStr, wallet)
	}

	onAmountChanged := func() {
		amtStr, _ := amountEntry.GetText()
		updateFeeHint(amountHint, amtStr, wallet, payload.TypeTransfer)
	}

	onSend := func() {
		sender := senderEntry.GetActiveID()
		receiver, _ := receiverEntry.GetText()
		amountStr, _ := amountEntry.GetText()

		amount, err := util.StringToChange(amountStr)
		if err != nil {
			errorCheck(err)
			return
		}

		trx, err := wallet.MakeTransferTx(sender, receiver, amount)
		if err != nil {
			errorCheck(err)
			return
		}
		msg := fmt.Sprintf(`
You are going to sign and broadcast this transaction:

From: <b>%v</b>
To: <b>%v</b>
Amount: <b>%v</b>
Fee: <b>%v</b>

THIS ACTION IS NOT REVERSIBLE. Do you want to continue?`, sender, receiver,
			util.ChangeToString(amount), util.ChangeToString(trx.Fee()))

		signAndBroadcastTransaction(dlg, msg, wallet, trx)

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
