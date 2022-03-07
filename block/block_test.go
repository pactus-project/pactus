package block

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto"
	"github.com/zarbchain/zarb-go/crypto/hash"
	simplemerkle "github.com/zarbchain/zarb-go/libs/merkle"
	"github.com/zarbchain/zarb-go/sortition"
	"github.com/zarbchain/zarb-go/util"
)

func TestBlockSanityCheck(t *testing.T) {
	b, _ := GenerateTestBlock(nil, nil)
	assert.NoError(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.TxIDs = TxIDs{}
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.StateHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.PrevBlockHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.PrevCertificateHash = hash.UndefHash
	b.data.Header.data.PrevBlockHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.SortitionSeed = sortition.UndefVerifiableSeed
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.ProposerAddress = crypto.TreasuryAddress
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.PrevCertificateHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.PrevCertificate = GenerateTestCertificate(hash.UndefHash)
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.PrevCertificate = GenerateTestCertificate(hash.GenerateTestHash())
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.Header.data.TxIDsHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.PrevCertificate = nil
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

	bz, err := b1.MarshalJSON()
	assert.NoError(t, err)
	assert.NotNil(t, bz)
}

func TestDecode(t *testing.T) {
	var b1 Block
	/*
		{
			1: {
				1: 1,
				2: 1641381449,
				3: h'CC023E43C4BB7B111A7029EA224887537B4C70AA201054721D463091A4266698',
				4: h'A3058DC8E31E3D908522287EF34B0FB8E7AE6E70CB376FCA6BD113F74843CE3C',
				5: h'07CF8EFBBA30D33A433937C9A88816187647B844E674CBEE7231C8E3F7188A0D',
				6: h'DB5815DC702EA6ED3DD5E141DD90DAAC82C9B7EC7394E7B6F576CF65E843857C',
				7: h'A54825BC7F77EAC94E34203FE43A6D0222BEEA34224DCF6C86E9E55E479709FA3E53267860BA0FCC14DD911E8C0208EE',
				8: h'01B05EC5B163832FFE1B588470D49AFF2A86BE1F8E'
			},
			2: {
				1: h'CC023E43C4BB7B111A7029EA224887537B4C70AA201054721D463091A4266698',
				2: 3,
				3: [10, 18, 12, 16],
				4: [18],
				5: h'B7A82C851FD0B35A66683A1D4B0D7AAC4F5B97D4C72E882E2B10C8B25ADA75C1AB1C8F97CA2615F29B89BC8F25DDB6B9'
			},
			3: {
				1: [
					h'051EE8781A7352E5BEFC7183123FF1A6F436F6E3C5FDD4A6AACCDBAA4CEC551A',
					h'EB121430ECDB1DAD019475B7C76D03A864F1B279CE9AE12EBDBF93EB8EA76E5C',
					h'2A6AE8B2C89D9F88B8AB7E7FA6DB4058383D6DE1AA7C8F8DCA9C98FFF8FE4F89',
					h'4C59BB0C74EB67C974AA8F69E856E2B6B02618E1E92B369EFD97F761E05CDB70'
				]
			}
		}
	*/
	d, _ := hex.DecodeString("a301a80101021a61d57e49035820cc023e43c4bb7b111a7029ea224887537b4c70aa201054721d463091a4266698045820a3058dc8e31e3d908522287ef34b0fb8e7ae6e70cb376fca6bd113f74843ce3c05582007cf8efbba30d33a433937c9a88816187647b844e674cbee7231c8e3f7188a0d065820db5815dc702ea6ed3dd5e141dd90daac82c9b7ec7394e7b6f576cf65e843857c075830a54825bc7f77eac94e34203fe43a6d0222beea34224dcf6c86e9e55e479709fa3e53267860ba0fcc14dd911e8c0208ee085501b05ec5b163832ffe1b588470d49aff2a86be1f8e02a5015820cc023e43c4bb7b111a7029ea224887537b4c70aa201054721d463091a4266698020303840a120c10048112055830b7a82c851fd0b35a66683a1d4b0d7aac4f5b97d4c72e882e2b10c8b25ada75c1ab1c8f97ca2615f29b89bc8f25ddb6b903a101845820051ee8781a7352e5befc7183123ff1a6f436f6e3c5fdd4a6aaccdbaa4cec551a5820eb121430ecdb1dad019475b7c76d03a864f1b279ce9ae12ebdbf93eb8ea76e5c58202a6ae8b2c89d9f88b8ab7e7fa6db4058383d6de1aa7c8f8dca9c98fff8fe4f8958204c59bb0c74eb67c974aa8f69e856e2b6b02618e1e92b369efd97f761e05cdb70")
	assert.NoError(t, b1.Decode(d))
	d2, _ := b1.Encode()
	assert.Equal(t, d, d2)

	// block header: a80101021a61d57e49035820cc023e43c4bb7b111a7029ea224887537b4c70aa201054721d463091a4266698045820a3058dc8e31e3d908522287ef34b0fb8e7ae6e70cb376fca6bd113f74843ce3c05582007cf8efbba30d33a433937c9a88816187647b844e674cbee7231c8e3f7188a0d065820db5815dc702ea6ed3dd5e141dd90daac82c9b7ec7394e7b6f576cf65e843857c075830a54825bc7f77eac94e34203fe43a6d0222beea34224dcf6c86e9e55e479709fa3e53267860ba0fcc14dd911e8c0208ee085501b05ec5b163832ffe1b588470d49aff2a86be1f8e
	expected1 := hash.CalcHash(d[2:225])
	assert.True(t, b1.HashesTo(expected1))

	expected2, _ := hash.FromString("f82c75379f1cf017299f722213236fc7b4711366f23dbf733025880baa1d05c5")
	assert.Equal(t, b1.Hash(), expected2)
	assert.Equal(t, b1.Stamp(), hash.Stamp{0xf8, 0x2c, 0x75, 0x37})

	// hash TxIDs
	merkleTree := simplemerkle.NewTreeFromHashes([]hash.Hash{b1.TxIDs().IDs()[0], b1.TxIDs().IDs()[1], b1.TxIDs().IDs()[2], b1.TxIDs().IDs()[3]})
	expected3 := merkleTree.Root()
	assert.Equal(t, b1.Header().TxIDsHash(), expected3)
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
			tmp.Header().PrevBlockHash(),
			tmp.Header().StateHash(),
			nil,
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})
}
