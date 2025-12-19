//go:build gtk

package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/types/amount"
)

//go:embed assets/ui/dialog_wallet_set_default_fee.ui
var uiWalletSetDefaultFeeDialog []byte

func setDefaultFee(wdgWallet *widgetWallet) {
	builder, err := gtk.BuilderNewFromString(string(uiWalletSetDefaultFeeDialog))
	fatalErrorCheck(err)

	dlg := getDialogObj(builder, "id_dialog_wallet_set_default_fee")
	feeEntry := getEntryObj(builder, "id_entry_default_fee")
	currentFeeLabel := getLabelObj(builder, "id_label_current_fee_value")

	getButtonObj(builder, "id_button_ok").SetImage(OkIcon())
	getButtonObj(builder, "id_button_cancel").SetImage(CancelIcon())

	info, err := wdgWallet.model.WalletInfo()
	if err != nil {
		showErrorDialog(dlg, fmt.Sprintf("Failed to get wallet info: %v", err))

		return
	}

	currentFee := info.DefaultFee
	currentFeeLabel.SetText(currentFee.String())

	// Set initial value in entry without unit
	feeEntry.SetText(strings.ReplaceAll(currentFee.String(), " PAC", ""))

	onOk := func() {
		feeStr := getEntryText(feeEntry)

		feeAmount, err := amount.FromString(feeStr)
		if err != nil {
			showErrorDialog(dlg, fmt.Sprintf("Invalid fee amount: %v", err))

			return
		}

		err = wdgWallet.model.SetDefaultFee(feeAmount)
		if err != nil {
			showErrorDialog(dlg, fmt.Sprintf("Failed to set default fee: %v", err))

			return
		}

		dlg.Close()

		wdgWallet.rebuild()
	}

	onCancel := func() {
		dlg.Close()
	}

	signals := map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	}
	builder.ConnectSignals(signals)

	runDialog(dlg)
}
