package main

import (
	"errors"

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
func isApplicationWindow(obj glib.IObject) (*gtk.ApplicationWindow, error) {
	if win, ok := obj.(*gtk.ApplicationWindow); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.ApplicationNew")
}

func isDialog(obj glib.IObject) (*gtk.Dialog, error) {
	if dlg, ok := obj.(*gtk.Dialog); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.Dialog")
}

func isAboutDialog(obj glib.IObject) (*gtk.AboutDialog, error) {
	if dlg, ok := obj.(*gtk.AboutDialog); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.AboutDialog")
}

func isEntry(obj glib.IObject) (*gtk.Entry, error) {
	if dlg, ok := obj.(*gtk.Entry); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.Entry")
}

func isTreeView(obj glib.IObject) (*gtk.TreeView, error) {
	if dlg, ok := obj.(*gtk.TreeView); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.TreeView")
}

func isBox(obj glib.IObject) (*gtk.Box, error) {
	if dlg, ok := obj.(*gtk.Box); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.Box")
}
