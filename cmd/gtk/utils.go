package main

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func showInfoDialog(msg string) {
	dlg := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL, gtk.MESSAGE_INFO, gtk.BUTTONS_OK, "%s", msg)
	dlg.Run()
	dlg.Destroy()
}

func showErrorDialog(msg string) {
	dlg := gtk.MessageDialogNew(nil, gtk.DIALOG_MODAL, gtk.MESSAGE_ERROR, gtk.BUTTONS_OK, "%s", msg)
	dlg.Run()
	dlg.Destroy()
}

func errorCheck(err error) {
	if err != nil {
		showErrorDialog(err.Error())
		gtk.MainQuit()
	}
}

func getObj(builder *gtk.Builder, name string) glib.IObject {
	obj, err := builder.GetObject(name)
	errorCheck(err)
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

func getBoxObj(builder *gtk.Builder, name string) *gtk.Box {
	return getObj(builder, name).(*gtk.Box)
}

func getLabelObj(builder *gtk.Builder, name string) *gtk.Label {
	return getObj(builder, name).(*gtk.Label)
}

func getProgressBarObj(builder *gtk.Builder, name string) *gtk.ProgressBar {
	return getObj(builder, name).(*gtk.ProgressBar)
}
