//go:build gtk

package controller

import (
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
)

const configRestartMessage = "Configuration saved successfully.\n\n" +
	"You must restart the node for the changes to take effect."

type ConfigEditorDialogController struct {
	view  *view.ConfigEditorDialogView
	model *model.ConfigModel
}

func NewConfigEditorDialogController(
	view *view.ConfigEditorDialogView,
	configModel *model.ConfigModel,
) *ConfigEditorDialogController {
	return &ConfigEditorDialogController{
		view:  view,
		model: configModel,
	}
}

func (c *ConfigEditorDialogController) Run() {
	c.setEditorContent(c.model.SavedContent())
	c.updateSaveButton()

	buf := c.view.TextView.Buffer()
	if buf != nil {
		buf.Connect("changed", func() {
			c.updateSaveButton()
		})
	}

	onSave := func() {
		content := gtkutil.GetTextViewContent(c.view.TextView)
		if err := c.model.Save(content); err != nil {
			gtkutil.ShowErrorDialog(c.view.Window, err.Error(), nil)

			return
		}
		c.updateSaveButton()
		gtkutil.ShowInfoDialog(c.view.Window, configRestartMessage, nil)
	}

	onRestore := func() {
		c.setEditorContent(c.model.DefaultTOML())
		c.updateSaveButton()
	}

	c.view.ConnectSignals(map[string]any{
		"on_save_config":     onSave,
		"on_restore_default": onRestore,
	})

	gtkutil.ShowModalWindow(c.view.Window)
}

func (c *ConfigEditorDialogController) setEditorContent(content string) {
	gtkutil.SetTextViewContent(c.view.TextView, content)
	gtkutil.ApplyTOMLSyntaxHighlight(c.view.TextView)
}

func (c *ConfigEditorDialogController) updateSaveButton() {
	content := gtkutil.GetTextViewContent(c.view.TextView)
	c.view.ButtonSave.SetSensitive(c.model.IsDirty(content))
}
