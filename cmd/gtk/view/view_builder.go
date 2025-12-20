//go:build gtk

package view

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

// ViewBuilder is a small embedded helper for views created from a gtk.Builder.
// It provides ConnectSignals so every view doesn't repeat the same boilerplate.
type ViewBuilder struct {
	builder *gtk.Builder
}

func NewViewBuilder(ui []byte) ViewBuilder {
	builder, err := gtk.BuilderNewFromString(string(ui))
	gtkutil.FatalErrorCheck(err)

	return ViewBuilder{builder: builder}
}

func (vb *ViewBuilder) Builder() *gtk.Builder {
	return vb.builder
}

func (vb *ViewBuilder) GetObj(name string) glib.IObject {
	obj, err := vb.builder.GetObject(name)
	gtkutil.FatalErrorCheck(err)

	return obj
}

func (vb *ViewBuilder) GetApplicationWindowObj(name string) *gtk.ApplicationWindow {
	return vb.GetObj(name).(*gtk.ApplicationWindow)
}

func (vb *ViewBuilder) GetDialogObj(name string) *gtk.Dialog {
	return vb.GetObj(name).(*gtk.Dialog)
}

func (vb *ViewBuilder) GetAboutDialogObj(name string) *gtk.AboutDialog {
	return vb.GetObj(name).(*gtk.AboutDialog)
}

func (vb *ViewBuilder) GetComboBoxTextObj(name string) *gtk.ComboBoxText {
	return vb.GetObj(name).(*gtk.ComboBoxText)
}

func (vb *ViewBuilder) GetEntryObj(name string) *gtk.Entry {
	return vb.GetObj(name).(*gtk.Entry)
}

func (vb *ViewBuilder) GetOverlayObj(name string) *gtk.Overlay {
	return vb.GetObj(name).(*gtk.Overlay)
}

func (vb *ViewBuilder) GetTreeViewObj(name string) *gtk.TreeView {
	return vb.GetObj(name).(*gtk.TreeView)
}

func (vb *ViewBuilder) GetTextViewObj(name string) *gtk.TextView {
	return vb.GetObj(name).(*gtk.TextView)
}

func (vb *ViewBuilder) GetBoxObj(name string) *gtk.Box {
	return vb.GetObj(name).(*gtk.Box)
}

func (vb *ViewBuilder) GetLabelObj(name string) *gtk.Label {
	return vb.GetObj(name).(*gtk.Label)
}

func (vb *ViewBuilder) GetToolButtonObj(name string) *gtk.ToolButton {
	return vb.GetObj(name).(*gtk.ToolButton)
}

func (vb *ViewBuilder) GetButtonObj(name string) *gtk.Button {
	return vb.GetObj(name).(*gtk.Button)
}

func (vb *ViewBuilder) GetImageObj(name string) *gtk.Image {
	return vb.GetObj(name).(*gtk.Image)
}

func (vb *ViewBuilder) GetProgressBarObj(name string) *gtk.ProgressBar {
	return vb.GetObj(name).(*gtk.ProgressBar)
}

func (vb *ViewBuilder) GetMenuItem(name string) *gtk.MenuItem {
	return vb.GetObj(name).(*gtk.MenuItem)
}

func (vb *ViewBuilder) BuildExtendedEntry(name string) *gtk.Entry {
	return gtkutil.BuildExtendedEntry(vb.builder, name)
}

func (vb *ViewBuilder) ConnectSignals(signals map[string]any) {
	if vb == nil || vb.builder == nil {
		return
	}
	vb.builder.ConnectSignals(signals)
}
