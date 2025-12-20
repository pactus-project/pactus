//go:build gtk

package controller

import (
	"fmt"
	"strings"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/wallet"
)

type WalletDefaultFeeModel interface {
	WalletInfo() (wallet.Info, error)
	SetDefaultFee(fee amount.Amount) error
}

type WalletDefaultFeeDialogController struct {
	view *view.WalletDefaultFeeDialogView
}

func NewWalletDefaultFeeDialogController(view *view.WalletDefaultFeeDialogView) *WalletDefaultFeeDialogController {
	return &WalletDefaultFeeDialogController{view: view}
}

func (c *WalletDefaultFeeDialogController) Run(model WalletDefaultFeeModel, afterSave func()) {
	info, err := model.WalletInfo()
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Dialog, fmt.Sprintf("Failed to get wallet info: %v", err))

		return
	}

	currentFee := info.DefaultFee
	c.view.CurrentFeeLabel.SetText(currentFee.String())
	c.view.FeeEntry.SetText(strings.ReplaceAll(currentFee.String(), " PAC", ""))

	onOk := func() {
		feeStr := gtkutil.GetEntryText(c.view.FeeEntry)
		feeAmount, err := amount.FromString(feeStr)
		if err != nil {
			gtkutil.ShowErrorDialog(c.view.Dialog, fmt.Sprintf("Invalid fee amount: %v", err))

			return
		}
		if err := model.SetDefaultFee(feeAmount); err != nil {
			gtkutil.ShowErrorDialog(c.view.Dialog, fmt.Sprintf("Failed to set default fee: %v", err))

			return
		}
		c.view.Dialog.Close()
		if afterSave != nil {
			afterSave()
		}
	}

	onCancel := func() { c.view.Dialog.Close() }

	c.view.ConnectSignals(map[string]any{
		"on_ok":     onOk,
		"on_cancel": onCancel,
	})

	gtkutil.RunDialog(c.view.Dialog)
}
