//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_transaction_bond.ui
var uiTransactionBondDialog []byte

func broadcastTransactionBond(wallet *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionBondDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_bond")

	senderEntry := getComboBoxTextObj(builder, "id_combo_sender")
	balanceLabel := getLabelObj(builder, "id_label_balance")
	receiverEntry := getEntryObj(builder, "id_entry_receiver")
	publicKeyEntry := getEntryObj(builder, "id_entry_public_key")
	amountEntry := getEntryObj(builder, "id_entry_amount")
	payableLabel := getLabelObj(builder, "id_label_payable")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, i := range wallet.AddressLabels() {
		senderEntry.Append(i.Address, i.Address)
	}
	senderEntry.SetActive(0)

	onSenderChanged := func() {
		senderStr := senderEntry.GetActiveID()
		info := wallet.AddressInfo(senderStr)
		balance, _ := wallet.Balance(senderStr)
		balanceLabel.SetMarkup(
			fmt.Sprintf("<span foreground='gray' size='small'>balance: %v, label: %v</span>",
				util.ChangeToString(balance), info.Label))
	}

	onAmountChanged := func() {
		amountStr, _ := amountEntry.GetText()
		amount, _ := util.StringToChange(amountStr)
		fee := wallet.CalculateFee(amount)

		payableLabel.SetMarkup(
			fmt.Sprintf("<span foreground='gray' size='small'>payable: %v, fee: %v</span>",
				util.ChangeToString(fee+amount), util.ChangeToString(fee)))
	}
	onSend := func() {
		sender := senderEntry.GetActiveID()
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
		"on_sender_changed": onSenderChanged,
		"on_amount_changed": onAmountChanged,
		"on_send":           onSend,
		"on_cancel":         onClose,
	}
	builder.ConnectSignals(signals)

	onSenderChanged()

	dlg.Run()
}
