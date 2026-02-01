//go:build gtk

package gtkutil

import (
	"fmt"
	"strings"
)

func escapeMarkup(text string) string {
	// Minimal escaping for Pango markup contexts.
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&apos;",
	)

	return r.Replace(text)
}

// SmallGray wraps text in a small gray Pango markup span.
func SmallGray(text string) string {
	return fmt.Sprintf("<span foreground='gray' size='small'>%s</span>", escapeMarkup(text))
}

// ImportedLabel appends the imported suffix if needed.
func ImportedLabel(label string, imported bool) string {
	if !imported {
		return label
	}
	if strings.TrimSpace(label) == "" {
		return "(Imported)"
	}

	return label + " (Imported)"
}

func YesNo(v bool) string {
	if v {
		return "Yes"
	}

	return "No"
}
