//go:build gtk

package main

import (
	"fmt"
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/wallet"
)

func showQuestionDialog(parent gtk.IWindow, msg string) bool {
	dlg := gtk.MessageDialogNewWithMarkup(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_QUESTION, gtk.BUTTONS_YES_NO, "")
	dlg.SetMarkup(msg)
	res := dlg.Run()
	dlg.Destroy()
	return res == gtk.RESPONSE_YES
}

func showInfoDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNewWithMarkup(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "")
	dlg.SetMarkup(msg)
	dlg.Run()
	dlg.Destroy()
}

func showWarningDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNewWithMarkup(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "")
	dlg.SetMarkup(msg)
	dlg.Run()
	dlg.Destroy()
}

func showErrorDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNewWithMarkup(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "")
	dlg.SetMarkup(msg)
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
	stake, err := w.Stake(addr)
	if err != nil {
		updateHintLabel(lbl, "")
	} else {
		hint := fmt.Sprintf("stake: %v", util.ChangeToString(stake))

		info := w.AddressInfo(addr)
		if info != nil && info.Label != "" {
			hint += ", label: " + info.Label
		}
		updateHintLabel(lbl, hint)
	}
}
func updateAccountHint(lbl *gtk.Label, addr string, w *wallet.Wallet) {
	balance, err := w.Balance(addr)
	if err != nil {
		updateHintLabel(lbl, "")
	} else {
		hint := fmt.Sprintf("balance: %v", util.ChangeToString(balance))

		info := w.AddressInfo(addr)
		if info != nil && info.Label != "" {
			hint += ", label: " + info.Label
		}
		updateHintLabel(lbl, hint)
	}
}

func updateFeeHint(lbl *gtk.Label, amtStr string, w *wallet.Wallet) {
	amount, err := util.StringToChange(amtStr)
	if err != nil {
		updateHintLabel(lbl, "")
	} else {
		fee := w.CalculateFee(amount)
		hint := fmt.Sprintf("payable: %v, fee: %v",
			util.ChangeToString(fee+amount), util.ChangeToString(fee))
		updateHintLabel(lbl, hint)
	}
}

func updateHintLabel(lbl *gtk.Label, hint string) {
	lbl.SetMarkup(
		fmt.Sprintf("<span foreground='gray' size='small'>%s</span>", hint))
}
