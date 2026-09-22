package account

import (
	"bytes"
	"testing"

	"github.com/pactus-project/pactus/types/amount"
	"github.com/stretchr/testify/require"
)

func TestCloneDoesNotShareAnchorStorage(t *testing.T) {
	acc := NewAccount(1)
	acc.AddToBalance(2)
	require.NoError(t, acc.SetAnchor(AnchorData{
		RootHash:      bytes.Repeat([]byte{0xAB}, 32),
		ManifestURI:   "shared",
		LockedDeposit: 1,
	}))

	cloned := acc.Clone()
	cloned.data.Anchor.RootHash[0] = 0
	cloned.data.Anchor.ManifestURI = "changed"
	cloned.data.Anchor.LockedDeposit = amount.Amount(9)

	require.Equal(t, byte(0xAB), acc.data.Anchor.RootHash[0])
	require.Equal(t, "shared", acc.data.Anchor.ManifestURI)
	require.Equal(t, amount.Amount(1), acc.data.Anchor.LockedDeposit)
	require.Nil(t, NewAccount(1).Clone().data.Anchor)
}

func TestDecodedAccountsDoNotShareRootHash(t *testing.T) {
	acc := NewAccount(1)
	acc.AddToBalance(2)
	require.NoError(t, acc.SetAnchor(AnchorData{
		RootHash:      bytes.Repeat([]byte{0xCD}, 32),
		ManifestURI:   "once",
		LockedDeposit: 1,
	}))
	encoded, err := acc.Bytes()
	require.NoError(t, err)

	first, err := FromBytes(encoded)
	require.NoError(t, err)
	second, err := FromBytes(encoded)
	require.NoError(t, err)

	first.data.Anchor.RootHash[0] ^= 0xFF
	require.Equal(t, byte(0xCD), second.data.Anchor.RootHash[0])
	require.Equal(t, byte(0xCD), acc.data.Anchor.RootHash[0])

	encoded[14] ^= 0xFF
	require.Equal(t, byte(0xCD), second.RootHash()[0])
}
