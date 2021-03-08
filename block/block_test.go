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
	d, _ := hex.DecodeString("a301aa0101021a603b4e7503582080a7ef2b756b6ddc11e1035671dd25362ba746036ed9fed65004f5fac988fff20458202700ec88439ec052f1c209bdf318d84837eb19532958f01396d96779a99cbb56055820cd8bdb9a0294b9e6f89432f222fa393aa73e3213763f5b717b4cd86c41387dfb065820ece5ccec89371b996a99940942f883e88226bbc291ce4c9310db47d25509c10c0758202f46746b1c3dac686b024f669602f9e9f65f46890d429bdb0baa2aed3a7ea6a6085820cf10683918a60e31824af1f101851a0605d2c84426fa8c575091302dad864d41095830852c3956ee7bb1794423e10fff0218d9e9c4e52a1d20da039ff6585036bbc827d3d4b428f340a67dec0ffc0949bbbb910a544b6c7faef3da967d1d349d97b2ba346b862d77d502a401582080a7ef2b756b6ddc11e1035671dd25362ba746036ed9fed65004f5fac988fff202020384a201000200a201010201a201020201a20103020104583066732a44f7f816d5accd3e2da7cd088652d4d4a93efaa0059786e3c7cea30608d8d820e5459bb186d56aa17cad80680f03a101845820dcb2b201186d91cb600aaad2d7452f969cc42d4cf9b57ed980176fa34759cd0d5820dc679a3e457039504a511fa99abafcc493e185382c05b7f2c317dd4ac38440935820816142193a73d2196c6bc86fc8e5cfd97fbdf0c113293c21e3a193848f22f8365820074fad84a213c6b720128131bb585cf96612f7fbd08d07db00b10ca4e3682553")
	assert.NoError(t, b1.Decode(d))
	d2, _ := b1.Encode()
	assert.Equal(t, d, d2)
	h, _ := crypto.HashFromString("8568b336462912580a2023423b1171c48f2ff8927bde2a52fd9993d42f1227b6")
	assert.True(t, b1.HashesTo(h))
}

func TestSanityCheck(t *testing.T) {
	tmp, _ := GenerateTestBlock(nil, nil)
	t.Run("Invalid block information, Commit is missed, Should panic", func(t *testing.T) {
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
