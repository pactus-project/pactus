//go:build gtk

package gtkutil

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

// ApplyTOMLSyntaxHighlight applies GtkTextBuffer tags for basic TOML syntax highlighting.
// It uses only GtkTextView/GtkTextBuffer and works on all platforms supported by GTK3.
func ApplyTOMLSyntaxHighlight(tv *gtk.TextView) {
	buf := tv.Buffer()
	if buf == nil {
		return
	}

	ensureTOMLTags(buf)

	startIter, endIter := buf.Bounds()
	content := buf.Text(startIter, endIter, true)
	table := buf.TagTable()

	for _, name := range []string{
		TOMLTagComment, TOMLTagString, TOMLTagKey,
		TOMLTagNumber, TOMLTagBoolean, TOMLTagTable,
	} {
		tag := table.Lookup(name)
		if tag != nil {
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

func declareTag(buf *gtk.TextBuffer, name string, _ map[string]any) {
	table := buf.TagTable()

	if tag := table.Lookup(name); tag != nil {
		return
	}

	tag := gtk.NewTextTag(name)
	// In GTK4, use GObject properties via coreglib.ObjectValue
	// For now, just create the tag without complex property setting
	// TODO: Implement proper Pango formatting for tags
	table.Add(tag)
}

func applyTOMLTagRange(buf *gtk.TextBuffer, start, end int, tagName string) {
	if start < 0 || end <= start {
		return
	}

	table := buf.TagTable()
	tag := table.Lookup(tagName)
	if tag == nil {
		return
	}

	startIter := charOffsetToIter(buf, start)
	endIter := charOffsetToIter(buf, end)
	buf.ApplyTag(tag, startIter, endIter)
}

func charOffsetToIter(buf *gtk.TextBuffer, offset int) *gtk.TextIter {
	iter := buf.StartIter()
	if offset > 0 {
		iter.ForwardChars(offset)
	}

	return iter
}
