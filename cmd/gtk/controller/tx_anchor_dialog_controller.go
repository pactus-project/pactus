//go:build gtk

package controller

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

// TxAnchorDialogController drives the anchor dialog (PIP-50), which can be
// prefilled from an anchor link (PIP-57).
type TxAnchorDialogController struct {
	view  *view.TxAnchorDialogView
	model *model.WalletModel

	accounts []*pactus.AddressInfo
	balance  amount.Amount
	current  model.CurrentAnchor
	anchor   *pactus.AnchorInfo

	// linkTarget is the total deposit a link asks for. When it is set, the
	// deposit to add is computed again for every account the user picks.
	linkTarget *amount.Amount
}

func NewTxAnchorDialogController(
	view *view.TxAnchorDialogView,
	model *model.WalletModel,
) *TxAnchorDialogController {
	return &TxAnchorDialogController{view: view, model: model}
}

func (c *TxAnchorDialogController) Show() {
	c.accounts = c.model.ListAccountAddresses()
	if len(c.accounts) == 0 {
		gtkutil.ShowErrorDialog(gtkutil.MainWindow(),
			"This wallet has no account address. Create one from Wallet > New Address first.", nil)

		return
	}

	setDefaultFee(c.model, c.view.FeeEntry)
	gtkutil.DropDownFromAddressListWithAmount(c.view.SenderDrop, c.accounts, c.model.Balance)
	gtkutil.DropDownSetup(c.view.TypeDrop, anchorTypeItems())
	c.view.TypeDrop.SetSelected(anchorTypeIndex(0x01))
	c.view.TypeSpin.SetValue(1)
	c.view.TypeSpin.SetVisible(false)

	gtkutil.DropDownOnChanged(c.view.SenderDrop, c.onSenderChanged)
	gtkutil.DropDownOnChanged(c.view.TypeDrop, c.onTypeChanged)
	c.view.TypeSpin.ConnectValueChanged(c.refresh)
	c.view.SetToggle.ConnectToggled(c.refresh)
	c.view.DeleteToggle.ConnectToggled(c.refresh)
	gtkutil.EntryOnChanged(c.view.RootEntry, c.refresh)
	gtkutil.EntryOnChanged(c.view.URIEntry, c.refresh)
	gtkutil.EntryOnChanged(c.view.DepositEntry, c.refresh)
	gtkutil.EntryOnChanged(c.view.FeeEntry, c.refresh)
	c.view.LinkEntry.ConnectActivate(c.onApplyLink)

	gtkutil.ConnectButtonSignal(c.view.LinkApplyButton, c.onApplyLink)
	gtkutil.ConnectButtonSignal(c.view.ButtonSend, c.onSend)
	gtkutil.ConnectButtonSignal(c.view.ButtonCancel, c.onCancel)

	c.view.SenderDrop.SetSelected(0)
	c.onSenderChanged()
	c.onTypeChanged()

	gtkutil.ShowNonModalWindow(c.view.Window)
}

func anchorTypeItems() []string {
	items := make([]string, 0, len(model.KnownAnchorTypes)+1)
	for _, typ := range model.KnownAnchorTypes {
		items = append(items, model.AnchorTypeLabel(typ))
	}

	return append(items, "Other (enter the number)")
}

// anchorTypeIndex returns the drop-down index of a known type, or the index of
// the "Other" item.
func anchorTypeIndex(anchorType uint8) uint {
	for i, typ := range model.KnownAnchorTypes {
		if typ == anchorType {
			return uint(i)
		}
	}

	return uint(len(model.KnownAnchorTypes))
}

func (c *TxAnchorDialogController) isDelete() bool {
	return c.view.DeleteToggle.Active()
}

func (c *TxAnchorDialogController) selectedType() uint8 {
	idx := c.view.TypeDrop.Selected()
	if idx < uint(len(model.KnownAnchorTypes)) {
		return model.KnownAnchorTypes[idx]
	}

	return uint8(c.view.TypeSpin.ValueAsInt())
}

func (c *TxAnchorDialogController) sender() string {
	return gtkutil.DropDownGetSelectedAddress(c.view.SenderDrop)
}

func (c *TxAnchorDialogController) onSenderChanged() {
	sender := c.sender()

	hint := ""
	if info := c.model.AddressInfo(sender); info != nil && info.Label != "" {
		hint = gtkutil.EscapeMarkup("label: "+info.Label) + "   "
	}
	balance, err := c.model.Balance(sender)
	if err != nil {
		balance = 0
	}
	c.balance = balance
	hint += gtkutil.LabelValueMarkup("Account Balance", balance.String())
	c.view.SenderHint.SetMarkup(hint)

	c.anchor = nil
	c.current = model.CurrentAnchor{}
	anchor, found, err := c.model.Anchor(sender)
	if err != nil {
		c.view.NoAnchorLabel.SetText("Unable to read the anchor of this account: " + err.Error())
	} else {
		c.view.NoAnchorLabel.SetText("This account has no anchor yet.")
	}
	if found && anchor != nil {
		c.anchor = anchor
		c.current = model.CurrentAnchor{
			Exists: true,
			Type:   uint8(anchor.AnchorType),
			Locked: amount.Amount(anchor.LockedDeposit),
		}
	}
	c.showCurrentAnchor()

	c.view.DepositEntry.SetText(model.DefaultAnchorDeposit(c.current, c.linkTarget).Format())
	c.refresh()
}

func (c *TxAnchorDialogController) showCurrentAnchor() {
	hasAnchor := c.anchor != nil
	c.view.NoAnchorLabel.SetVisible(!hasAnchor)
	c.view.CurrentGrid.SetVisible(hasAnchor)
	if !hasAnchor {
		return
	}

	uri := c.anchor.ManifestUri
	if uri == "" {
		uri = "(none)"
	}
	c.view.CurrentType.SetText(model.AnchorTypeLabel(uint8(c.anchor.AnchorType)))
	c.view.CurrentRoot.SetText(hex.EncodeToString(c.anchor.RootHash))
	c.view.CurrentURI.SetText(uri)
	c.view.CurrentLocked.SetText(amount.Amount(c.anchor.LockedDeposit).String())
	c.view.CurrentUpdate.SetText(fmt.Sprintf("block %d, %s", c.anchor.UpdatedAtHeight,
		time.Unix(int64(c.anchor.UpdatedAtTime), 0).Format("02 Jan 2006 15:04")))
}

func (c *TxAnchorDialogController) onTypeChanged() {
	isOther := c.view.TypeDrop.Selected() >= uint(len(model.KnownAnchorTypes))
	c.view.TypeSpin.SetVisible(isOther)
	c.refresh()
}

// refresh shows the fields of the chosen action, checks them and says what the
// transaction will do. Send is enabled only when the transaction is valid.
func (c *TxAnchorDialogController) refresh() {
	deleting := c.isDelete()
	c.view.SetBox.SetVisible(!deleting)
	c.view.DeleteInfo.SetVisible(deleting)

	fee, feeErr := amount.FromString(gtkutil.EntryGetText(c.view.FeeEntry))
	if feeErr != nil {
		setHint(c.view.FeeHint, "Enter a valid fee.")
	} else {
		setHint(c.view.FeeHint, "")
	}

	var warnings []string
	valid := feeErr == nil
	if deleting {
		valid = c.refreshDelete(fee, valid, &warnings)
	} else {
		valid = c.refreshSet(fee, valid, &warnings)
	}

	c.showWarnings(warnings)
	c.view.ButtonSend.SetSensitive(valid)
}

func (c *TxAnchorDialogController) refreshDelete(fee amount.Amount, valid bool, warnings *[]string) bool {
	check := model.CheckAnchorDelete(c.current, fee)
	if check.Err != nil {
		c.view.DeleteInfo.SetMarkup(gtkutil.EscapeMarkup(upperFirst(check.Err.Error()) + "."))

		return false
	}

	c.view.DeleteInfo.SetMarkup(fmt.Sprintf(
		"The anchor shown above will be removed from the chain.\n"+
			"Returned to the balance: <b>%s</b> (locked deposit %s minus the fee %s).",
		check.Refund, c.current.Locked, fee))
	if check.Refund == 0 {
		*warnings = append(*warnings, "Nothing comes back: the whole locked deposit pays the fee.")
	}

	return valid
}

func (c *TxAnchorDialogController) refreshSet(fee amount.Amount, valid bool, warnings *[]string) bool {
	anchorType := c.selectedType()
	if anchorType >= 0x04 && anchorType <= 0xFE {
		setHint(c.view.TypeHint, fmt.Sprintf("Type %d is reserved by PIP-50; apps treat it as unknown.", anchorType))
	} else {
		setHint(c.view.TypeHint, model.AnchorTypeName(anchorType))
	}

	if _, err := parseRootHash(gtkutil.EntryGetText(c.view.RootEntry)); err != nil {
		setHint(c.view.RootHint, upperFirst(err.Error())+".")
		valid = false
	} else {
		setHint(c.view.RootHint, fmt.Sprintf("%d bytes", len(strings.TrimSpace(
			gtkutil.EntryGetText(c.view.RootEntry)))/2))
	}

	uriLen := len(gtkutil.EntryGetText(c.view.URIEntry))
	setHint(c.view.URIHint, fmt.Sprintf("%d of 128 bytes", uriLen))
	if uriLen > 128 {
		valid = false
	}

	deposit, err := amount.FromString(gtkutil.EntryGetText(c.view.DepositEntry))
	if err != nil {
		setHint(c.view.DepositHint, "Enter a valid deposit.")

		return false
	}

	check := model.CheckAnchorSet(c.current, anchorType, deposit, fee, c.balance)
	hint := fmt.Sprintf("Locked now: %s. Locked after this transaction: %s.",
		c.current.Locked, check.LockedAfter)
	if !c.current.Exists {
		hint += " A new anchor needs at least 1 PAC."
	}
	setHint(c.view.DepositHint, hint)

	if check.Err != nil {
		*warnings = append(*warnings, upperFirst(check.Err.Error())+".")
		valid = false
	}
	if check.TypeChanges {
		*warnings = append(*warnings, fmt.Sprintf(
			"<b>This replaces a %s anchor with a %s anchor.</b> "+
				"An account has a single anchor slot, so the current one will no longer be on chain.",
			gtkutil.EscapeMarkup(model.AnchorTypeName(c.current.Type)),
			gtkutil.EscapeMarkup(model.AnchorTypeName(anchorType))))
	} else if check.Replaces {
		*warnings = append(*warnings, "The current anchor of this account will be replaced.")
	}
	if check.LargeDeposit {
		*warnings = append(*warnings, fmt.Sprintf(
			"This locks <b>%s</b>. It comes back only when you delete the anchor.", deposit))
	}

	return valid
}

// showWarnings shows the warnings in the banner, or hides it when there are none.
// Each warning is markup the controller built, with untrusted values escaped.
func (c *TxAnchorDialogController) showWarnings(warnings []string) {
	if len(warnings) == 0 {
		c.view.WarningLabel.SetVisible(false)

		return
	}
	c.view.WarningLabel.SetMarkup(strings.Join(warnings, "\n"))
	c.view.WarningLabel.SetVisible(true)
}

// onApplyLink checks an anchor link and fills in the form from it.
// The user still reviews every field before anything is signed.
func (c *TxAnchorDialogController) onApplyLink() {
	link := strings.TrimSpace(gtkutil.EntryGetText(c.view.LinkEntry))
	if link == "" {
		c.view.LinkHint.SetMarkup("")

		return
	}

	req, err := model.ParseAnchorRequest(link, c.model.WalletNetwork())
	if err != nil {
		c.view.LinkHint.SetMarkup(fmt.Sprintf("<span foreground='#e0483a'>%s</span>",
			gtkutil.EscapeMarkup(upperFirst(err.Error())+".")))

		return
	}

	c.view.MemoEntry.SetText(req.Memo)
	if req.IsDelete() {
		c.linkTarget = nil
		c.view.DeleteToggle.SetActive(true)
	} else {
		target := req.Locked
		c.linkTarget = &target
		c.view.SetToggle.SetActive(true)
		c.view.TypeDrop.SetSelected(anchorTypeIndex(req.AnchorType))
		c.view.TypeSpin.SetValue(float64(req.AnchorType))
		c.view.RootEntry.SetText(hex.EncodeToString(req.RootHash))
		c.view.URIEntry.SetText(req.ManifestURI)
	}

	note := "Link applied. Check every field before you send."
	if req.Sender != "" {
		if idx, ok := c.accountIndex(req.Sender); ok {
			c.view.SenderDrop.SetSelected(idx)
			note += " The account the link suggests is selected; you can pick another one."
		} else {
			note += " The account the link suggests is not in this wallet; pick an account."
		}
	}
	// Selecting the same account again does not fire a change, so update the
	// deposit for the selected account here.
	c.onSenderChanged()
	c.view.LinkHint.SetMarkup(gtkutil.SmallGray(note))
}

func (c *TxAnchorDialogController) accountIndex(addr string) (uint, bool) {
	for i, info := range c.accounts {
		if info.Address == addr {
			return uint(i), true
		}
	}

	return 0, false
}

func (c *TxAnchorDialogController) onSend() {
	sender := c.sender()
	memo := gtkutil.EntryGetText(c.view.MemoEntry)

	fee, err := amount.FromString(gtkutil.EntryGetText(c.view.FeeEntry))
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	if c.isDelete() {
		c.sendDelete(sender, fee, memo)

		return
	}
	c.sendSet(sender, fee, memo)
}

func (c *TxAnchorDialogController) sendSet(sender string, fee amount.Amount, memo string) {
	anchorType := c.selectedType()
	manifestURI := gtkutil.EntryGetText(c.view.URIEntry)

	rootHash, err := parseRootHash(gtkutil.EntryGetText(c.view.RootEntry))
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, upperFirst(err.Error())+".", nil)

		return
	}
	deposit, err := amount.FromString(gtkutil.EntryGetText(c.view.DepositEntry))
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}
	check := model.CheckAnchorSet(c.current, anchorType, deposit, fee, c.balance)
	if check.Err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, upperFirst(check.Err.Error())+".", nil)

		return
	}

	trx, err := c.model.MakeAnchorSetTx(sender, anchorType, rootHash, manifestURI, deposit, fee, memo)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}
	want := &payload.AnchorPayload{
		Action: payload.AnchorActionSet, RootHash: rootHash, ManifestURI: manifestURI,
		AnchorType: anchorType, Deposit: deposit,
	}
	if err := verifyAnchorTx(trx, sender, want, memo); err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	uri := manifestURI
	if uri == "" {
		uri = "(none)"
	}
	depositText := deposit.String()
	if check.LargeDeposit {
		depositText = "<b>" + depositText + "</b>"
	}
	msg := fmt.Sprintf(`
<b>Transaction Details</b>
<tt>
Type:         Anchor (set)
Network:      %s
Account:      %s
Anchor type:  %s
Root hash:    %s
Manifest URI: %s
Locked now:   %s
Deposit:      %s
Locked after: %s
Fee:          %s
Memo:         %s
</tt>
%s
You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`,
		c.model.WalletNetwork(), sender,
		gtkutil.EscapeMarkup(model.AnchorTypeLabel(anchorType)),
		wrapHex(hex.EncodeToString(rootHash), 64, "              "),
		gtkutil.EscapeMarkup(uri),
		c.current.Locked, depositText, check.LockedAfter, trx.Fee(),
		gtkutil.EscapeMarkup(trx.Memo()),
		confirmWarnings(check))

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxAnchorDialogController) sendDelete(sender string, fee amount.Amount, memo string) {
	check := model.CheckAnchorDelete(c.current, fee)
	if check.Err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, upperFirst(check.Err.Error())+".", nil)

		return
	}

	trx, err := c.model.MakeAnchorDeleteTx(sender, fee, memo)
	if err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}
	if err := verifyAnchorTx(trx, sender, &payload.AnchorPayload{Action: payload.AnchorActionDelete}, memo); err != nil {
		gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

		return
	}

	uri := c.anchor.ManifestUri
	if uri == "" {
		uri = "(none)"
	}
	msg := fmt.Sprintf(`
<b>Transaction Details</b>
<tt>
Type:         Anchor (delete)
Network:      %s
Account:      %s
Removes:
  Anchor type:  %s
  Root hash:    %s
  Manifest URI: %s
Locked:       %s
Fee:          %s
Returned:     <b>%s</b>
Memo:         %s
</tt>

You are going to sign and broadcast this transaction.
<b>This action cannot be undone.</b>
Do you want to continue with this transaction?`,
		c.model.WalletNetwork(), sender,
		gtkutil.EscapeMarkup(model.AnchorTypeLabel(c.current.Type)),
		wrapHex(hex.EncodeToString(c.anchor.RootHash), 64, "                "),
		gtkutil.EscapeMarkup(uri),
		c.current.Locked, trx.Fee(), check.Refund,
		gtkutil.EscapeMarkup(trx.Memo()))

	confirmAndSend(c.view.Window, c.model, msg, trx)
}

func (c *TxAnchorDialogController) onCancel() {
	c.view.Window.Close()
}

// verifyAnchorTx checks that the node built the transaction the user asked for.
// The node may be remote, so the wallet does not sign a transaction it did not check.
func verifyAnchorTx(trx *tx.Tx, sender string, want *payload.AnchorPayload, memo string) error {
	pld, ok := trx.Payload().(*payload.AnchorPayload)
	if !ok || pld.From.String() != sender || pld.Action != want.Action || trx.Memo() != memo {
		return errTxMismatch
	}
	if want.Action == payload.AnchorActionSet &&
		(!bytes.Equal(pld.RootHash, want.RootHash) || pld.ManifestURI != want.ManifestURI ||
			pld.AnchorType != want.AnchorType || pld.Deposit != want.Deposit) {
		return errTxMismatch
	}

	return nil
}

// parseRootHash decodes a root hash typed as hex. It must be 32 to 64 bytes.
func parseRootHash(text string) ([]byte, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("enter the root hash")
	}
	rootHash, err := hex.DecodeString(text)
	if err != nil {
		return nil, errors.New("the root hash is not valid hex")
	}
	if len(rootHash) < 32 || len(rootHash) > 64 {
		return nil, fmt.Errorf("the root hash is %d bytes; it must be 32 to 64 bytes", len(rootHash))
	}

	return rootHash, nil
}

func confirmWarnings(check model.AnchorSetCheck) string {
	var lines []string
	if check.TypeChanges {
		lines = append(lines, "<b>Warning:</b> the current anchor has another type and will be replaced.")
	} else if check.Replaces {
		lines = append(lines, "The current anchor of this account will be replaced.")
	}
	if check.LargeDeposit {
		lines = append(lines, "<b>The deposit is above 1 PAC.</b> It comes back only when you delete the anchor.")
	}
	if len(lines) == 0 {
		return ""
	}

	return "\n" + strings.Join(lines, "\n") + "\n"
}

// wrapHex splits a long hex string into lines of the given width, indenting
// the lines after the first one so they line up in the details table.
func wrapHex(text string, width int, indent string) string {
	if len(text) <= width {
		return text
	}

	var builder strings.Builder
	for start := 0; start < len(text); start += width {
		if start > 0 {
			builder.WriteString("\n" + indent)
		}
		builder.WriteString(text[start:min(start+width, len(text))])
	}

	return builder.String()
}
