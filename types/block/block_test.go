package block_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/pactus-project/pactus/crypto/hash"
	"github.com/pactus-project/pactus/types/block"
	"github.com/pactus-project/pactus/types/certificate"
	"github.com/pactus-project/pactus/util"
	"github.com/pactus-project/pactus/util/simplemerkle"
	"github.com/pactus-project/pactus/util/testsuite"
	"github.com/stretchr/testify/assert"
)

// TODO: check error type.
func TestBasicCheck(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	t.Run("No transactions", func(t *testing.T) {
		b0 := ts.GenerateTestBlock(nil, nil)
		b := block.NewBlock(b0.Header(), b0.PrevCertificate(), block.Txs{})

		assert.Error(t, b.BasicCheck())
	})

	t.Run("Without the previous certificate", func(t *testing.T) {
		b0 := ts.GenerateTestBlock(nil, nil)
		b := block.NewBlock(b0.Header(), nil, b0.Transactions())

		assert.Error(t, b.BasicCheck())
	})

	t.Run("Invalid certificate", func(t *testing.T) {
		b0 := ts.GenerateTestBlock(nil, nil)
		cert0 := b0.PrevCertificate()
		invCert := certificate.NewCertificate(0, 0, cert0.Committers(), cert0.Absentees(), cert0.Signature())
		b := block.NewBlock(b0.Header(), invCert, b0.Transactions())

		assert.Error(t, b.BasicCheck())
	})

	t.Run("Invalid transaction", func(t *testing.T) {
		b0 := ts.GenerateTestBlock(nil, nil)
		trxs0 := b0.Transactions()
		invalidSigner := ts.RandSigner()
		invalidSigner.SignMsg(trxs0[0])
		b := block.NewBlock(b0.Header(), b0.PrevCertificate(), trxs0)

		assert.Error(t, b.BasicCheck())
	})

	t.Run("Invalid state root hash", func(t *testing.T) {
		d := ts.DecodingHex(
			"0134ec9b649f1fe7230ecf98ede2eb097587f69d7e29fade7d1ba0d5d9383bd74c1704655b00000000000000000000000000000000000000" +
				"00000000000000000000000000b9f39a3a63edeeeb24cec7daae01575168024a96b15def00da7ff84332acba24e3269dcfb5a592824ec174" +
				"5551ffcbb9011c8e4ed0bfe587a05f2bffa6965455a6fd483a1a010000000100040312142101128d1e40bd0135faab80e264d2171b1b6b87" +
				"00e87e1c7520e3f8c15930c7615eee6ef829ed5034c7d9254092f58b32d7ef010180b04fda9903c0f5ccf018010128d651affb92e704fcb9" +
				"93249c076e276fa89f4301e4b5fa19f70684a771b6587ae2b4fcd14ece4dc3adb2be84f4f0010c746573742073656e642d7478a56ee5b7af" +
				"5a4baeddd68a82556ca09ca0b1e3ae2587a0d202f5204193cf582a38a2e9bd08468a0321e6fd723741924c8144710257ea8629e7181a478f" +
				"348664e6484a20655b2eca7aec933e53d537fc082b34da84b899aba4aa6a980dcd15e80cebc5837cf48ee01d4a65bad3e272855c6f9d4d2b" +
				"a9b7b8ebce45f3e9113be1337d7e67d9fc9029705e0c2d3d824b31")

		b, err := block.FromBytes(d)
		assert.NoError(t, err)

		assert.Error(t, b.BasicCheck())
	})

	t.Run("Invalid previous block hash", func(t *testing.T) {
		d := ts.DecodingHex(
			"0134ec9b64000000000000000000000000000000000000000000000000000000000000000082916efdb068bbec819457f5ce73bf8d2ca743" +
				"e12c1946ee844e696e47bbe164b9f39a3a63edeeeb24cec7daae01575168024a96b15def00da7ff84332acba24e3269dcfb5a592824ec174" +
				"5551ffcbb9011c8e4ed0bfe587a05f2bffa6965455a6fd483a1a010000000100040312142101128d1e40bd0135faab80e264d2171b1b6b87" +
				"00e87e1c7520e3f8c15930c7615eee6ef829ed5034c7d9254092f58b32d7ef010180b04fda9903c0f5ccf018010128d651affb92e704fcb9" +
				"93249c076e276fa89f4301e4b5fa19f70684a771b6587ae2b4fcd14ece4dc3adb2be84f4f0010c746573742073656e642d7478a56ee5b7af" +
				"5a4baeddd68a82556ca09ca0b1e3ae2587a0d202f5204193cf582a38a2e9bd08468a0321e6fd723741924c8144710257ea8629e7181a478f" +
				"348664e6484a20655b2eca7aec933e53d537fc082b34da84b899aba4aa6a980dcd15e80cebc5837cf48ee01d4a65bad3e272855c6f9d4d2b" +
				"a9b7b8ebce45f3e9113be1337d7e67d9fc9029705e0c2d3d824b31")

		_, err := block.FromBytes(d)
		assert.Error(t, err)
	})

	t.Run("Invalid proposer address (type is 2)", func(t *testing.T) {
		d := ts.DecodingHex(
			"0134ec9b649f1fe7230ecf98ede2eb097587f69d7e29fade7d1ba0d5d9383bd74c1704655b82916efdb068bbec819457f5ce73bf8d2ca743" +
				"e12c1946ee844e696e47bbe164b9f39a3a63edeeeb24cec7daae01575168024a96b15def00da7ff84332acba24e3269dcfb5a592824ec174" +
				"5551ffcbb9021c8e4ed0bfe587a05f2bffa6965455a6fd483a1a010000000100040312142101128d1e40bd0135faab80e264d2171b1b6b87" +
				"00e87e1c7520e3f8c15930c7615eee6ef829ed5034c7d9254092f58b32d7ef010180b04fda9903c0f5ccf018010128d651affb92e704fcb9" +
				"93249c076e276fa89f4301e4b5fa19f70684a771b6587ae2b4fcd14ece4dc3adb2be84f4f0010c746573742073656e642d7478a56ee5b7af" +
				"5a4baeddd68a82556ca09ca0b1e3ae2587a0d202f5204193cf582a38a2e9bd08468a0321e6fd723741924c8144710257ea8629e7181a478f" +
				"348664e6484a20655b2eca7aec933e53d537fc082b34da84b899aba4aa6a980dcd15e80cebc5837cf48ee01d4a65bad3e272855c6f9d4d2b" +
				"a9b7b8ebce45f3e9113be1337d7e67d9fc9029705e0c2d3d824b31")

		b, err := block.FromBytes(d)
		assert.NoError(t, err)

		assert.Error(t, b.BasicCheck())
	})

	t.Run("Ok", func(t *testing.T) {
		b := ts.GenerateTestBlock(nil, nil)
		assert.NoError(t, b.BasicCheck())
		assert.LessOrEqual(t, b.Header().Time(), time.Now())
		assert.NotZero(t, b.Header().UnixTime())
		assert.Equal(t, b.Header().Version(), uint8(1))
	})
}

func TestCBORMarshaling(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	b1 := ts.GenerateTestBlock(nil, nil)
	bz1, err := cbor.Marshal(b1)
	assert.NoError(t, err)
	var b2 block.Block
	err = cbor.Unmarshal(bz1, &b2)
	assert.NoError(t, err)
	assert.NoError(t, b2.BasicCheck())
	assert.Equal(t, b1.Hash(), b2.Hash())

	assert.Equal(t, b1.Hash(), b2.Hash())

	err = cbor.Unmarshal([]byte{1}, &b2)
	assert.Error(t, err)
}

func TestEncodingBlock(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk := ts.GenerateTestBlock(nil, nil)
	length := blk.SerializeSize()

	for i := 0; i < length; i++ {
		w := util.NewFixedWriter(i)
		assert.Error(t, blk.Encode(w), "encode test %v failed", i)
	}
	w := util.NewFixedWriter(length)
	assert.NoError(t, blk.Encode(w))

	for i := 0; i < length; i++ {
		blk2 := new(block.Block)
		r := util.NewFixedReader(i, w.Bytes())
		assert.Error(t, blk2.Decode(r), "decode test %v failed", i)
	}

	blk2 := new(block.Block)
	r := util.NewFixedReader(length, w.Bytes())
	assert.NoError(t, blk2.Decode(r))
	assert.Equal(t, blk.Hash(), blk2.Hash())
	assert.Equal(t, blk.Header(), blk2.Header())
}

func TestTxFromBytes(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	blk := ts.GenerateTestBlock(nil, nil)
	bs, _ := blk.Bytes()
	_, err := block.FromBytes(bs)
	assert.NoError(t, err)
	_, err = blk.Bytes()
	assert.NoError(t, err)

	_, err = block.FromBytes([]byte{1})
	assert.Error(t, err)
}

func TestBlockHash(t *testing.T) {
	d, _ := hex.DecodeString(
		"0140da9b641551048b59a859946ca7f9ab95c9cf84da488a1a5c49ba643b29b653dc223bc20a4e9ff03158165f3d42" +
			"4e2a74677bfe24a7295d1ce2e55ca3644cbe9a5a5e7d913b8e1ba6a020afbd5a25024a12b37cf8e1ed0b9498f91d75b294db0f95123d8593" +
			"05aa5deea3d4216777e74310b6a601bb4d4d6b13c9b295781ab1533aea032978d4f89300001234050004060f1b23010fab4f72234cc7c120" +
			"48bbbc616c005573d8ad4d5c6997996d6f488946cdd78410f0a400c4a7f9bdb41506bdf717a892fa00")
	b, err := block.FromBytes(d)
	assert.NoError(t, err)
	assert.Equal(t, b.SerializeSize(), len(d))
	d2, _ := b.Bytes()
	assert.Equal(t, d, d2)

	headerSize := b.Header().SerializeSize()
	headerData := d[:headerSize]
	certSize := b.PrevCertificate().SerializeSize()
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
	expected2, _ := hash.FromString("03b77a2ac3d944c6fbef6a7e927ce285491bb5be413472adae56dee1f9ee3c97")
	assert.Equal(t, b.Hash(), expected1)
	assert.Equal(t, b.Hash(), expected2)
	assert.Equal(t, b.Stamp(), hash.Stamp{0x03, 0xb7, 0x7a, 0x2a})
}

func TestMakeBlock(t *testing.T) {
	ts := testsuite.NewTestSuite(t)

	b0 := ts.GenerateTestBlock(nil, nil)
	b1 := block.MakeBlock(1, util.Now(), b0.Transactions(),
		b0.Header().PrevBlockHash(),
		b0.Header().StateRoot(),
		b0.PrevCertificate(),
		b0.Header().SortitionSeed(),
		b0.Header().ProposerAddress())

	assert.Equal(t, b0.Hash(), b1.Hash())
}
