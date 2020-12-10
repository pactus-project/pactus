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
	expected1, _ := crypto.HashFromString("bff193b1129bb1f15b73964d81916d71d4572f12fd413ca007e7370cf9bdfb00")
	assert.Equal(t, c.CommitersHash(), expected1)

	expected2, _ := crypto.HashFromString("576b40382ae8cae729c64168aa91f159c75392ddb851ca48ae5b1e1a79ba62c5")
	assert.Equal(t, c.Hash(), expected2)
}

func TestCommitMerkle(t *testing.T) {
	b, _ := GenerateTestBlock(nil)

	commiters := b.LastCommit().Commiters()
	data := make([]crypto.Hash, len(commiters))
	for i, c := range commiters {
		b := c.Address.RawBytes()
		data[i] = crypto.HashH(b)
	}
	merkle := simpleMerkle.NewTreeFromHashes(data)
	assert.Equal(t, merkle.Root(), b.LastCommit().CommitersHash())
}

func TestCommitSanityCheck(t *testing.T) {
	b, _ := GenerateTestBlock(nil)
	c := b.LastCommit()
	assert.NoError(t, c.SanityCheck())
	c.data.Commiters[0].Status = 0 // not signed
	// Not enough signer
	assert.Error(t, c.SanityCheck())
	c.data.Commiters[3].Status = 1 // signed
	assert.NoError(t, c.SanityCheck())
	c.data.Commiters[3].Status = 2 // invalid status
	assert.Error(t, c.SanityCheck())
}
