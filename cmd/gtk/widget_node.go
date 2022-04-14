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

	box := getBoxObj(builder, "id_box_node")

	w := &widgetNode{
		Box: box,
	}

	signals := map[string]interface{}{}
	builder.ConnectSignals(signals)

	return w
}
