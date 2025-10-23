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

//go:embed assets/ui/dialog_transaction_transfer.ui
var uiTransactionTransferDialog []byte

func broadcastTransactionTransfer(wlt *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionTransferDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_transfer")

	senderEntry := getComboBoxTextObj(builder, "id_combo_sender")
	senderHint := getLabelObj(builder, "id_hint_sender")
	receiverEntry := getEntryObj(builder, "id_entry_receiver")
	receiverHint := getLabelObj(builder, "id_hint_receiver")
	amountEntry := getEntryObj(builder, "id_entry_amount")
	amountHint := getLabelObj(builder, "id_hint_amount")
	feeEntry := getEntryObj(builder, "id_entry_fee")
	feeHint := getLabelObj(builder, "id_hint_fee")
	memoEntry := getEntryObj(builder, "id_entry_memo")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	feeEntry.SetText(fmt.Sprintf("%g", wlt.Info().DefaultFee.ToPAC()))

	for _, i := range wlt.AllAccountAddresses() {
		senderEntry.Append(i.Address, i.Address)
	}
	senderEntry.SetActive(0)

	onSenderChanged := func() {
		senderStr := senderEntry.GetActiveID()
		updateAccountHint(senderHint, senderStr, wlt)
		updateBalanceHint(amountHint, senderEntry.GetActiveID(), wlt)
	}

	onReceiverChanged := func() {
		receiverStr, _ := receiverEntry.GetText()
		updateAccountHint(receiverHint, receiverStr, wlt)
	}
	onFeeChanged := func() {
		updateFeeHint(feeHint, wlt, payload.TypeTransfer)
	}

	onSend := func() {
		sender := senderEntry.GetActiveID()
		receiver, _ := receiverEntry.GetText()
		amountStr, _ := amountEntry.GetText()
		memo, _ := memoEntry.GetText()

		amt, err := amount.FromString(amountStr)
		if err != nil {
			showError(err)

			return
		}

		feeStr, _ := feeEntry.GetText()
		opts := []wallet.TxOption{
			wallet.OptionMemo(memo),
			wallet.OptionFee(feeStr),
		}

		trx, err := wlt.MakeTransferTx(sender, receiver, amt, opts...)
		if err != nil {
			showError(err)

			return
		}
		msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:   Transfer
From:   %s
To:     %s
Amount: %s
Fee:    %s
Memo:   %s
</tt>

You are going to sign and broadcast this transaction.
<b>⚠️ This action cannot be undone.</b>
Do you want to continue with this transaction?`,
			sender, receiver, amt, trx.Fee(), trx.Memo())

		signAndBroadcastTransaction(dlg, msg, wlt, trx)
	}

	onClose := func() {
		dlg.Close()
	}

	signals := map[string]any{
		"on_sender_changed":   onSenderChanged,
		"on_receiver_changed": onReceiverChanged,
		"on_fee_changed":      onFeeChanged,
		"on_send":             onSend,
		"on_cancel":           onClose,
	}
	builder.ConnectSignals(signals)

	onSenderChanged()

	runDialog(dlg)
}
