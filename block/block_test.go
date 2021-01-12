package block

import (
	"encoding/hex"
	"testing"

	"github.com/zarbchain/zarb-go/util"

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

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.LastCommit = nil
	assert.Error(t, b.SanityCheck())
}

func TestMarshaling(t *testing.T) {
	b1, _ := GenerateTestBlock(nil, nil)

	bz1, err := b1.MarshalCBOR()
	assert.NoError(t, err)
	var b2 Block
	err = b2.UnmarshalCBOR(bz1)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())

	bz2, _ := b1.MarshalCBOR()
	assert.Equal(t, bz1, bz2)
}

func TestJSONMarshaling(t *testing.T) {
	b1, _ := GenerateTestBlock(nil, nil)

	bz1, err := b1.MarshalJSON()
	assert.NoError(t, err)
	var b2 Block
	err = b2.UnmarshalJSON(bz1)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())
}

func TestDecode(t *testing.T) {
	var b1 Block
	d, _ := hex.DecodeString("a301a90101021a5ffc5ee30358207f3ba277e58a59d104beb6e3eac20651f19532feb5d65db2599986e7b6749be40458204f709b4fd94ad2f0c926dd07338f2d0daf84d9ce327eff6fd7ef83a34b79c0650558202c189fad5e2338db67c5f748762d0b342cda8f189f54b1d47984257679cecedb065820be5732573591992d8c471cb8ccce8311da20f1343f23c32df5a763f08da4caaf075820ad27e7ab698dd476b85b03123bf5e9713db4c5593a66c7a110614aa85cc893e90858201098b3bab1d3d8a972cb4ea0ec9d92baa82fdd4b6464c3c36e75a9c23c0013440954bb432da41d0427c0adcad7b9a64ebfa3c5f81c3a02a50158207f3ba277e58a59d104beb6e3eac20651f19532feb5d65db2599986e7b6749be402030384000102030480055830c8349ae83c19f5d27e82e446bc6aca6852e5ef6f7d94350535f1f3c0776fd8fe25d3b1dbab3f9ba4787d2abbba9fc59903a101845820f5e2075bd72f68181ccf6e631d1410510c86e1e8778e8958c2648d72eadda038582014a8c257a6df0dd0ad5d2c2efb13604ef9bc1e739d604bbdb4c5d4a2dfac9cae582024b4ac207f9003951c491a60be1d1c78df0263515ca012314b6d432fb05901395820c982c23ab3f02c7735d2de8336d1f6699284a41164f17504fb884344be110678")
	assert.NoError(t, b1.Decode(d))
	d2, _ := b1.Encode()
	assert.Equal(t, d, d2)
	h, _ := crypto.HashFromString("aa09f4acc00bbeb9edf15a26af94861ef055a098e6d981f2375d5424998047b9")
	assert.True(t, b1.HashesTo(h))
}

func TestSanityCheck(t *testing.T) {
	tmp, _ := GenerateTestBlock(nil, nil)
	t.Run("Invalid block information, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		MakeBlock(util.Now(), tmp.TxIDs(),
			tmp.Header().LastBlockHash(),
			tmp.Header().CommittersHash(),
			tmp.Header().StateHash(),
			tmp.Header().LastReceiptsHash(), nil, tmp.Header().ProposerAddress())
	})
}

func TestBlockFingerprint(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.Contains(t, b.Fingerprint(), b.Hash().Fingerprint())
	assert.Contains(t, b.Fingerprint(), b.Header().CommittersHash().Fingerprint())
}
