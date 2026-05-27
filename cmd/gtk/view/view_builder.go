//go111:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

func GetObj[T glib.Objector](b *gtk.Builder, name string) T {
	return b.GetObject(name).Cast().(T)
}

// ViewBuilder is a small embedded helper for views created from a gtk.Builder.
type ViewBuilder struct {
	builder *gtk.Builder
}

func NewViewBuilder(ui []byte) ViewBuilder {
	builder := gtk.NewBuilderFromString(string(ui))

	return ViewBuilder{builder: builder}
}

func (vb *ViewBuilder) Builder() *gtk.Builder {
	return vb.builder
}

func (vb *ViewBuilder) GetApplicationWindowObj(name string) *gtk.ApplicationWindow {
	return GetObj[*gtk.ApplicationWindow](vb.builder, name)
}

func (vb *ViewBuilder) GetWindowObj(name string) *gtk.Window {
	return GetObj[*gtk.Window](vb.builder, name)
}

func (vb *ViewBuilder) GetAboutDialogObj(name string) *gtk.AboutDialog {
	return GetObj[*gtk.AboutDialog](vb.builder, name)
}

func (vb *ViewBuilder) GetComboBoxTextObj(name string) *gtk.ComboBoxText {
	return GetObj[*gtk.ComboBoxText](vb.builder, name)
}

func (vb *ViewBuilder) GetEntryObj(name string) *gtk.Entry {
	return GetObj[*gtk.Entry](vb.builder, name)
}

func (vb *ViewBuilder) GetOverlayObj(name string) *gtk.Overlay {
	return GetObj[*gtk.Overlay](vb.builder, name)
}

func (vb *ViewBuilder) GetColumnViewObj(name string) *gtk.ColumnView {
	return GetObj[*gtk.ColumnView](vb.builder, name)
}

func (vb *ViewBuilder) GetTextViewObj(name string) *gtk.TextView {
	return GetObj[*gtk.TextView](vb.builder, name)
}

func (vb *ViewBuilder) GetBoxObj(name string) *gtk.Box {
	return GetObj[*gtk.Box](vb.builder, name)
}

func (vb *ViewBuilder) GetLabelObj(name string) *gtk.Label {
	return GetObj[*gtk.Label](vb.builder, name)
}

func (vb *ViewBuilder) GetButtonObj(name string) *gtk.Button {
	return GetObj[*gtk.Button](vb.builder, name)
}

func (vb *ViewBuilder) GetImageObj(name string) *gtk.Image {
	return GetObj[*gtk.Image](vb.builder, name)
}

func (vb *ViewBuilder) GetProgressBarObj(name string) *gtk.ProgressBar {
	return GetObj[*gtk.ProgressBar](vb.builder, name)
}

func (vb *ViewBuilder) GetPopoverMenu(name string) *gtk.PopoverMenu {
	return GetObj[*gtk.PopoverMenu](vb.builder, name)
}

func (vb *ViewBuilder) BuildExtendedEntry(name string) *gtk.Entry {
	return gtkutil.BuildExtendedEntry(vb.builder, name)
}

func (vb *ViewBuilder) ConnectSignals(signals map[string]any) {
	for key, val := range signals {
		vb.builder.Connect(key, val)
	}
}
