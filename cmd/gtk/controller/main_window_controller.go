//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

type MainWindowController struct {
	view *view.MainWindowView
}

func NewMainWindowController(view *view.MainWindowView) *MainWindowController {
	return &MainWindowController{view: view}
}

func (c *MainWindowController) BuildView() {
	c.view.Window.ShowMenubar()
}
