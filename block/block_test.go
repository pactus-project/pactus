package block

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/util"
)

func TestBlockSanityCheck(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.NoError(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.TxIDs = TxIDs{}
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.StateHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.CommitteeHash = crypto.UndefHash
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
	b.data.LastCertificate.data.Round = b.data.LastCertificate.data.Round + 1
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = crypto.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.LastCertificate = nil
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

	assert.Equal(t, b1.Hash(), b2.Hash())
	assert.Equal(t, b1.Header().Time(), b2.Header().Time())
	assert.Equal(t, b1.Header().Version(), b2.Header().Version())
	assert.Equal(t, b2.Header().Version(), 1)
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
	d, _ := hex.DecodeString("a301aa0101021a604645da0358209637ca15ffe5315c10b8e6e3e4bf25bc7c640f74868b1991fbd32db9c8d37b8b0458200791da273cca8b346babc659efc3a63c56f237151018383558401fd0ad87cf3a055820a5d857164beeb990dcfffc3bb4eb491f996b07214f1695bb2982e8289d6cf6980658205e6dd4fc470aaa2628943d569ef781f0009936d17fc45e3de8cf813cf3e44071075820669ddf68da8a273654ad63d9f1afc334a93479c5102995557aac23634f50b7750858207ae84cba480d3e28fef96c1bef43f363e329537643a16fcc6859a6a7870dd545095830ee5dbb1990c9daca3216206a631881edc274cd5ada5d105f7f1b045f538fed8be36b919a8ef5c9c90ae261f2855eecdd0a542d1cab8e470c21d31cf1ef52696002343a84710a02a50158209637ca15ffe5315c10b8e6e3e4bf25bc7c640f74868b1991fbd32db9c8d37b8b020003840001020304810008583013c067dfe7d99c0055788e9ea4a3bcefa657fa25d803c0ae0e5bb8e8ed7b1a7427830285e80bc2461c53959f4bce0b0b03a101845820c8e67e955d7a20e4cc9193f5d7a366cf6b905edf5dc7138d7309c74422025ad2582045d62a1e8c47ec7f2e8a9cf59d22db5b88601fc0d661f8cacbde0b6cb12541fd58200fdfadd817200342615dfc0bacbe0859914015825e0cdd850d999405796227ba5820c0df2bc3a82cb6330b6940077618a2d6dfddd8c90b48210278e3d699faa72863")
	assert.NoError(t, b1.Decode(d))
	d2, _ := b1.Encode()
	assert.Equal(t, d, d2)
	h, _ := crypto.HashFromString("93f65851a091dfa9cac6788e95901fa3e008d4d994e0690cb42260d2cc13f906")
	assert.True(t, b1.HashesTo(h))
}

func TestSanityCheck(t *testing.T) {
	tmp, _ := GenerateTestBlock(nil, nil)
	t.Run("Invalid block information, Certificate is missed, Should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		MakeBlock(1, util.Now(), tmp.TxIDs(),
			tmp.Header().LastBlockHash(),
			tmp.Header().CommitteeHash(),
			tmp.Header().StateHash(),
			tmp.Header().LastReceiptsHash(),
			nil,
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})
}
