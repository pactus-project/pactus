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

func broadcastTransactionWithdraw(wlt *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionTransferDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_withdraw")

	validatorEntry := getComboBoxTextObj(builder, "id_combo_validator")
	validatorHint := getLabelObj(builder, "id_hint_validator")
	receiverEntry := getEntryObj(builder, "id_entry_receiver")
	receiverHint := getLabelObj(builder, "id_hint_receiver")
	stakeEntry := getEntryObj(builder, "id_entry_stake")
	stakeHint := getLabelObj(builder, "id_hint_stake")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, i := range wlt.AddressInfos() {
		validatorEntry.Append(i.Address, i.Address)
	}
	validatorEntry.SetActive(0)

	onSenderChanged := func() {
		senderStr := validatorEntry.GetActiveID()
		updateAccountHint(validatorHint, senderStr, wlt)
	}

	onReceiverChanged := func() {
		receiverStr, _ := receiverEntry.GetText()
		updateAccountHint(receiverHint, receiverStr, wlt)
	}

	onAmountChanged := func() {
		amtStr, _ := stakeEntry.GetText()
		updateFeeHint(stakeHint, amtStr, wlt, payload.TypeTransfer)
	}

	onSend := func() {
		sender := validatorEntry.GetActiveID()
		receiver, _ := receiverEntry.GetText()
		amountStr, _ := stakeEntry.GetText()

		amount, err := util.StringToChange(amountStr)
		if err != nil {
			errorCheck(err)
			return
		}

		trx, err := wlt.MakeWithdrawTx(sender, receiver, amount)
		if err != nil {
			errorCheck(err)
			return
		}
		msg := fmt.Sprintf(`
You are going to sign and broadcast this transaction:

From:   %v
To:     %v
Amount: %v
Fee:    %v

THIS ACTION IS NOT REVERSIBLE. Do you want to continue?`, sender, receiver,
			util.ChangeToString(amount), util.ChangeToString(trx.Fee()))

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
