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

func broadcastTransactionBond(model *walletModel) {
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
	feeEntry := getEntryObj(builder, "id_entry_fee")
	feeHint := getLabelObj(builder, "id_hint_fee")
	memoEntry := getEntryObj(builder, "id_entry_memo")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	feeEntry.SetText(fmt.Sprintf("%g", wlt.Info().DefaultFee.ToPAC()))

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
		updateBalanceHint(amountHint, senderEntry.GetActiveID(), wlt)
	}

	onReceiverChanged := func() {
		receiverEntry, _ := receiverCombo.GetEntry()
		receiverStr := getEntryText(receiverEntry)
		updateValidatorHint(receiverHint, receiverStr, wlt)
	}

	onFeeChanged := func() {
		updateFeeHint(feeHint, wlt, payload.TypeTransfer)
	}

	onSend := func() {
		sender := senderEntry.GetActiveID()
		receiverEntry, _ := receiverCombo.GetEntry()
		receiver := getEntryText(receiverEntry)
		publicKey := getEntryText(publicKeyEntry)
		amountStr := getEntryText(amountEntry)
		memo := getEntryText(memoEntry)

		amt, err := amount.FromString(amountStr)
		if err != nil {
			showError(err)

			return
		}

		feeStr := getEntryText(feeEntry)
		opts := []wallet.TxOption{
			wallet.OptionMemo(memo),
			wallet.OptionFee(feeStr),
		}

		trx, err := wlt.MakeBondTx(sender, receiver, publicKey, amt, opts...)
		if err != nil {
			showError(err)

			return
		}
		msg := fmt.Sprintf(`
📝 Transaction Details:
<tt>
Type:   Bond
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
