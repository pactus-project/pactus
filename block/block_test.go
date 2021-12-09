package block

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/crypto/hash"
	simplemerkle "github.com/zarbchain/zarb-go/libs/merkle"
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
	assert.Error(t, b.SanityCheck())

	b, _ = GenerateTestBlock(nil, nil)
	b.data.PrevCertificate = GenerateTestCertificate(hash.UndefHash)
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
				2: 1618459134,
				3: h'D846EF49A6C72390645F12970987865A795A55FA19C92DBB9CBE24D6503ECA9F', // last block hash
				4: h'8AE0A4883808290510BB77678BB24A2527D22D7DCF2D5D605EA57595260BFDF0',
				5: h'8E442E0F18A7797D7C289EAD53B7C02D9F77147003BEBBF7B0572A72FB004BBB',
				6: h'85C4963C28750EEF54BA1B14DD03FC85DBE482A280D06E0EEFB427FCB15B616C',
				7: h'DB66DDCE5CD16EC9710294769C7386977E48EEF2BC38C5A93B49EA06AC9FA8FC502976397ABC00C5DF21D2D1C757D80D',
				8: h'3BFC7DF5C9915C56E399FBE47BE7D25AEFF238B3'
			},
			2: {
				1: h'D846EF49A6C72390645F12970987865A795A55FA19C92DBB9CBE24D6503ECA9F',
				2: 6,
				3: [10, 18, 2, 6],
				4: [10],
				5: h'85C368E9E6DF4EA1B16E29AEBBF74A3DA45A033683E753C93130336E035C2181BF469DAB5E0448064FB64F6282B28296'
			},
			3: {
				1: [
					h'05D5455C116D98A90C452365E7E9A4CD03847BF7BA0EABAB4CD8ACAA417A4712',
					h'F1C6FD7464BB6D57C3B4B56995BFBB864B64EC82C8E77B454449637B5BE98FF5',
					h'3DE71E737F5AA31D148DB8A9D29BE3662584276EDD3530C92643D220997FCF40',
					h'B65B38AFD18345DB104A11AC06AC6FA692ABD808F7F7615F69182AD636B1EA08'
				]
			}
		}
	*/
	d, _ := hex.DecodeString("a301a80101021a6077b9fe035820d846ef49a6c72390645f12970987865a795a55fa19c92dbb9cbe24d6503eca9f0458208ae0a4883808290510bb77678bb24a2527d22d7dcf2d5d605ea57595260bfdf00558208e442e0f18a7797d7c289ead53b7c02d9f77147003bebbf7b0572a72fb004bbb06582085c4963c28750eef54ba1b14dd03fc85dbe482a280d06e0eefb427fcb15b616c075830db66ddce5cd16ec9710294769c7386977e48eef2bc38c5a93b49ea06ac9fa8fc502976397abc00c5df21d2d1c757d80d08543bfc7df5c9915c56e399fbe47be7d25aeff238b302a5015820d846ef49a6c72390645f12970987865a795a55fa19c92dbb9cbe24d6503eca9f020603840a12020604810a05583085c368e9e6df4ea1b16e29aebbf74a3da45a033683e753c93130336e035c2181bf469dab5e0448064fb64f6282b2829603a10184582005d5455c116d98a90c452365e7e9a4cd03847bf7ba0eabab4cd8acaa417a47125820f1c6fd7464bb6d57c3b4b56995bfbb864b64ec82c8e77b454449637b5be98ff558203de71e737f5aa31d148db8a9d29be3662584276edd3530c92643d220997fcf405820b65b38afd18345db104a11ac06ac6fa692abd808f7f7615f69182ad636b1ea08")
	assert.NoError(t, b1.Decode(d))
	d2, _ := b1.Encode()
	assert.Equal(t, d, d2)

	// a80101021a6077b9fe035820d846ef49a6c72390645f12970987865a795a55fa19c92dbb9cbe24d6503eca9f0458208ae0a4883808290510bb77678bb24a2527d22d7dcf2d5d605ea57595260bfdf00558208e442e0f18a7797d7c289ead53b7c02d9f77147003bebbf7b0572a72fb004bbb06582085c4963c28750eef54ba1b14dd03fc85dbe482a280d06e0eefb427fcb15b616c075830db66ddce5cd16ec9710294769c7386977e48eef2bc38c5a93b49ea06ac9fa8fc502976397abc00c5df21d2d1c757d80d08543bfc7df5c9915c56e399fbe47be7d25aeff238b3
	expected1 := hash.CalcHash(d[2:224])
	assert.True(t, b1.HashesTo(expected1))

	expected2, _ := hash.FromString("e161e98f19001e1f3fbcb2ae72fa8ff4a57bf8337ebebb65b0632f87763368ba")
	assert.Equal(t, b1.Hash(), expected2)

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
