//go111:build gtk

package controller

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/types/tx"
)

func confirmAndSend(parent *gtk.Window, model *model.WalletModel, msg string, trx *tx.Tx) bool {
	if !gtkutil.ShowQuestionDialog(parent, msg) {
		return false
	}

	password, ok := PasswordProvider(model)
	if !ok {
		return false
	}

	if err := model.SignTransaction(password, trx); err != nil {
		gtkutil.ShowError(err)

		return false
	}

	txID, err := model.BroadcastTransaction(trx)
	if err != nil {
		gtkutil.ShowError(err)

		return false
	}

	gtkutil.ShowInfoDialog(parent,
		fmt.Sprintf("✅ Transaction sent successfully!\n\n"+
			"Transaction ID: <a href=\"https://pactusscan.com/transaction/%s\">%s</a>", txID, txID))

	return true
}
