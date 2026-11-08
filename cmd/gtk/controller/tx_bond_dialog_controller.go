//go:build gtk

package controller

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
)

type TxBondDialogController struct {
	view  *view.TxBondDialogView
	model *model.WalletModel

	// chain is the status of the node, read in the background for a delegated
	// bond. It is nil until the first read finishes.
	chain    *model.ChainStatus
	chainErr error
	loading  bool

	// receiverStake is the stake of the selected validator, or 0 when it is
	// not on the blockchain yet.
	receiverStake amount.Amount
}

func NewTxBondDialogController(
	view *view.TxBondDialogView,
	model *model.WalletModel,
) *TxBondDialogController {
	return &TxBondDialogController{view: view, model: model}
}

func (c *TxBondDialogController) Show() {
	c.applyDefaults()
	c.populateCombos()

	gtkutil.DropDownOnChanged(c.view.SenderDrop, c.onSenderChanged)
	gtkutil.ComboBoxOnChanged(c.view.ReceiverCombo, c.onReceiverChanged)
	gtkutil.EntryOnChanged(c.view.FeeEntry, c.onFeeChanged)

	c.view.DelegateCheck.ConnectToggled(c.onDelegateToggled)
	gtkutil.EntryOnChanged(c.view.ShareEntry, c.refreshDelegation)
	gtkutil.EntryOnChanged(c.view.ExpiryEntry, c.refreshDelegation)
	gtkutil.DropDownOnChanged(c.view.TermDrop, c.onTermChanged)

	gtkutil.ConnectButtonSignal(c.view.AmountMaxButton, c.onMax)
	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	c.view.SenderDrop.SetSelected(0)
	c.onSenderChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func (c *TxBondDialogController) applyDefaults() {
	setDefaultFee(c.model, c.view.FeeEntry)

	terms := make([]string, 0, len(model.DelegationTerms)+1)
	for _, term := range model.DelegationTerms {
		terms = append(terms, term.Label)
	}
	gtkutil.DropDownSetup(c.view.TermDrop, append(terms, "Until a block height"))
	c.view.TermDrop.SetSelected(1)
}

func (c *TxBondDialogController) populateCombos() {
	gtkutil.DropDownFromAddressListWithAmount(c.view.SenderDrop, c.model.ListAccountAddresses(), c.model.Balance)
	gtkutil.ComboBoxFromAddressListWithAmount(c.view.ReceiverCombo, c.model.ListValidatorAddresses(), c.model.Stake)
}

func (c *TxBondDialogController) onSenderChanged() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	if info := c.model.AddressInfo(sender); info != nil && info.Label != "" {
		setHint(c.view.SenderHint, fmt.Sprintf("label: %s", info.Label))
	} else {
		setHint(c.view.SenderHint, "")
	}

	bal, err := c.model.Balance(sender)
	if err == nil {
		setBalanceHint(c.view.AmountHint, bal)
	} else {
		setHint(c.view.AmountHint, "")
	}

	c.refreshDelegation()
}

func (c *TxBondDialogController) onReceiverChanged() {
	hint := ""

	receiver := gtkutil.ComboBoxGetSelectedAddress(c.view.ReceiverCombo)
	stake, err := c.model.Stake(receiver)
	c.receiverStake = 0
	if err == nil {
		c.receiverStake = stake
		hint = fmt.Sprintf("stake: %s", stake)
	}
	if info := c.model.AddressInfo(receiver); info != nil && info.Label != "" {
		if hint != "" {
			hint += ", "
		}
		hint += "label: " + info.Label
	}
	setHint(c.view.ReceiverHint, hint)

	c.refreshDelegation()
}

// onMax fills in the whole balance of the sender, minus the fee, but not more
// than the stake the validator can still receive.
func (c *TxBondDialogController) onMax() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	balance, err := c.model.Balance(sender)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	fee, err := amount.FromString(gtkutil.EntryGetText(c.view.FeeEntry))
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	// A validator that is not on the blockchain yet has no stake.
	receiver := gtkutil.ComboBoxGetSelectedAddress(c.view.ReceiverCombo)
	stake, err := c.model.Stake(receiver)
	if err != nil {
		stake = 0
	}

	maxStake := genesis.DefaultGenesisParams().MaximumStake
	maxAmount, err := model.MaxBondAmount(balance, fee, stake, maxStake)
	setMaxAmount(c.view.Window, c.view.AmountEntry, maxAmount, err)
}

func (c *TxBondDialogController) isDelegated() bool {
	return c.view.DelegateCheck.Active()
}

// onDelegateToggled switches between a normal bond and a delegated bond. A
// delegated bond always stakes the maximum stake, so the amount is fixed.
func (c *TxBondDialogController) onDelegateToggled() {
	delegated := c.isDelegated()
	c.view.DelegateBox.SetVisible(delegated)
	c.view.AmountEntry.SetSensitive(!delegated)
	c.view.AmountMaxButton.SetSensitive(!delegated)
	if delegated {
		c.view.AmountEntry.SetText(genesis.DefaultGenesisParams().MaximumStake.Format())
		c.loadChainStatus()
	}
	c.onReceiverChanged()
}

func (c *TxBondDialogController) onTermChanged() {
	custom := c.view.TermDrop.Selected() >= uint(len(model.DelegationTerms))
	c.view.ExpiryEntry.SetVisible(custom)
	c.loadChainStatus()
	c.refreshDelegation()
}

// loadChainStatus reads the height of the node in the background. The node can
// take a while to answer when it is busy, for example while it syncs, so the
// dialog must not wait for it on the UI thread.
func (c *TxBondDialogController) loadChainStatus() {
	if c.loading {
		return
	}
	c.loading = true

	go func() {
		status, err := c.model.ChainStatus()
		gtkutil.IdleAddAsync(func() {
			c.loading = false
			c.chainErr = err
			if err == nil {
				c.chain = &status
			}
			c.refreshDelegation()
		})
	}()
}

// delegationValues reads the owner share and the expiry height from the form.
func (c *TxBondDialogController) delegationValues(height uint32) (amount.Amount, uint32, error) {
	percent, err := strconv.ParseFloat(strings.TrimSpace(gtkutil.EntryGetText(c.view.ShareEntry)), 64)
	if err != nil {
		return 0, 0, model.ErrDelegateShareRange
	}
	share, err := model.DelegateShareFromPercent(percent)
	if err != nil {
		return 0, 0, err
	}

	idx := c.view.TermDrop.Selected()
	var expiry uint32
	if idx < uint(len(model.DelegationTerms)) {
		expiry = model.DelegationExpiryAfterDays(height, model.DelegationTerms[idx].Days)
	} else {
		value, err := strconv.ParseUint(strings.TrimSpace(gtkutil.EntryGetText(c.view.ExpiryEntry)), 10, 32)
		if err != nil {
			return share, 0, model.ErrDelegateExpiryPast
		}
		expiry = uint32(value)
	}

	// The bond is in the next block at the earliest, and the node rejects an
	// expiry that is not after the height of that block.
	return share, expiry, model.CheckDelegationExpiry(height+1, expiry)
}

// refreshDelegation explains the delegation the form describes and enables
// Send only when it is valid.
func (c *TxBondDialogController) refreshDelegation() {
	if !c.isDelegated() {
		c.view.ButtonSend.SetSensitive(true)

		return
	}

	owner := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	c.view.OwnerLabel.SetMarkup(gtkutil.LabelValueMarkup("Stake owner (the selected account)", owner))

	valid := true
	height := uint32(0)
	if c.chain != nil {
		height = c.chain.Height
	}
	share, expiry, err := c.delegationValues(height)
	shareErr := errors.Is(err, model.ErrDelegateShareRange)
	if shareErr {
		setHint(c.view.ShareHint, upperFirst(err.Error())+".")
		valid = false
	} else {
		ownerReward, operatorReward := model.DelegationRewardSplit(share)
		setHint(c.view.ShareHint, fmt.Sprintf(
			"Of each 1 PAC block reward, you receive %s and the operator %s; 0.3 PAC goes to the treasury.",
			ownerReward, operatorReward))
	}

	switch {
	case c.chain == nil && c.chainErr != nil:
		setHint(c.view.ExpiryHint, "Unable to read the current block height: "+c.chainErr.Error())
		valid = false
	case c.chain == nil:
		setHint(c.view.ExpiryHint, "Reading the current block height from the node...")
		valid = false
	case !c.chain.IsSynced():
		setHint(c.view.ExpiryHint, fmt.Sprintf(
			"The node is still syncing (%d blocks left). The delegation end is set from the "+
				"current block, so wait until the node is synced.", c.chain.BlocksLeft))
		valid = false
	case err != nil && !shareErr:
		setHint(c.view.ExpiryHint, fmt.Sprintf("%s (current block %d).", upperFirst(err.Error()), height))
		valid = false
	default:
		ends := model.EstimateBlockTime(height, expiry, time.Now())
		setHint(c.view.ExpiryHint, fmt.Sprintf(
			"Ends at block %d, around %s. You can unbond only after that.", expiry, ends.Format("02 Jan 2006")))
	}

	if err := model.CheckDelegatedBondValidator(c.receiverStake); err != nil {
		setHint(c.view.ReceiverHint, upperFirst(err.Error())+".")
		valid = false
	}

	c.view.ButtonSend.SetSensitive(valid)
}

func (c *TxBondDialogController) onFeeChanged() {
	_ = payload.TypeBond
	setHint(c.view.FeeHint, "")
}

func (c *TxBondDialogController) onSend() {
	sender := gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
	receiver := gtkutil.ComboBoxGetSelectedAddress(c.view.ReceiverCombo)
	publicKey := gtkutil.EntryGetText(c.view.PublicKeyEntry)
	amountStr := gtkutil.EntryGetText(c.view.AmountEntry)
	feeStr := gtkutil.EntryGetText(c.view.FeeEntry)
	memo := gtkutil.EntryGetText(c.view.MemoEntry)

	amt, err := amount.FromString(amountStr)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	fee, err := amount.FromString(feeStr)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	if c.isDelegated() {
		c.sendDelegated(sender, receiver, publicKey, amt, fee, memo)

		return
	}

	trx, err := c.model.MakeBondTx(sender, receiver, publicKey, amt, fee, memo, nil)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	msg := fmt.Sprintf(`
<b>Transaction Details</b>
<tt>
Type:   Bond
From:   %s
To:     %s
Amount: %s
Fee:    %s
Memo:   %s
</tt>

You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`, sender, receiver, amt, trx.Fee(),
		gtkutil.EscapeMarkup(trx.Memo()))

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

// sendDelegated builds a delegated bond (PIP-49) with the selected account as
// the stake owner, checks what the node built, and asks the user to confirm.
func (c *TxBondDialogController) sendDelegated(sender, receiver, publicKey string,
	amt, fee amount.Amount, memo string,
) {
	// Send is enabled only once the node status is known and synced.
	if c.chain == nil || !c.chain.IsSynced() {
		return
	}
	height := c.chain.Height
	share, expiry, err := c.delegationValues(height)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, upperFirst(err.Error())+".", nil)

		return
	}
	if err := model.CheckDelegatedBondValidator(c.receiverStake); err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, upperFirst(err.Error())+".", nil)

		return
	}

	delegation := &model.Delegation{Owner: sender, Share: share, Expiry: expiry}
	trx, err := c.model.MakeBondTx(sender, receiver, publicKey, amt, fee, memo, delegation)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}
	if err := verifyDelegatedBond(trx, sender, receiver, amt, delegation); err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	ownerReward, operatorReward := model.DelegationRewardSplit(share)
	ends := model.EstimateBlockTime(height, expiry, time.Now())
	msg := fmt.Sprintf(`
<b>Transaction Details</b>
<tt>
Type:           Bond (delegated, PIP-49)
From:           %s
To:             %s
Amount:         %s
Stake owner:    %s
Owner reward:   %s per block (%.4g%%)
Operator gets:  %s per block
Delegation end: block %d (around %s)
Fee:            %s
Memo:           %s
</tt>

Only the stake owner can unbond this validator, after block %d,
and only the stake owner can withdraw the stake.

You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`,
		sender, receiver, amt, sender,
		ownerReward, model.DelegateSharePercent(share), operatorReward,
		expiry, ends.Format("02 Jan 2006"), trx.Fee(), gtkutil.EscapeMarkup(trx.Memo()), expiry)

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

// verifyDelegatedBond checks that the node built the delegated bond the user
// asked for. The node may be remote, so the wallet does not sign a transaction
// it did not check.
func verifyDelegatedBond(trx *tx.Tx, sender, receiver string, amt amount.Amount,
	delegation *model.Delegation,
) error {
	pld, ok := trx.Payload().(*payload.BondPayload)
	if !ok || pld.From.String() != sender || pld.To.String() != receiver || pld.Stake != amt ||
		!pld.IsDelegated() || pld.DelegateOwner.String() != delegation.Owner ||
		pld.DelegateShare != delegation.Share || uint32(pld.DelegateExpiry) != delegation.Expiry {
		return errTxMismatch
	}

	return nil
}

func (c *TxBondDialogController) onCancel() {
	c.view.Window.Close()
}
