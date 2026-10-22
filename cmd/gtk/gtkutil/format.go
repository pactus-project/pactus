//go:build gtk

package gtkutil

import (
	"fmt"
	"strings"

	"github.com/pactus-project/pactus/types/amount"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
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

// LabelValueMarkup returns "label: value" as Pango markup, with the value in bold.
func LabelValueMarkup(label, value string) string {
	return fmt.Sprintf("%s: <b>%s</b>", escapeMarkup(label), escapeMarkup(value))
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

// AvailabilityScorePercent formats an availability score (0–1) as a percentage string with 2 decimal places.
func AvailabilityScorePercent(score float64) string {
	return fmt.Sprintf("%.2f%%", score*100)
}

// AddressItems returns the items of an address list as "address - amount".
// The amount comes from getAmount; when it fails, for example for an address
// that is not on the blockchain yet, the item is the address alone.
func AddressItems(infos []*pactus.AddressInfo, getAmount func(addr string) (amount.Amount, error)) []string {
	items := make([]string, 0, len(infos))
	for _, info := range infos {
		amt, err := getAmount(info.Address)
		if err != nil {
			items = append(items, info.Address)

			continue
		}
		items = append(items, fmt.Sprintf("%s - %s", info.Address, amt))
	}

	return items
}

// AddressFromItem returns the address of an address list item. It accepts both
// "address - amount" items and plain addresses typed by the user.
func AddressFromItem(item string) string {
	addr, _, _ := strings.Cut(strings.TrimSpace(item), " ")

	return addr
}
