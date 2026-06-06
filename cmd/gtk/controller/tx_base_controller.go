//go:build gtk

package controller

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/types/tx"
)

func confirmAndSend(parent *gtk.Window, model *model.WalletModel,
	msg string, trx *tx.Tx,
) {
	gtkutil.ShowQuestionDialog(parent, msg, func(res gtk.ResponseType) {
		if res != gtk.ResponseYes {
			return
		}

		// Proceed with the transaction
		PasswordProvider(model, func(password string, ok bool) {
			if !ok {
				return
			}

			if err := model.SignTransaction(password, trx); err != nil {
				gtkutil.ShowErrorDialog(parent, err.Error(), nil)

				return
			}

			txID, err := model.BroadcastTransaction(trx)
			if err != nil {
				gtkutil.ShowErrorDialog(parent, err.Error(), nil)

				return
			}

			sentMsg := fmt.Sprintf("✅ Transaction sent successfully!\n\n"+
				"Transaction ID: <a href=\"https://pactusscan.com/transaction/%s\">%s</a>", txID, txID)
			gtkutil.ShowInfoDialog(parent, sentMsg, func(gtk.ResponseType) {
				parent.Close()
			})
		})
	})
}

func setDefaultFee(model *model.WalletModel, entry *gtk.Entry) {
	if info, err := model.WalletInfo(); err == nil {
		entry.SetText(fmt.Sprintf("%g", info.DefaultFee.ToPAC()))
	}
}

func setHint(lbl *gtk.Label, hint string) {
	if hint == "" {
		lbl.SetMarkup("")

		return
	}
	lbl.SetMarkup(gtkutil.SmallGray(hint))
}
