package main

import (
	_ "embed"

	"github.com/gotk3/gotk3/gtk"
)

type widgetNode struct {
	*gtk.Box
}

//go:embed assets/ui/widget_node.ui
var uiWidgetNode []byte

func buildWidgetNode() *widgetNode {
	builder, err := gtk.BuilderNewFromString(string(uiWidgetNode))
	errorCheck(err)

	objBox, err := builder.GetObject("id_box_node")
	errorCheck(err)

	box, err := isBox(objBox)
	errorCheck(err)

	w := &widgetNode{
		Box: box,
	}

	signals := map[string]interface{}{}
	builder.ConnectSignals(signals)

	return w
}
