//go:build gtk

package view

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/pactus-project/pactus/cmd/gtk/assets"
	"github.com/pactus-project/pactus/cmd/gtk/gtkutil"
)

type ValidatorWidgetView struct {
	ViewBuilder

	Box *gtk.Box

	ColViewValidators *gtk.ColumnView
}

func NewValidatorWidgetView() *ValidatorWidgetView {
	builder := NewViewBuilder(assets.ValidatorWidgetUI)

	view := &ValidatorWidgetView{
		ViewBuilder:       builder,
		Box:               builder.GetBoxObj("id_box_validator"),
		ColViewValidators: builder.GetColumnViewObj("id_columnview_validators"),
	}

	gtkutil.ColumnViewSetDefaultProperties(view.ColViewValidators)

	return view
}
