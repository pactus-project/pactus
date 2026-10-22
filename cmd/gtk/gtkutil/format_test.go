//go:build gtk

package gtkutil

import (
	"errors"
	"testing"
	"time"

	"github.com/pactus-project/pactus/types/amount"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
)

func TestSmallGray_EscapesMarkup(t *testing.T) {
	got := SmallGray(`a<b&"c"`)
	assert.NotEmpty(t, got)
	assert.NotContains(t, got, "<b", "expected markup to be escaped")
	assert.NotContains(t, got, "&\"", "expected markup to be escaped")
	assert.NotContains(t, got, `"`, "expected markup to be escaped")
	assert.Contains(t, got, "&lt;")
	assert.Contains(t, got, "&amp;")
}

func TestAddressItems(t *testing.T) {
	infos := []*pactus.AddressInfo{
		{Address: "pc1zfunded"},
		{Address: "pc1znew"},
	}
	getAmount := func(addr string) (amount.Amount, error) {
		if addr == "pc1zfunded" {
			return amount.Amount(1_234_500_000_000), nil
		}

		return 0, errors.New("not found")
	}

	items := AddressItems(infos, getAmount)
	assert.Equal(t, []string{"pc1zfunded - 1,234.5 PAC", "pc1znew"}, items)
}

func TestAddressFromItem(t *testing.T) {
	assert.Equal(t, "pc1zfunded", AddressFromItem("pc1zfunded - 1,234.5 PAC"))
	assert.Equal(t, "pc1znew", AddressFromItem("pc1znew"))
	assert.Equal(t, "pc1ptyped", AddressFromItem("  pc1ptyped  "))
	assert.Empty(t, AddressFromItem(""))
}

func TestLabelValueMarkup(t *testing.T) {
	assert.Equal(t, "Account Balance: <b>1,234.5 PAC</b>", LabelValueMarkup("Account Balance", "1,234.5 PAC"))
	assert.Equal(t, "a&lt;b: <b>c&amp;d</b>", LabelValueMarkup("a<b", "c&d"))
}

func TestFormatTimeAgo(t *testing.T) {
	now := time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)
	ago := func(d time.Duration) int64 {
		return now.Add(-d).Unix()
	}

	tests := []struct {
		name      string
		timestamp int64
		want      string
	}{
		{"zero", 0, ""},
		{"negative", -1, ""},
		{"seconds", ago(30 * time.Second), "just now"},
		{"future", now.Add(time.Minute).Unix(), "just now"},
		{"one minute", ago(time.Minute), "1 minute ago"},
		{"five minutes", ago(5 * time.Minute), "5 minutes ago"},
		{"fifty nine minutes", ago(59 * time.Minute), "59 minutes ago"},
		{"one hour", ago(time.Hour), "1 hour ago"},
		{"five hours", ago(5 * time.Hour), "5 hours ago"},
		{"almost a day", ago(23*time.Hour + 59*time.Minute), "23 hours ago"},
		{"one day", ago(24 * time.Hour), time.Unix(ago(24*time.Hour), 0).Format("02 Jan 2006 15:04")},
		{"two days", ago(48 * time.Hour), time.Unix(ago(48*time.Hour), 0).Format("02 Jan 2006 15:04")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatTimeAgo(tt.timestamp, now))
		})
	}
}
