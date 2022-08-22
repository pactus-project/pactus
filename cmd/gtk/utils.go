//go:build gtk

package main

import (
	"log"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func showInfoDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", msg)
	dlg.Run()
	dlg.Destroy()
}

func showWarningDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_WARNING, gtk.BUTTONS_OK, "%s", msg)
	dlg.Run()
	dlg.Destroy()
}

func showErrorDialog(parent gtk.IWindow, msg string) {
	dlg := gtk.MessageDialogNew(parent, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "%s", msg)
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
