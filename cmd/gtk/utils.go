//go:build gtk

package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet"
)

// https://stackoverflow.com/questions/3249053/copying-the-text-from-a-gtk-messagedialog
func updateMessageDialog(dlg *gtk.MessageDialog) {
	area, err := dlg.GetMessageArea()
	if err == nil {
		children := area.GetChildren()
		children.Foreach(func(item any) {
			label, err := gtk.WidgetToLabel(item.(*gtk.Widget))
			if err == nil {
				label.SetSelectable(true)
				label.SetUseMarkup(true)
			}
		})
	}
}

func showQuestionDialog(parent gtk.IWindow, msg string) bool {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "%s", msg)
	updateMessageDialog(dlg)
	res := runDialog(&dlg.Dialog)

	return res == gtk.RESPONSE_YES
}

func showInfoDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", msg)
	updateMessageDialog(dlg)
	runDialog(&dlg.Dialog)
}

func showWarningDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", msg)
	updateMessageDialog(dlg)
	runDialog(&dlg.Dialog)
}

func showErrorDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "%s", msg)
	updateMessageDialog(dlg)
	runDialog(&dlg.Dialog)
}

// showError displays an error dialog and logs the error message.
func showError(err error) {
	showErrorDialog(nil, err.Error())
	log.Print(err.Error())
}

// fatalErrorCheck checks for an error, shows an error dialog and terminates the program.
// Use with caution.
func fatalErrorCheck(err error) {
	if err != nil {
		showErrorDialog(nil, err.Error())
		log.Fatal(err.Error())
	}
}

func getObj(builder *gtk.Builder, name string) glib.IObject {
	obj, err := builder.GetObject(name)
	fatalErrorCheck(err)

	return obj
}

func getApplicationWindowObj(builder *gtk.Builder, name string) *gtk.ApplicationWindow {
	return getObj(builder, name).(*gtk.ApplicationWindow)
}

func getDialogObj(builder *gtk.Builder, name string) *gtk.Dialog {
	return getObj(builder, name).(*gtk.Dialog)
}

func getAboutDialogObj(builder *gtk.Builder, name string) *gtk.AboutDialog {
	return getObj(builder, name).(*gtk.AboutDialog)
}

func getComboBoxTextObj(builder *gtk.Builder, name string) *gtk.ComboBoxText {
	return getObj(builder, name).(*gtk.ComboBoxText)
}

func getEntryObj(builder *gtk.Builder, name string) *gtk.Entry {
	return getObj(builder, name).(*gtk.Entry)
}

func getOverlayObj(builder *gtk.Builder, name string) *gtk.Overlay {
	return getObj(builder, name).(*gtk.Overlay)
}

func getTreeViewObj(builder *gtk.Builder, name string) *gtk.TreeView {
	return getObj(builder, name).(*gtk.TreeView)
}

func getTextViewObj(builder *gtk.Builder, name string) *gtk.TextView {
	return getObj(builder, name).(*gtk.TextView)
}

func getBoxObj(builder *gtk.Builder, name string) *gtk.Box {
	return getObj(builder, name).(*gtk.Box)
}

func getLabelObj(builder *gtk.Builder, name string) *gtk.Label {
	return getObj(builder, name).(*gtk.Label)
}

func getToolButtonObj(builder *gtk.Builder, name string) *gtk.ToolButton {
	return getObj(builder, name).(*gtk.ToolButton)
}

func getButtonObj(builder *gtk.Builder, name string) *gtk.Button {
	return getObj(builder, name).(*gtk.Button)
}

func getImageObj(builder *gtk.Builder, name string) *gtk.Image {
	return getObj(builder, name).(*gtk.Image)
}

func getProgressBarObj(builder *gtk.Builder, name string) *gtk.ProgressBar {
	return getObj(builder, name).(*gtk.ProgressBar)
}

func getMenuItem(builder *gtk.Builder, name string) *gtk.MenuItem {
	return getObj(builder, name).(*gtk.MenuItem)
}

func getTextViewContent(tv *gtk.TextView) string {
	buf, _ := tv.GetBuffer()
	startIter, endIter := buf.GetBounds()
	content, err := buf.GetText(startIter, endIter, true)
	if err != nil {
		// TODO: Log error
		return ""
	}

	return content
}

func setTextViewContent(tv *gtk.TextView, content string) {
	buf, err := tv.GetBuffer()
	if err != nil {
		// TODO: Log error
		return
	}
	buf.SetText(content)
}

func updateValidatorHint(lbl *gtk.Label, addr string, wlt *wallet.Wallet) {
	stake, _ := wlt.Stake(addr)
	hint := fmt.Sprintf("stake: %s", stake)

	info := wlt.AddressInfo(addr)
	if info != nil && info.Label != "" {
		hint += ", label: " + info.Label
	}
	updateHintLabel(lbl, hint)
}

func updateAccountHint(lbl *gtk.Label, addr string, wlt *wallet.Wallet) {
	info := wlt.AddressInfo(addr)
	if info != nil && info.Label != "" {
		updateHintLabel(lbl, fmt.Sprintf("label: %s", info.Label))
	}
}

func updateFeeHint(_ *gtk.Label, _ *wallet.Wallet, _ payload.Type) {
	// Nothing for now!
	// The goal is to show an estimate of how long it takes for a transaction
	// with the given fee to be confirmed (confirmation time).
	// We can analyze data from past blocks to estimate the confirmation time.
}

func updateBalanceHint(lbl *gtk.Label, addr string, wlt *wallet.Wallet) {
	balance, err := wlt.Balance(addr)
	if err == nil {
		updateHintLabel(lbl, fmt.Sprintf("Account Balance: %s", balance))
	} else {
		updateHintLabel(lbl, "")
	}
}

func updateStakeHint(lbl *gtk.Label, addr string, wlt *wallet.Wallet) {
	stake, err := wlt.Stake(addr)
	if err == nil {
		updateHintLabel(lbl, fmt.Sprintf("Validator Stake: %s", stake))
	} else {
		updateHintLabel(lbl, "")
	}
}

func updateHintLabel(lbl *gtk.Label, hint string) {
	lbl.SetMarkup(
		fmt.Sprintf("<span foreground='gray' size='small'>%s</span>", hint))
}

func signAndBroadcastTransaction(parent *gtk.Dialog, msg string, wlt *wallet.Wallet, trx *tx.Tx) {
	if showQuestionDialog(parent, msg) {
		password, ok := getWalletPassword(wlt)
		if !ok {
			return
		}
		err := wlt.SignTransaction(password, trx)
		if err != nil {
			showError(err)

			return
		}
		txID, err := wlt.BroadcastTransaction(trx)
		if err != nil {
			showError(err)

			return
		}

		err = wlt.Save()
		if err != nil {
			showError(err)

			return
		}

		showInfoDialog(parent,
			fmt.Sprintf("✅ Transaction sent successfully!\n\nTransaction Hash: <a href=\"https://pacviewer.com/transaction/%s\">%s</a>", txID, txID))

		parent.Close()
	}
}

// openURLInBrowser open specific url in browser base on os.
func openURLInBrowser(address string) error {
	var cmd string
	args := make([]string, 0)

	addr, err := url.Parse(address)
	if err != nil {
		return err
	}

	switch addr.Scheme {
	case "http", "https":
	default:
		return errors.New("address scheme is invalid")
	}

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, address)

	return exec.CommandContext(context.Background(), cmd, args...).Start()
}

func buildExtendedEntry(builder *gtk.Builder, overlayID string) *gtk.Entry {
	overlay := getOverlayObj(builder, overlayID)

	// Create a new Entry
	entry, err := gtk.EntryNew()
	fatalErrorCheck(err)
	entry.SetCanFocus(true)
	entry.SetHExpand(true)
	entry.SetEditable(false)

	setCSSClass(&entry.Widget, "copyable_entry")

	// Create a new Button
	button, err := gtk.ButtonNewFromIconName("edit-copy-symbolic", gtk.ICON_SIZE_BUTTON)
	fatalErrorCheck(err)

	button.SetTooltipText("Copy to Clipboard") // TODO: Not working!
	button.SetHAlign(gtk.ALIGN_END)
	button.SetVAlign(gtk.ALIGN_CENTER)
	button.SetHExpand(false)
	button.SetVExpand(false)
	button.SetBorderWidth(0)

	setCSSClass(&button.Widget, "inline_button")

	// Set the click event for the Button
	button.Connect("clicked", func() {
		buffer, _ := entry.GetText()
		clipboard, _ := gtk.ClipboardGet(gdk.SELECTION_CLIPBOARD)
		clipboard.SetText(buffer)
	})

	overlay.Add(entry)
	overlay.AddOverlay(button)

	overlay.ShowAll() // Ensure all child widgets are shown

	return entry
}

func setCSSClass(widget *gtk.Widget, name string) {
	styleContext, err := widget.GetStyleContext()
	fatalErrorCheck(err)

	styleContext.AddClass(name)
}

func runDialog(dlg *gtk.Dialog) gtk.ResponseType {
	response := dlg.Run()

	// Destroy should be done after the dialog is closed
	// Read more here: https://docs.gtk.org/gtk3/method.Dialog.run.html
	dlg.Destroy()

	return response
}
