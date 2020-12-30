package block

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
)

func TestRandomBlock(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.NoError(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.TxIDs = TxIDs{}
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.Version = 2
	assert.Error(t, b.SanityCheck())

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

func TestDecode(t *testing.T) {
	var b1 Block
	d, _ := hex.DecodeString("a301a90101021a5fdfb45203582000000000000000000000000000000000000000000000000000000000000000000458205678536a4d13daad8e4961de56bc278c0fc1e826c85eae89e82910aea829541205582088e2d058f12a87cb919a9f6d0c6aa2d8f85fbff02ea6b307420b6fdaa87f83df06582000000000000000000000000000000000000000000000000000000000000000000758200000000000000000000000000000000000000000000000000000000000000000085820b44cb31963c4de8065421bfdae7ab89daa9b52b4c57d7e56ea427cff2c8ab20a095425187ebec72fc662d422e2afea5b2f70494a3e5702f603a10181582088e2d058f12a87cb919a9f6d0c6aa2d8f85fbff02ea6b307420b6fdaa87f83df")
	assert.NoError(t, b1.Decode(d))
	d2, _ := b1.Encode()
	assert.Equal(t, d, d2)
	h, _ := crypto.HashFromString("c376ec4c54d891fd8118cfc5a4c4022201a0a229461846ffa26dad408a48fe70")
	assert.True(t, b1.HashesTo(h))
}

func TestBlockFingerprint(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.Contains(t, b.Fingerprint(), b.Hash().Fingerprint())
	assert.Contains(t, b.Fingerprint(), b.Header().CommittersHash().Fingerprint())
}
