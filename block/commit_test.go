package block

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"

	"github.com/zarbchain/zarb-go/crypto"
	simpleMerkle "github.com/zarbchain/zarb-go/libs/merkle"
)

func TestNilCommitHash(t *testing.T) {
	var c Commit
	assert.Equal(t, c.Hash(), crypto.UndefHash)
}

func TestCommitMarshaling(t *testing.T) {
	d, _ := hex.DecodeString("a301020258304e02ec725dd2b2ae23fc6a90bb391e3a61c1488eb0bffdffb75088da8dc7e9b267f396898dbe3e7971dd2af9395200810384a20154a7dd14a976e2894602a0c04081b66258bd930faa0201a20154f2150d173fdf8e5435712e3731237e4751675ef30201a20154c1ecaa8747f46553556d484d1345f7e152eddee20201a20154817a3ea1b55ebb68c29d45592d41da6bedb7f3350200")
	c := new(Commit)
	err := cbor.Unmarshal(d, c)
	assert.NoError(t, err)
	d2, err := cbor.Marshal(c)
	assert.NoError(t, err)
	assert.Equal(t, d, d2)
	expected1, _ := crypto.HashFromString("52568b5383dd52e75bd8a956c49d80986b49802c15fc51f87a6f0650420a1bac")
	assert.Equal(t, c.CommittersHash(), expected1)

	expected2, _ := crypto.HashFromString("e57cf7135136274b6db0165408dc2ff3ccfd22f2a13b85292e0984e79a310e17")
	assert.Equal(t, c.Hash(), expected2)
}

func TestCommitMerkle(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)

	committers := b.LastCommit().Committers()
	data := make([]crypto.Hash, len(committers))
	for i, c := range committers {
		b := c.Address.RawBytes()
		data[i] = crypto.HashH(b)
	}
	merkle := simpleMerkle.NewTreeFromHashes(data)
	assert.Equal(t, merkle.Root(), b.LastCommit().CommittersHash())
}

func TestCommitSanityCheck(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	c := b.LastCommit()
	assert.NoError(t, c.SanityCheck())
	c.data.Committers[0].Status = 0 // not signed
	// Not enough signer
	assert.Error(t, c.SanityCheck())
	c.data.Committers[3].Status = 1 // signed
	assert.NoError(t, c.SanityCheck())
	c.data.Committers[3].Status = 2 // invalid status
	assert.Error(t, c.SanityCheck())
}
