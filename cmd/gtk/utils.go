package main

import (
	"errors"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

func errorCheck(err error) {
	if err != nil {
		showErrorDialog(nil, err.Error())
		gtk.MainQuit()
	}
}
func isApplicationWindow(obj glib.IObject) (*gtk.ApplicationWindow, error) {
	// Make type assertion (as per gtk.go).
	if win, ok := obj.(*gtk.ApplicationWindow); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.ApplicationNew")
}

func isDialog(obj glib.IObject) (*gtk.Dialog, error) {
	// Make type assertion (as per gtk.go).
	if dlg, ok := obj.(*gtk.Dialog); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.Dialog")
}

func isEntry(obj glib.IObject) (*gtk.Entry, error) {
	// Make type assertion (as per gtk.go).
	if dlg, ok := obj.(*gtk.Entry); ok {
		return dlg, nil
	}
	return nil, errors.New("not a *gtk.Entry")
}
