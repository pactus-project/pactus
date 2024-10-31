//go:build gtk

package main

import (
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet"
)

//go:embed assets/ui/dialog_transaction_withdraw.ui
var uiTransactionWithdrawDialog []byte

func broadcastTransactionWithdraw(wlt *wallet.Wallet) {
	builder, err := gtk.BuilderNewFromString(string(uiTransactionWithdrawDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_transaction_withdraw")

	validatorEntry := getComboBoxTextObj(builder, "id_combo_validator")
	validatorHint := getLabelObj(builder, "id_hint_validator")
	receiverCombo := getComboBoxTextObj(builder, "id_combo_receiver")
	receiverHint := getLabelObj(builder, "id_hint_receiver")
	stakeEntry := getEntryObj(builder, "id_entry_stake")
	stakeHint := getLabelObj(builder, "id_hint_stake")
	feeEntry := getEntryObj(builder, "id_entry_fee")
	feeHint := getLabelObj(builder, "id_hint_fee")
	memoEntry := getEntryObj(builder, "id_entry_memo")
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())
	getButtonObj(builder, "id_button_send").SetImage(SendIcon())

	for _, ai := range wlt.AllValidatorAddresses() {
		validatorEntry.Append(ai.Address, ai.Address)
	}
	validatorEntry.SetActive(0)

	for _, ai := range wlt.AllAccountAddresses() {
		receiverCombo.Append(ai.Address, ai.Address)
	}

	onSenderChanged := func() {
		senderStr := validatorEntry.GetActiveID()
		updateValidatorHint(validatorHint, senderStr, wlt)
	}

	onReceiverChanged := func() {
		receiverEntry, _ := receiverCombo.GetEntry()
		receiverStr, _ := receiverEntry.GetText()
		updateAccountHint(receiverHint, receiverStr, wlt)
	}

	onAmountChanged := func() {
		updateAmountHint(stakeHint, validatorEntry.GetActiveID(), wlt)
	}

	onFeeChanged := func() {
		updateFeeHint(feeHint, wlt, payload.TypeTransfer)
	}

	onSend := func() {
		sender := validatorEntry.GetActiveID()
		receiverEntry, _ := receiverCombo.GetEntry()
		receiver, _ := receiverEntry.GetText()
		amountStr, _ := stakeEntry.GetText()
		memo, _ := memoEntry.GetText()

		amt, err := amount.FromString(amountStr)
		if err != nil {
			errorCheck(err)

			return
		}

		opts := []wallet.TxOption{
			wallet.OptionMemo(memo),
		}

		fee, _ := feeEntry.GetText()
		if fee != "" {
			feeAmount, err := amount.FromString(fee)
			cmd.FatalErrorCheck(err)
			opts = append(opts, wallet.OptionFee(feeAmount))
		}

		trx, err := wlt.MakeWithdrawTx(sender, receiver, amt, opts...)
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

	signals := map[string]any{
		"on_sender_changed":   onSenderChanged,
		"on_receiver_changed": onReceiverChanged,
		"on_stake_changed":    onAmountChanged,
		"on_fee_changed":      onFeeChanged,
		"on_send":             onSend,
		"on_cancel":           onClose,
	}
	builder.ConnectSignals(signals)

	onSenderChanged()

	dlg.Run()
}
