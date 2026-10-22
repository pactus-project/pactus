//go:build gtk

package gtkutil

import (
	"errors"
	"testing"

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
