//go:build gtk

package controller

import (
	"fmt"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
)

func confirmAndSend(parent *gtk.Window, model *model.WalletModel,
	msg string, trx *tx.Tx,
) {
	gtkutil.ShowQuestionDialog(parent, msg, func(res int) {
		if res != int(gtk.ResponseYes) {
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
			gtkutil.ShowInfoDialog(parent, sentMsg, func(int) {
				parent.Close()
				gtkutil.PresentMainWindow()
			})
		})
	})
}

func setDefaultFee(model *model.WalletModel, entry *gtk.Entry) {
	if info, err := model.WalletInfo(); err == nil {
		entry.SetText(fmt.Sprintf("%g", amount.Amount(info.DefaultFee).ToPAC()))
	}
}

// setBalanceHint shows the account balance below the amount field, in normal
// size with the amount in bold, so it is easy to see.
func setBalanceHint(lbl *gtk.Label, balance amount.Amount) {
	lbl.SetMarkup(gtkutil.LabelValueMarkup("Account Balance", balance.String()))
}

// setMaxAmount fills the amount field with the result of a Max calculation,
// or shows the error when there is nothing to fill in.
func setMaxAmount(window *gtk.Window, entry *gtk.Entry, maxAmount amount.Amount, err error) {
	if err != nil {
		gtkutil.ShowErrorDialog(window, err.Error(), nil)

		return
	}
	entry.SetText(maxAmount.Format())
}

func setHint(lbl *gtk.Label, hint string) {
	if hint == "" {
		lbl.SetMarkup("")

		return
	}
	lbl.SetMarkup(gtkutil.SmallGray(hint))
}
