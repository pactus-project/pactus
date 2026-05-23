//go:build gtk

package gtkutil

import (
	"github.com/gotk3/gotk3/gtk"
)

// ApplyTOMLSyntaxHighlight applies GtkTextBuffer tags for basic TOML syntax highlighting.
// It uses only GtkTextView/GtkTextBuffer and works on all platforms supported by GTK3.
func ApplyTOMLSyntaxHighlight(tv *gtk.TextView) {
	buf, err := tv.GetBuffer()
	if err != nil {
		return
	}

	ensureTOMLTags(buf)

	startIter, endIter := buf.GetBounds()
	content, err := buf.GetText(startIter, endIter, true)
	if err != nil {
		return
	}

	table, err := buf.GetTagTable()
	if err != nil {
		return
	}

	for _, name := range []string{
		TOMLTagComment, TOMLTagString, TOMLTagKey,
		TOMLTagNumber, TOMLTagBoolean, TOMLTagTable,
	} {
		tag, lookupErr := table.Lookup(name)
		if lookupErr == nil && tag != nil {
			buf.RemoveTag(tag, startIter, endIter)
		}
	}

	for _, hr := range ComputeTOMLHighlightRanges(content) {
		applyTOMLTagRange(buf, hr.Start, hr.End, hr.Tag)
	}
}

func ensureTOMLTags(buf *gtk.TextBuffer) {
	declareTag(buf, TOMLTagComment, map[string]any{"foreground": "#6A9955"})
	declareTag(buf, TOMLTagString, map[string]any{"foreground": "#CE9178"})
	declareTag(buf, TOMLTagKey, map[string]any{"foreground": "#9CDCFE"})
	declareTag(buf, TOMLTagNumber, map[string]any{"foreground": "#B5CEA8"})
	declareTag(buf, TOMLTagBoolean, map[string]any{"foreground": "#569CD6"})
	declareTag(buf, TOMLTagTable, map[string]any{"foreground": "#DCDCAA", "weight": 700})
}

func declareTag(buf *gtk.TextBuffer, name string, props map[string]any) {
	table, err := buf.GetTagTable()
	if err != nil {
		return
	}

	if _, lookupErr := table.Lookup(name); lookupErr == nil {
		return
	}

	propsIface := make(map[string]any, len(props))
	for k, v := range props {
		propsIface[k] = v
	}

	_ = buf.CreateTag(name, propsIface)
}

func applyTOMLTagRange(buf *gtk.TextBuffer, start, end int, tagName string) {
	if start < 0 || end <= start {
		return
	}

	table, err := buf.GetTagTable()
	if err != nil {
		return
	}

	tag, lookupErr := table.Lookup(tagName)
	if lookupErr != nil || tag == nil {
		return
	}

	startIter := charOffsetToIter(buf, start)
	endIter := charOffsetToIter(buf, end)
	buf.ApplyTag(tag, startIter, endIter)
}

func charOffsetToIter(buf *gtk.TextBuffer, offset int) *gtk.TextIter {
	iter := buf.GetStartIter()

	if offset > 0 {
		iter.ForwardChars(offset)
	}

	return iter
}
