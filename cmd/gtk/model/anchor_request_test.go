//go:build gtk

package model

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util/bech32m"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The vectors below come from the PIP-57 reference test vectors
// (pactus-project/PIPs, assets/pip-57/testvectors.py).

const (
	pipRoot   = "2727edf1004e656e34a07948dfccde9ae583b5f240c88f9f5dc1ae36c5ae9261"
	pipCID    = "bafk2bzaceatsp3prabhgk3ruub4urx6m32nola5v6jamrd47lxa24nwfv2jgc"
	pipURI    = "ipfs%3A%2F%2F" + pipCID
	pipSet    = "pactus:anchor?network=mainnet&action=set&type=1&root="
	pipTestEd = "tpc1rw03cn7xkd7jjl5ajjkwj67wwhsf7e252fu50nz"
)

var pipHash = strings.Repeat("ab", 32)

func pipAddress(t *testing.T, typ crypto.AddressType) string {
	t.Helper()

	addr, err := bech32m.EncodeFromBase256WithType("pc", byte(typ), make([]byte, 20))
	require.NoError(t, err)

	return addr
}

func TestParseAnchorRequestValid(t *testing.T) {
	mainBLS := pipAddress(t, crypto.AddressTypeBLSAccount)
	mainEd := pipAddress(t, crypto.AddressTypeEd25519Account)
	mainSecp := pipAddress(t, crypto.AddressTypeSecp256k1Account)

	tests := []struct {
		link    string
		network string
		check   func(t *testing.T, req *AnchorRequest)
	}{
		{
			pipSet + pipRoot + "&uri=" + pipURI + "&memo=pactus.page", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) {
				t.Helper()
				assert.Equal(t, AnchorActionSet, req.Action)
				assert.Equal(t, uint8(1), req.AnchorType)
				assert.Equal(t, pipRoot, hex.EncodeToString(req.RootHash))
				assert.Equal(t, "ipfs://"+pipCID, req.ManifestURI)
				assert.Equal(t, MinAnchorDeposit, req.Locked)
				assert.Equal(t, "pactus.page", req.Memo)
			},
		},
		{
			pipSet + pipRoot + "&uri=" + pipURI + "&locked=100&memo=pactus.page", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) {
				t.Helper()
				assert.Equal(t, amount.Amount(100*amount.NanoPACPerPAC), req.Locked)
			},
		},
		{
			"pactus:anchor?network=mainnet&action=delete&memo=pactus.page", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.True(t, req.IsDelete()) },
		},
		{
			"pactus:anchor?network=mainnet&action=delete&utm_source=x", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.True(t, req.IsDelete()) },
		},
		{
			pipSet + strings.ToUpper(pipRoot) + "&locked=2.5", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) {
				t.Helper()
				assert.Equal(t, pipRoot, hex.EncodeToString(req.RootHash))
				assert.Equal(t, amount.Amount(2_500_000_000), req.Locked)
			},
		},
		{
			pipSet + pipHash + "&locked=0.000000001", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, amount.Amount(1), req.Locked) },
		},
		{
			pipSet + pipHash, NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Empty(t, req.ManifestURI) },
		},
		{
			pipSet + pipHash + "&uri=" + strings.Repeat("x", 128), NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) {
				t.Helper()
				assert.Equal(t, strings.Repeat("x", 128), req.ManifestURI)
			},
		},
		{
			pipSet + strings.Repeat("cd", 64), NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Len(t, req.RootHash, 64) },
		},
		{
			"pactus:anchor?network=mainnet&action=set&type=200&root=" + pipHash, NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, uint8(200), req.AnchorType) },
		},
		{
			pipSet + pipHash + "&memo=a+b", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, "a+b", req.Memo) },
		},
		{
			pipSet + pipHash + "&memo=10%20%E2%82%AC", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, "10 €", req.Memo) },
		},
		{
			pipSet + pipHash + "&sender=" + mainBLS, NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, mainBLS, req.Sender) },
		},
		{
			pipSet + pipHash + "&sender=" + mainEd, NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, mainEd, req.Sender) },
		},
		{
			pipSet + pipHash + "&sender=" + mainSecp, NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, mainSecp, req.Sender) },
		},
		{
			"pactus:anchor?network=testnet&action=delete&sender=" + pipTestEd, NetworkTestnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, pipTestEd, req.Sender) },
		},
		{
			pipSet + pipHash + "&label=ignored", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, uint8(1), req.AnchorType) },
		},
		{
			pipSet + pipHash + "&memo=caf%C3%A9%20%F0%9F%8C%BB", NetworkMainnet,
			func(t *testing.T, req *AnchorRequest) { t.Helper(); assert.Equal(t, "café 🌻", req.Memo) },
		},
	}

	for _, tt := range tests {
		req, err := ParseAnchorRequest(tt.link, tt.network)
		require.NoError(t, err, tt.link)
		tt.check(t, req)
	}
}

func TestParseAnchorRequestRejected(t *testing.T) {
	mainEd := pipAddress(t, crypto.AddressTypeEd25519Account)
	mainValidator := pipAddress(t, crypto.AddressTypeValidator)
	mainTreasury := pipAddress(t, crypto.AddressTypeTreasury)
	lastChar := "q"
	if strings.HasSuffix(mainEd, "q") {
		lastChar = "p"
	}
	badChecksum := mainEd[:len(mainEd)-1] + lastChar

	links := []string{
		"pactus:anchor?network=testnet&action=set&type=1&root=" + pipHash,
		"pactus:anchor?action=set&type=1&root=" + pipHash,
		"pactus:anchor?network=betanet&action=set&type=1&root=" + pipHash,
		"pactus:anchor?network=mainnet&action=update&type=1&root=" + pipHash,
		"pactus:anchor?network=mainnet&action=set&root=" + pipHash,
		"pactus:anchor?network=mainnet&action=set&type=1",
		"pactus:anchor?network=mainnet&action=set&type=256&root=" + pipHash,
		"pactus:anchor?network=mainnet&action=set&type=01&root=" + pipHash,
		pipSet + strings.Repeat("a", 63),
		pipSet + strings.Repeat("a", 62),
		pipSet + strings.Repeat("a", 130),
		pipSet + strings.Repeat("g", 64),
		pipSet + pipHash + "&uri=" + strings.Repeat("x", 129),
		pipSet + pipHash + "&memo=" + strings.Repeat("m", 65),
		pipSet + pipHash + "&memo=%FF",
		"pactus:anchor?network=mainnet&action=set&type=1&type=2&root=" + pipHash,
		"pactus:anchor?network=mainnet&action=delete&ACTION=set",
		"pactus:anchor?network=mainnet&action=delete&%61ction=set",
		"pactus:anchor?network=mainnet&action=delete&=x",
		"pactus:anchor?network=mainnet&action=delete&type=1&root=" + pipHash,
		"pactus:anchor?network=mainnet&action=delete&type=x",
		"pactus:anchor?network=mainnet&action=delete&uri=",
		"pactus:anchor?network=mainnet&action=delete&locked=100",
		pipSet + pipHash + "&memo=a%0Ab",
		pipSet + pipHash + "&memo=a%7Fb",
		pipSet + pipHash + "&memo=a%E2%80%8Bb",
		pipSet + pipHash + "&memo=a%E2%80%A8b",
		pipSet + pipHash + "&uri=ipfs%3A%2F%2F%E2%80%AEfdp.x",
		"pactus:anchor?network=mainnet&action=delete&memo=%09",
		pipSet + pipHash + "&locked=-1",
		pipSet + pipHash + "&locked=+1",
		pipSet + pipHash + "&locked=1e2",
		pipSet + pipHash + "&locked=007",
		pipSet + pipHash + "&locked=1.",
		pipSet + pipHash + "&locked=1,5",
		pipSet + pipHash + "&locked=0.0000000001",
		pipSet + pipHash + "&locked=42000001",
		pipSet + pipHash + "&req-expiry=1",
		pipSet + pipHash + "&sender=" + mainValidator,
		pipSet + pipHash + "&sender=" + mainTreasury,
		pipSet + pipHash + "&sender=" + pipTestEd,
		pipSet + pipHash + "&sender=" + badChecksum,
		"pactus:" + mainEd + "?amount=1",
	}

	for _, link := range links {
		_, err := ParseAnchorRequest(link, NetworkMainnet)
		assert.ErrorIs(t, err, ErrInvalidAnchorLink, link)
	}
}

func TestAnchorDepositToAdd(t *testing.T) {
	tests := []struct {
		lockedNow, target, want amount.Amount
	}{
		{0, pac(1), pac(1)},
		{0, pac(0.5), pac(1)},
		{0, pac(100), pac(100)},
		{pac(1), pac(1), 0},
		{pac(1), pac(100), pac(99)},
		{pac(100), pac(1), 0},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, AnchorDepositToAdd(tt.target, tt.lockedNow),
			"locked %s, target %s", tt.lockedNow, tt.target)
	}
}

func TestCanDeleteAnchor(t *testing.T) {
	assert.False(t, CanDeleteAnchor(0, pac(0.01)))
	assert.True(t, CanDeleteAnchor(pac(100), pac(0.01)))
	assert.True(t, CanDeleteAnchor(pac(0.01), pac(0.01)))
	assert.False(t, CanDeleteAnchor(pac(0.005), pac(0.01)))
}

func TestAnchorTypeLabel(t *testing.T) {
	assert.Equal(t, "1 (Service manifest)", AnchorTypeLabel(1))
	assert.Equal(t, "255 (Private or encrypted)", AnchorTypeLabel(255))
	assert.Equal(t, "7 (Reserved)", AnchorTypeLabel(7))
}
