package block

import (
	"encoding/hex"
	"testing"

	"github.com/fxamacker/cbor/v2"
	"github.com/stretchr/testify/assert"
	"github.com/zarbchain/zarb-go/types/crypto/bls"
	"github.com/zarbchain/zarb-go/types/crypto/hash"
	"github.com/zarbchain/zarb-go/util"
	simplemerkle "github.com/zarbchain/zarb-go/util/merkle"
)

func TestSanityCheck(t *testing.T) {
	b := GenerateTestBlock(nil, nil)
	assert.NoError(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Txs = Txs{}
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.StateRoot = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.PrevBlockHash = hash.UndefHash
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.Header.data.ProposerAddress[0] = 0x2
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.PrevCert = nil
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	b.data.PrevCert.data.Round = -1
	assert.Error(t, b.SanityCheck())

	b = GenerateTestBlock(nil, nil)
	invalidSigner := bls.GenerateTestSigner()
	invalidSigner.SignMsg(b.data.Txs[0])
	assert.Error(t, b.SanityCheck())
}

func TestCBORMarshaling(t *testing.T) {
	b1 := GenerateTestBlock(nil, nil)

	bz1, err := cbor.Marshal(b1)
	assert.NoError(t, err)
	var b2 Block
	err = cbor.Unmarshal(bz1, &b2)
	assert.NoError(t, err)
	assert.NoError(t, b2.SanityCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())

	assert.Equal(t, b1.Hash(), b2.Hash())
	assert.Equal(t, b1.Header().Time(), b2.Header().Time())
	assert.Equal(t, b1.Header().Version(), b2.Header().Version())

	err = cbor.Unmarshal([]byte{1}, &b2)
	assert.Error(t, err)
}

func TestEncodingBlock(t *testing.T) {
	blk := GenerateTestBlock(nil, nil)
	len := blk.SerializeSize()

	for i := 0; i < len; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, blk.Encode(w), "encode test %v failed", i)
	}
	w := util.NewFixedWriter(len)
	assert.NoError(t, blk.Encode(w))

	for i := 0; i < len; i++ {
		blk2 := new(Block)
		r := util.NewFixedReader(i, w.Bytes())
		assert.Error(t, blk2.Decode(r), "decode test %v failed", i)
	}

	blk2 := new(Block)
	r := util.NewFixedReader(len, w.Bytes())
	assert.NoError(t, blk2.Decode(r))
	assert.Equal(t, blk.Hash(), blk2.Hash())
	assert.Equal(t, blk.Header(), blk2.Header())
}

func TestTxFromBytes(t *testing.T) {
	blk := GenerateTestBlock(nil, nil)
	bs, _ := blk.Bytes()
	_, err := FromBytes(bs)
	assert.NoError(t, err)
	assert.Equal(t, blk.memorizedData, bs)
	_, err = blk.Bytes()
	assert.NoError(t, err)

	_, err = FromBytes([]byte{1})
	assert.Error(t, err)
}

func TestBlockHash(t *testing.T) {
	d, _ := hex.DecodeString("011a873d62b69e39b4e06567b6ad3a58f61df4c3c05920a29043277af01264c9e1e7693068bbf7b5e010ca98da562965a1a3411a48fee70bd0dbbe11d9867fa9e13b3e005e99bbd54999c7cd6bb176b160962080ee130c455c88507bd51a878a0b85c656cfc1a542cbbe0105708389ca68269bda290119cba9960c6ad28aaaa140377f652bdea0551e3b0104d6041607b207011685ee6c00e9554451b2d46665ed5ca6a9a8bb1843f485a3323e05757a79a94518185dc1ce48c87634672928f0b90815f10401742b6b20dd04e981f7f1efb0140101d4dc361132b551ef27c514e44558c7da36e08794010fbd61774e607aa691785d3567f4b111fd8b283ef693d3e9de91cd1f0c746573742073656e642d74788657d98718fd5794eb37115b1d9418a63358c3a7bcf1e17146fdb71f8ce0525a657a78ee41886ddaf4887de6f0d0478685f2fbe28bb1770344e50a7b3edb34631741c403f7d19dff02a728e585fcb1284d23ad4c557329abc1873afc889e1d681423a930fdaf52a8ea591d8c13a88ec3126a3a2463b147f22878588789e58b5c07f6d5c5afbded864c4066659e30176401f10c077fcc04f5ef819fc9d6080101d3e45d249a39d806a1faec2fd85820db340b98e30168fc72a1a961933e694439b2e3c8751d27de5ad3b9c3dc91b9c9b59b010c746573742073656e642d7478b53d79e156e9417e010fa21f2b2a96bee6be46fcd233295d2f697cdb9e782b6112ac01c80d0d9d64c2320664c77fa2a68d82fa4fcac04a3b565267685e90db1b01420285d2f8295683c138c092c209479983ba1591370778846681b7b558e0611776208c0718006311c84b4a113335c70d1f5c7c5dd93a5625c4af51c48847abd0b590c055306162d2a03ca1cbf7bcc101485641b2cd02f3e088df92eb020101a771feb6079155a30c57861e6e3683c3b94d1ccd0182ef6ced84b2ccb2ddf717172e1b4c3371b83c92e780dd8cebff92660c746573742073656e642d7478b18c974a6cb2dfcfb1f2e664005c3637eb8b17def29c947c5e816d00b48367438a705e62abb3a4f62e482b3ead3d15058f0e51102e9b3eab12f6f3058d679cb22d6dd8a2ce6202d9853a84690d4293f8c328e972758ed7e08488a942372462920568b52d42e02928b8a7ccd3391ac460ca0e3b032ba624d8e7e2010d1cd31392819f5a24b9c6370afcf0e777cef97c82011fd50521d007a8ead3e7fc891001012109b6290b25650951067871b027068bfc439418014a2916f7d8aa26061311f42a594d59ddfd6aaf6ab9cc91c796a1801c0c746573742073656e642d7478977353953a0925b0d1046b94a5fd7ea9766db30a94237a08dd15daf77977b4fa1d16b776c1181c9193b2bd72b960bf0298f011ca48cca9850225c6002a5e976b2ffb2e6f63b62b183e74a2d236c69d01847f51493d04d6d3f9e991e4552e47660e54cd8a14dc07c093448aac18e0551883eedb6bb97a44a549e65dcde8969763ab62eb15600629931f880fab20f8621c")
	b, err := FromBytes(d)
	assert.NoError(t, err)
	assert.Equal(t, b.SerializeSize(), len(d))
	d2, _ := b.Bytes()
	assert.Equal(t, d, d2)

	headerSize := b.Header().SerializeSize()
	certSize := b.PrevCertificate().SerializeSize()
	headerData := d[:headerSize]
	certData := d[headerSize : headerSize+certSize]
	certHash := hash.CalcHash(certData)

	txHashes := make([]hash.Hash, 0)
	for _, trx := range b.Transactions() {
		txHashes = append(txHashes, trx.ID())
	}
	txRoot := simplemerkle.NewTreeFromHashes(txHashes).Root()

	hashData := headerData
	hashData = append(hashData, certHash.Bytes()...)
	hashData = append(hashData, txRoot.Bytes()...)
	hashData = append(hashData, util.Int32ToSlice(int32(b.Transactions().Len()))...)

	expected1 := hash.CalcHash(hashData)
	expected2, _ := hash.FromString("c856bd7d0ce99842ec450ffad6416c45320d667b8b7c4927727860a43290595c")
	assert.Equal(t, b.Hash(), expected1)
	assert.Equal(t, b.Hash(), expected2)
	assert.Equal(t, b.Stamp(), hash.Stamp{0xc8, 0x56, 0xbd, 0x7d})
}

func TestMakeBlock(t *testing.T) {
	tmp := GenerateTestBlock(nil, nil)
	t.Run("Valid block information, should not panic", func(t *testing.T) {
		MakeBlock(1, util.Now(), tmp.Transactions(),
			tmp.Header().PrevBlockHash(),
			tmp.Header().StateRoot(),
			tmp.PrevCertificate(),
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})

	t.Run("Invalid block information, should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// Certificate is missed,
		MakeBlock(1, util.Now(), tmp.Transactions(),
			tmp.Header().PrevBlockHash(),
			tmp.Header().StateRoot(),
			nil,
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})

	t.Run("Invalid block information, should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// Invalid state root
		MakeBlock(1, util.Now(), tmp.Transactions(),
			tmp.Header().PrevBlockHash(),
			hash.UndefHash,
			tmp.PrevCertificate(),
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})

	t.Run("Invalid block information, should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("The code did not panic")
			}
		}()

		// Invalid previous block hash
		MakeBlock(1, util.Now(), tmp.Transactions(),
			hash.UndefHash,
			tmp.Header().PrevBlockHash(),
			tmp.PrevCertificate(),
			tmp.Header().SortitionSeed(),
			tmp.Header().ProposerAddress())
	})
}
