//go:build gtk

package controller

import (
	"log"

	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
	"github.com/pactus-project/pactus/cmd/gtk/model"
	"github.com/pactus-project/pactus/cmd/gtk/view"
	"libdb.so/gotk4-sourceview/pkg/gtksource/v5"
)

const configRestartMessage = "Configuration saved successfully.\n\n" +
	"You must restart the node for the changes to take effect."

type ConfigEditorDialogController struct {
	view   *view.ConfigEditorDialogView
	model  *model.ConfigModel
	buffer *gtksource.Buffer
}

func NewConfigEditorDialogController(
	view *view.ConfigEditorDialogView,
	configModel *model.ConfigModel,
) *ConfigEditorDialogController {
	buffer := gtksource.NewBuffer(nil)

	view.SourceView.SetBuffer(&buffer.TextBuffer)

	langManager := gtksource.LanguageManagerGetDefault()

	tomlLang := langManager.Language("toml")
	if tomlLang == nil {
		log.Println("Warning: TOML syntax highlighting not found.")
	} else {
		buffer.SetLanguage(tomlLang)
		buffer.SetHighlightSyntax(true)
	}

	manager := gtksource.StyleSchemeManagerGetDefault()
	scheme := manager.Scheme("Adwaita-dark")
	buffer.SetStyleScheme(scheme)

	return &ConfigEditorDialogController{
		view:   view,
		model:  configModel,
		buffer: buffer,
	}
}

func (c *ConfigEditorDialogController) Show() {
	c.setEditorContent(c.model.SavedContent())
	c.updateSaveButton()

	c.buffer.Connect("changed", func() {
		c.updateSaveButton()
	})

	onSave := func() {
		if err := c.model.Save(c.content()); err != nil {
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

	gtkutil.ConnectButtonSignal(c.view.ButtonRestore, onRestore)
	gtkutil.ConnectButtonSignal(c.view.ButtonSave, onSave)

	gtkutil.ShowModalWindow(c.view.Window)
}

func (c *ConfigEditorDialogController) setEditorContent(content string) {
	c.buffer.SetText(content)
}

func (c *ConfigEditorDialogController) updateSaveButton() {
	c.view.ButtonSave.SetSensitive(c.model.IsDirty(c.content()))
}

func (c *ConfigEditorDialogController) content() string {
	start := c.buffer.StartIter()
	end := c.buffer.EndIter()

	return c.buffer.Text(start, end, true) // 'true' includes hidden characters
}
