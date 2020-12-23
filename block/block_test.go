package block

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestRandomBlock(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.NoError(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.StateHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.CommittersHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.LastReceiptsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.LastBlockHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.LastCommitHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.LastCommit.data.Round = b.data.LastCommit.data.Round + 1
	assert.Error(t, b.SanityCheck())
}

func TestMarshaling(t *testing.T) {
	b1, _ := GenerateTestBlock(nil, nil)

	bz1, err := b1.MarshalCBOR()
	fmt.Printf("%x", bz1)
	assert.NoError(t, err)
	var b2 Block
	err = b2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())

	bz2, _ := b1.MarshalCBOR()
	assert.Equal(t, bz1, bz2)
}

func TestBlockFingerprint(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.Contains(t, b.Fingerprint(), b.Hash().Fingerprint())
	assert.Contains(t, b.Fingerprint(), b.Header().CommittersHash().Fingerprint())
}
