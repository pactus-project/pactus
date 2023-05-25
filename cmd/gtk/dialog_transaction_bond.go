//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_transaction_bond.ui
var uiTransactionBondDialog []byte

func broadcastTransactionBond(wallet *wallet.Wallet, valAddrs []crypto.Address) {
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
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, i := range wallet.AddressLabels() {
		senderEntry.Append(i.Address, i.Address)
	}

	for _, addr := range valAddrs {
		receiverCombo.Append(addr.String(), addr.String())
	}

	senderEntry.SetActive(0)

	onSenderChanged := func() {
		senderStr := senderEntry.GetActiveID()
		updateAccountHint(senderHint, senderStr, wallet)
	}

	onReceiverChanged := func() {
		receiverEntry, _ := receiverCombo.GetEntry()
		receiverStr, _ := receiverEntry.GetText()
		updateValidatorHint(receiverHint, receiverStr, wallet)
	}

	onAmountChanged := func() {
		amtStr, _ := amountEntry.GetText()
		updateFeeHint(amountHint, amtStr, wallet)
	}

	onSend := func() {
		sender := senderEntry.GetActiveID()
		receiverEntry, _ := receiverCombo.GetEntry()
		receiver, _ := receiverEntry.GetText()
		publicKey, _ := publicKeyEntry.GetText()
		amountStr, _ := amountEntry.GetText()

		amount, err := util.StringToChange(amountStr)
		if err != nil {
			errorCheck(err)
			return
		}

		trx, err := wallet.MakeBondTx(sender, receiver, publicKey, amount)
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

	THIS ACTION IS NOT REVERSIBLE Do you want to continue?`, sender, receiver,
			util.ChangeToString(amount), util.ChangeToString(trx.Fee()))

		if showQuestionDialog(dlg, msg) {
			password, ok := getWalletPassword(wallet)
			if !ok {
				return
			}
			err := wallet.SignTransaction(password, trx)
			if err != nil {
				errorCheck(err)
				return
			}
			_, err = wallet.BroadcastTransaction(trx)
			if err != nil {
				errorCheck(err)
				return
			}
		}

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
