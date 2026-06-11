//go:build gtk

package controller

import (
	"fmt"
	"strings"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
)

type WalletDefaultFeeDialogController struct {
	view  *view.WalletDefaultFeeDialogView
	model *model.WalletModel
}

func NewWalletDefaultFeeDialogController(
	view *view.WalletDefaultFeeDialogView,
	model *model.WalletModel,
) *WalletDefaultFeeDialogController {
	return &WalletDefaultFeeDialogController{view: view, model: model}
}

func (c *WalletDefaultFeeDialogController) Run(onUpdate func()) {
	info, err := c.model.WalletInfo()
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, fmt.Sprintf("Failed to get wallet info: %v", err), nil)

		return
	}

	currentFee := amount.Amount(info.DefaultFee)
	c.view.CurrentFeeLabel.SetText(currentFee.String())
	c.view.FeeEntry.SetText(strings.ReplaceAll(currentFee.String(), " PAC", ""))

	onOk := func() {
		feeStr := gtkutil.EntryGetText(c.view.FeeEntry)
		feeAmount, err := amount.FromString(feeStr)
		if err != nil {
			gtkutil.ShowErrorDialog(c.view.Window, fmt.Sprintf("Invalid fee amount: %v", err), nil)

			return
		}
		if err := c.model.SetDefaultFee(feeAmount); err != nil {
			gtkutil.ShowErrorDialog(c.view.Window, fmt.Sprintf("Failed to set default fee: %v", err), nil)

			return
		}
		c.view.Window.Close()
		onUpdate()
	}

	onCancel := func() { c.view.Window.Close() }

	gtkutil.ConnectButtonSignal(c.view.ButtonOK, onOk)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, onCancel)

	gtkutil.ShowModalWindow(c.view.Window)
}
