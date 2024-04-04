//go:build gtk

package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"runtime"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet"
)

// https://stackoverflow.com/questions/3249053/copying-the-text-from-a-gtk-messagedialog
func updateMessageDialog(dlg *gtk.MessageDialog) {
	area, err := dlg.GetMessageArea()
	if err == nil {
		children := area.GetChildren()
		children.Foreach(func(item interface{}) {
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
		gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, msg)
	updateMessageDialog(dlg)
	res := dlg.Run()
	dlg.Destroy()

	return res == gtk.RESPONSE_YES
}

func showInfoDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, msg)
	updateMessageDialog(dlg)
	dlg.Run()
	dlg.Destroy()
}

func showWarningDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, msg)
	updateMessageDialog(dlg)
	dlg.Run()
	dlg.Destroy()
}

func showErrorDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent,
		gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, msg)
	updateMessageDialog(dlg)
	dlg.Run()
	dlg.Destroy()
}

func errorCheck(err error) {
	if err != nil {
		showErrorDialog(nil, err.Error())
		log.Print(err.Error())
	}
}

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

func updateValidatorHint(lbl *gtk.Label, addr string, w *wallet.Wallet) {
	stake, _ := w.Stake(addr)
	hint := fmt.Sprintf("stake: %s", stake)

	info := w.AddressInfo(addr)
	if info != nil && info.Label != "" {
		hint += ", label: " + info.Label
	}
	updateHintLabel(lbl, hint)
}

func updateAccountHint(lbl *gtk.Label, addr string, w *wallet.Wallet) {
	balance, _ := w.Balance(addr)
	hint := fmt.Sprintf("balance: %s", balance)

	info := w.AddressInfo(addr)
	if info != nil && info.Label != "" {
		hint += ", label: " + info.Label
	}
	updateHintLabel(lbl, hint)
}

func updateFeeHint(lbl *gtk.Label, amtStr string, w *wallet.Wallet, payloadType payload.Type) {
	amt, err := amount.FromString(amtStr)
	if err != nil {
		updateHintLabel(lbl, "")
	} else {
		fee, _ := w.CalculateFee(amt, payloadType)
		hint := fmt.Sprintf("payable: %s, fee: %s",
			fee+amt, fee)
		updateHintLabel(lbl, hint)
	}
}

func updateHintLabel(lbl *gtk.Label, hint string) {
	lbl.SetMarkup(
		fmt.Sprintf("<span foreground='gray' size='small'>%s</span>", hint))
}

func signAndBroadcastTransaction(parent *gtk.Dialog, msg string, w *wallet.Wallet, trx *tx.Tx) {
	if showQuestionDialog(parent, msg) {
		password, ok := getWalletPassword(w)
		if !ok {
			return
		}
		err := w.SignTransaction(password, trx)
		if err != nil {
			errorCheck(err)

			return
		}
		txID, err := w.BroadcastTransaction(trx)
		if err != nil {
			errorCheck(err)

			return
		}

		err = w.Save()
		if err != nil {
			errorCheck(err)

			return
		}

		showInfoDialog(parent,
			fmt.Sprintf("Transaction Hash: <a href=\"https://pacviewer.com/transactions/%s\">%s</a>", txID, txID))
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

	return exec.Command(cmd, args...).Start()
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
