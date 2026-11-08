//go:build gtk

package controller

import (
	"fmt"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TxUnbondDialogController struct {
	view  *view.TxUnbondDialogView
	model *model.WalletModel

	// delegateOwner is the stake owner that signs the unbond of a delegated
	// validator (PIP-49). It is empty for a validator that is not delegated.
	delegateOwner string

	// seq numbers the delegation reads, so a late answer for a validator the
	// user no longer has selected is ignored.
	seq uint64
}

func NewTxUnbondDialogController(
	view *view.TxUnbondDialogView,
	model *model.WalletModel,
) *TxUnbondDialogController {
	return &TxUnbondDialogController{view: view, model: model}
}

func (c *TxUnbondDialogController) Show() {
	c.applyDefaults()
	c.populateCombos()

	gtkutil.ComboBoxOnChanged(c.view.ValidatorCombo, c.onValidatorChanged)
	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	gtkutil.ComboBoxSelectFirst(c.view.ValidatorCombo)
	c.onValidatorChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func (*TxUnbondDialogController) applyDefaults() {
}

func (c *TxUnbondDialogController) populateCombos() {
	gtkutil.ComboBoxFromAddressList(c.view.ValidatorCombo, c.model.ListValidatorAddresses())
}

func (c *TxUnbondDialogController) onValidatorChanged() {
	validator := gtkutil.ComboBoxGetSelectedAddress(c.view.ValidatorCombo)
	setHint(c.view.ValidatorHint, validatorHint(c.model, validator))

	// Send stays off until the node says whether the validator is delegated.
	c.seq++
	seq := c.seq
	c.delegateOwner = ""
	c.view.ButtonSend.SetSensitive(false)
	loadDelegatedValidator(c.model, validator, "unbond", func(state *delegatedValidator) {
		if seq != c.seq {
			return // the user picked another validator meanwhile
		}
		c.delegateOwner = state.owner
		showDelegationNotice(c.view.DelegationHint, state)
		if state.delegated && state.canSign && !state.expired {
			// An unbond of a delegated validator is valid only after the delegation ends.
			c.view.DelegationHint.SetMarkup(c.view.DelegationHint.Label() + "\n" + fmt.Sprintf(
				"<b>You can unbond from block %d.</b>", state.expiry))
		}
		c.view.ButtonSend.SetSensitive(!state.delegated || (state.canSign && state.expired))
	})
}

func (c *TxUnbondDialogController) onSend() {
	validator := gtkutil.ComboBoxGetSelectedAddress(c.view.ValidatorCombo)
	memo := gtkutil.EntryGetText(c.view.MemoEntry)

	trx, err := c.model.MakeUnbondTx(validator, c.delegateOwner, memo)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}
	if err := verifyUnbondTx(trx, validator, c.delegateOwner); err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	signer := ""
	if c.delegateOwner != "" {
		signer = fmt.Sprintf("\nStake owner: %s (signs this transaction)", c.delegateOwner)
	}

	msg := fmt.Sprintf(`
<b>Transaction Details</b>
<tt>
Type:      Unbond
Validator: %s%s
Fee:       %s
Memo:      %s
</tt>

You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`, validator, signer, trx.Fee(),
		gtkutil.EscapeMarkup(trx.Memo()))

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxUnbondDialogController) onCancel() {
	c.view.Window.Close()
}

// verifyUnbondTx checks that the node built the unbond the user asked for.
func verifyUnbondTx(trx *tx.Tx, validator, delegateOwner string) error {
	pld, ok := trx.Payload().(*payload.UnbondPayload)
	if !ok || pld.Validator.String() != validator {
		return errTxMismatch
	}
	if delegateOwner == "" {
		if pld.IsDelegated() {
			return errTxMismatch
		}

		return nil
	}
	if !pld.IsDelegated() || pld.DelegateOwner.String() != delegateOwner {
		return errTxMismatch
	}

	return nil
}
